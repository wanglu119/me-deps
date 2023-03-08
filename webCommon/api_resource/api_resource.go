package api_resource

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
	"github.com/spf13/afero"
	"github.com/wanglu119/me-deps/webCommon"
	"github.com/wanglu119/me-deps/webCommon/api_resource/files"
	"github.com/wanglu119/me-deps/webCommon/api_resource/fileutils"
)

var BaseScope string

func init() {
	var err error
	BaseScope, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func GetScope(subDir string) (scope string) {
	scope = filepath.Join(BaseScope, subDir)
	err := os.MkdirAll(BaseScope, fs.ModeDir)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	return scope
}

func ResourcePostHandler() webCommon.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()

		defer func() {
			_, _ = io.Copy(ioutil.Discard, r.Body)
		}()

		// Directories creation on POST.
		if strings.HasSuffix(r.URL.Path, "/") {
			err := d.GetFs().MkdirAll(r.URL.Path, 0775)
			if err != nil {
				log.Error(err)
				webCommon.ProcError(res, err)
				return
			}
			return
		}

		_, err := files.NewFileInfo(files.FileOptions{
			Fs:     d.GetFs(),
			Path:   r.URL.Path,
			Modify: true,
			Expand: false,
		})

		var file io.Reader

		result := []string{}
		ctype := r.Header.Get("Content-Type")
		randomName := r.Header.Get("RandomName")
		if strings.Contains(ctype, "multipart/form-data") {
			r.ParseMultipartForm(1024 * 1024)

			for filename, fhs := range r.MultipartForm.File {
				filepath := path.Join(r.URL.Path, filename)
				if randomName == "true" {
					ext := path.Ext(filename)
					filename = xid.New().String() + "_" + base64.StdEncoding.EncodeToString([]byte(filename)) + ext
					filepath = path.Join(r.URL.Path, filename)
				}

				if len(fhs) > 0 {
					file, err := fhs[0].Open()
					if err != nil {
						_ = d.GetFs().RemoveAll(r.URL.Path)
						log.Error(err)
						webCommon.ProcError(res, err)
						return
					}
					_, err = writeFile(d.GetFs(), filepath, file)
					file.Close()
					if err != nil {
						_ = d.GetFs().RemoveAll(r.URL.Path)
						log.Error(err)
						webCommon.ProcError(res, err)
						return
					}
					result = append(result, filepath)
				}
			}
		} else {
			if err == nil {
				log.Info(r.URL.Path)
				if r.URL.Query().Get("override") != "true" {
					res.Status = http.StatusConflict
					return
				}
			}

			file = r.Body
			_, err = writeFile(d.GetFs(), r.URL.Path, file)
			if err != nil {
				_ = d.GetFs().RemoveAll(r.URL.Path)
				log.Error(err)
				webCommon.ProcError(res, err)
				return
			}
			result = append(result, r.URL.Path)
		}

		res.Data = result
	}
}

func writeFile(fs afero.Fs, dst string, in io.Reader) (os.FileInfo, error) {
	dir, _ := path.Split(dst)
	err := fs.MkdirAll(dir, 0775)
	if err != nil {
		return nil, err
	}

	file, err := fs.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, in)
	if err != nil {
		return nil, err
	}

	// Gets the info about the file.
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func ResourceDeleteHandler() webCommon.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()

		if r.URL.Path == "/" {
			res.Status = http.StatusForbidden
			return
		}

		_, err := files.NewFileInfo(files.FileOptions{
			Fs:     d.GetFs(),
			Path:   r.URL.Path,
			Modify: true,
			Expand: false,
			// Checker: d,
		})
		if err != nil {
			log.Error(err)
			webCommon.ProcError(res, err)
			return
		}

		err = d.GetFs().RemoveAll(r.URL.Path)
		if err != nil {
			log.Error(err)
			webCommon.ProcError(res, err)
			return
		}
	}
}

func ResourceGetHandler() webCommon.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()

		var err error
		defer func() {
			if err != nil {
				webCommon.ProcError(res, err)
				return
			}
		}()

		file, err := files.NewFileInfo(files.FileOptions{
			Fs:     d.GetFs(),
			Path:   r.URL.Path,
			Modify: true,
			Expand: true,
			// Checker: d,
		})
		if err != nil {
			if !os.IsExist(err) {
				err = nil
				return
			}
			log.Error(err)
			return
		}

		sorting := files.Sorting{}
		err = json.NewDecoder(r.Body).Decode(&sorting)
		if err != nil {
			log.Error(err)
			return
		}

		if file.IsDir {
			file.Listing.Sorting = sorting
			file.Listing.ApplySort()
			res.Data = file
			return
		}

		if checksum := r.URL.Query().Get("checksum"); checksum != "" {
			err = file.Checksum(checksum)
			if err != nil {
				log.Error(err)
				return
			}

			// do not waste bandwidth if we just want the checksum
			file.Content = ""
		}

		res.Data = file
	}
}

func checkParent(src, dst string) error {
	rel, err := filepath.Rel(src, dst)
	if err != nil {
		return err
	}

	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, "../") && rel != ".." && rel != "." {
		return errors.New("ErrSourceIsParent")
	}

	return nil
}

func addVersionSuffix(path string, fs afero.Fs) string {
	counter := 1
	dir, name := filepath.Split(path)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	for {
		if _, err := fs.Stat(path); err != nil {
			break
		}
		renamed := fmt.Sprintf("%s(%d)%s", base, counter, ext)
		path = filepath.ToSlash(dir) + renamed
		counter++
	}

	return path
}

func ResourcePatchHandler() webCommon.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()

		var err error
		defer func() {
			if err != nil {
				webCommon.ProcError(res, err)
				return
			}
		}()

		src := r.URL.Path
		dst := r.URL.Query().Get("destination")
		action := r.URL.Query().Get("action")

		if dst == "/" || src == "/" {
			res.Status = http.StatusForbidden
			return
		}
		if err = checkParent(src, dst); err != nil {
			log.Error(err)
			res.Status = http.StatusBadRequest
			return
		}

		override := r.URL.Query().Get("override") == "true"
		rename := r.URL.Query().Get("rename") == "true"
		if !override && !rename {
			if _, err = d.GetFs().Stat(dst); err == nil {
				res.Status = http.StatusConflict
				return
			}
		}

		if rename {
			dst = addVersionSuffix(dst, d.GetFs())
		}

		switch action {
		// TODO: use enum
		case "copy":
			err = fileutils.Copy(d.GetFs(), src, dst)
			if err != nil {
				log.Error(err)
			}
			return
		case "rename":
			dst = filepath.Clean("/" + dst)
			err = d.GetFs().Rename(src, dst)
			if err != nil {
				log.Error(err)
			}
			return
		default:
			err = errors.New(fmt.Sprintf("unsupported action %s", action))
			return
		}
	}
}
