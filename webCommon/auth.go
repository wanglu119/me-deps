package webCommon

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type BasicAuthVertifyFunc func(username string, password string, d WebData, r *http.Request) bool
type BasicAuth struct {
	AuthVertify BasicAuthVertifyFunc
}

func (ba *BasicAuth) WithBasicAuth(fn HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d WebData) {
		res := d.GetResponse()

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			res.Status = http.StatusNonAuthoritativeInfo
			res.Data = "Authorization not in header"
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			errMsg := fmt.Sprintf("%v", err)
			log.Error(errMsg)
			res.Status = http.StatusInternalServerError
			res.Data = errMsg
			return
		}
		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			res.Status = http.StatusNotAcceptable
			res.Data = "parse basic auth error"
			return
		}

		if ba.AuthVertify(pair[0], pair[1], d, r) {
			d.SetResponse(res)
			fn(w, r, d)
			return
		} else {
			if res.Status == 200 {
				res.Status = http.StatusNotAcceptable
				res.Data = "auth fail"
			}
			return
		}
	}
}

// =======================================================

type Extractor []string

func (e Extractor) ExtractToken(r *http.Request) (string, error) {
	token, _ := request.HeaderExtractor{"X-Auth"}.ExtractToken(r)

	// Checks if the token isn't empty and if it contains two dots.
	// The former prevents incompatibility with URLs that previously
	// used basic auth.
	if token != "" && strings.Count(token, ".") == 2 {
		return token, nil
	}

	auth := r.URL.Query().Get("auth")
	if auth == "" {
		return "", request.ErrNoTokenInRequest
	}

	return auth, nil
}

type JwtAuthClaims struct {
	*jwt.StandardClaims
	AuthData interface{} `json:"authData"`
}

type JwtAuth struct {
	ExpSec uint64
}

func (ja *JwtAuth) ExtractToken(r *http.Request) (string, error) {
	extractor := &Extractor{}
	return extractor.ExtractToken(r)
}

func (ja *JwtAuth) VerifyToken(r *http.Request, authToken string) (bool, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(authToken), nil
	}

	standardClaims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    "wl",
	}
	token, err := request.ParseFromRequest(r, &Extractor{}, keyFunc, request.WithClaims(standardClaims))
	if err != nil {
		log.Error(err)
		return false, err
	}
	if !token.Valid {
		err = errors.New("valid to token fail")
		log.Error(err)
		return false, err
	}
	expired := !standardClaims.VerifyExpiresAt(time.Now().Add(24*time.Hour).Unix(), true)
	if expired {
		err = errors.New("token expired")
		log.Error(err)
		return false, err
	}
	return true, nil
}

func (ja *JwtAuth) WithJwt(fn HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d WebData) {
		res := d.GetResponse()

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			jwtAuthClaims, _ := token.Claims.(*JwtAuthClaims)
			d.SetAuthData(jwtAuthClaims.AuthData)
			return []byte(d.GetAuthToken()), nil
		}

		standardClaims := ja.createJwtStandardClaims(d)
		token, err := request.ParseFromRequest(r, &Extractor{}, keyFunc, request.WithClaims(standardClaims))

		if err != nil || !token.Valid {
			errMsg := fmt.Sprintf("%v", err)
			log.Error(errMsg)
			res.Data = errMsg
			res.Status = http.StatusForbidden
			return
		}

		expired := !standardClaims.VerifyExpiresAt(time.Now().Add(24*time.Hour).Unix(), true)
		if expired {
			w.Header().Add("X-Renew-Token", "true")
		}
		fn(w, r, d)
	}
}

func (ja *JwtAuth) createJwtStandardClaims(d WebData) *JwtAuthClaims {
	expSec := ja.ExpSec
	if expSec == 0 {
		expSec = 24 * 60 * 60
	}
	standardClaims := &JwtAuthClaims{
		&jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(expSec) * time.Second).Unix(),
			Issuer:    "wl",
		},
		d.GetAuthData(),
	}
	return standardClaims
}

func (ja *JwtAuth) GenSigned(d WebData) (string, error) {
	standardClaims := ja.createJwtStandardClaims(d)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims)
	signed, err := token.SignedString([]byte(d.GetAuthToken()))

	return signed, err
}

func GenSignedByToken(t string, username string) (string, error) {
	standardClaims := &JwtAuthClaims{
		&jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "wl",
		},
		map[string]string{
			"username": username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims)
	signed, err := token.SignedString([]byte(t))

	return signed, err
}

// ==========================Captcha auth==============================
type CaptchaAuth struct {
}

func (ca *CaptchaAuth) WithCaptcha(fn HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d WebData) {
		res := d.GetResponse()

		captchaId := r.Header.Get("captchaId")
		captchaVal := r.Header.Get("captchaVal")

		if !captcha.VerifyString(captchaId, captchaVal) {
			errMsg := "verifyCode Error"
			log.Error(errMsg)
			res.Data = errMsg
			res.Status = http.StatusPreconditionRequired
			return
		}

		fn(w, r, d)
	}
}
