package api_resource

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	gopath "path"
	"path/filepath"
	"strings"

	"github.com/wanglu119/me-deps/webCommon"
	"github.com/wanglu119/me-deps/webCommon/api_resource/files"
	"github.com/wanglu119/me-deps/webCommon/api_resource/fileutils"

	"github.com/mholt/archiver"
	"github.com/spf13/afero"
)

func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return gopath.Clean(name)
}

func parseQueryFiles(r *http.Request, f *files.FileInfo) ([]string, error) {
	var fileSlice []string
	names := strings.Split(r.URL.Query().Get("files"), ",")

	if len(names) == 0 {
		fileSlice = append(fileSlice, f.Path)
	} else {
		for _, name := range names {
			name, err := url.QueryUnescape(strings.Replace(name, "+", "%2B", -1)) //nolint:govet
			if err != nil {
				return nil, err
			}

			name = slashClean(name)
			fileSlice = append(fileSlice, filepath.Join(f.Path, name))
		}
	}

	return fileSlice, nil
}

//nolint: goconst
func parseQueryAlgorithm(r *http.Request) (string, archiver.Writer, error) {
	// TODO: use enum
	switch r.URL.Query().Get("algo") {
	case "zip", "true", "":
		return ".zip", archiver.NewZip(), nil
	case "tar":
		return ".tar", archiver.NewTar(), nil
	case "targz":
		return ".tar.gz", archiver.NewTarGz(), nil
	case "tarbz2":
		return ".tar.bz2", archiver.NewTarBz2(), nil
	case "tarxz":
		return ".tar.xz", archiver.NewTarXz(), nil
	case "tarlz4":
		return ".tar.lz4", archiver.NewTarLz4(), nil
	case "tarsz":
		return ".tar.sz", archiver.NewTarSz(), nil
	default:
		return "", nil, errors.New("format not implemented")
	}
}

func setContentDisposition(w http.ResponseWriter, r *http.Request, file *files.FileInfo) {
	if r.URL.Query().Get("inline") == "true" {
		w.Header().Set("Content-Disposition", "inline")
	} else {
		// As per RFC6266 section 4.3
		w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(file.Name))
	}
}

func RawHandler() webCommon.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()

		file, err := files.NewFileInfo(files.FileOptions{
			Fs:     d.GetFs(),
			Path:   r.URL.Path,
			Modify: true,
			Expand: false,
			// Checker: d,
		})
		if err != nil {
			if !os.IsExist(err) {
				res.IsSend = false
				return
			}
			log.Error(err)
			webCommon.ProcError(res, err)
			return
		}

		if files.IsNamedPipe(file.Mode) {
			setContentDisposition(w, r, file)
			return
		}

		if !file.IsDir {
			if _, err := rawFileHandler(w, r, file); err != nil {
				log.Error(err)
				webCommon.ProcError(res, err)
				return
			} else {
				res.IsSend = false
				return
			}
		}

		if _, err := rawDirHandler(w, r, d, file); err != nil {
			log.Error(err)
			webCommon.ProcError(res, err)
			return
		}
	}
}

func addFile(ar archiver.Writer, d webCommon.WebData, path, commonPath string) error {

	info, err := d.GetFs().Stat(path)
	if err != nil {
		return err
	}

	var (
		file          afero.File
		arcReadCloser = ioutil.NopCloser(&bytes.Buffer{})
	)
	if !files.IsNamedPipe(info.Mode()) {
		file, err = d.GetFs().Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		arcReadCloser = file
	}

	if path != commonPath {
		filename := strings.TrimPrefix(path, commonPath)
		filename = strings.TrimPrefix(filename, string(filepath.Separator))
		err = ar.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   info,
				CustomName: filename,
			},
			ReadCloser: arcReadCloser,
		})
		if err != nil {
			return err
		}
	}

	if info.IsDir() {
		names, err := file.Readdirnames(0)
		if err != nil {
			return err
		}

		for _, name := range names {
			err = addFile(ar, d, filepath.Join(path, name), commonPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func rawDirHandler(w http.ResponseWriter, r *http.Request, d webCommon.WebData, file *files.FileInfo) (int, error) {
	filenames, err := parseQueryFiles(r, file)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extension, ar, err := parseQueryAlgorithm(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = ar.Create(w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer ar.Close()

	commonDir := fileutils.CommonPrefix(filepath.Separator, filenames...)

	name := filepath.Base(commonDir)
	if name == "." || name == "" || name == string(filepath.Separator) {
		name = file.Name
	}
	// Prefix used to distinguish a filelist generated
	// archive from the full directory archive
	if len(filenames) > 1 {
		name = "_" + name
	}
	name += extension
	w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(name))

	for _, fname := range filenames {
		err = addFile(ar, d, fname, commonDir)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return 0, nil
}

func rawFileHandler(w http.ResponseWriter, r *http.Request, file *files.FileInfo) (int, error) {
	isFresh := checkEtag(w, r, file.ModTime.Unix(), file.Size)
	if isFresh {
		return http.StatusNotModified, nil
	}

	fd, err := file.Fs.Open(file.Path)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer fd.Close()

	setContentDisposition(w, r, file)

	http.ServeContent(w, r, file.Name, file.ModTime, fd)
	return 0, nil
}

func checkEtag(w http.ResponseWriter, r *http.Request, fTime, fSize int64) bool {
	etag := fmt.Sprintf("%x%x", fTime, fSize)
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Etag", etag)

	return r.Header.Get("If-None-Match") == etag
}
