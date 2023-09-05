package webCommon

import (
	"net/http"

	"github.com/spf13/afero"
)

type WebData interface {
	GetResponse() *Response
	SetResponse(*Response)
	GetAuthToken() string
	GetAuthData() interface{}
	SetAuthData(interface{})
	SetHttpHeader(http.Header)
	SetHttpRequest(*http.Request)

	GetFs() afero.Fs
}

type WebDataImplPart struct {
	Resp   *Response
	Header http.Header
	Req    *http.Request
}

func (wdip *WebDataImplPart) GetResponse() *Response {
	return wdip.Resp
}

func (wdip *WebDataImplPart) SetResponse(resp *Response) {
	wdip.Resp = resp
}
func (wdip *WebDataImplPart) SetHttpHeader(header http.Header) {
	wdip.Header = header
}
func (wdip *WebDataImplPart) SetHttpRequest(req *http.Request) {
	wdip.Req = req
}

type HandleFunc func(w http.ResponseWriter, r *http.Request, d WebData)
type CreateWebDataFunc func() WebData
type WebHandler struct {
	CreateWebData CreateWebDataFunc
}

func (wh *WebHandler) handle(fn HandleFunc, prefix string) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := wh.CreateWebData()
		res := NewResponse(w)
		d.SetResponse(res)
		d.SetHttpHeader(r.Header)
		d.SetHttpRequest(r)
		defer res.SendJson()
		fn(w, r, d)
	})

	return http.StripPrefix(prefix, handler)
}

func (wh *WebHandler) Monkey(fn HandleFunc, prefix string) http.Handler {
	return wh.handle(fn, prefix)
}
