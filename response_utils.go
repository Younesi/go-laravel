package atlas

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

func (at *Atlas) WriteJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (at *Atlas) WriteXml(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (at *Atlas) DownloadFile(w http.ResponseWriter, r *http.Request, filePath, fileName string) error {
	fp := path.Join(filePath, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

func (at *Atlas) ErrNotFound(w http.ResponseWriter, r *http.Request) {
	at.ErrStatus(w, http.StatusNotFound)
}

func (at *Atlas) ErrInternalServer(w http.ResponseWriter, r *http.Request) {
	at.ErrStatus(w, http.StatusInternalServerError)
}

func (at *Atlas) ErrUnauthorized(w http.ResponseWriter, r *http.Request) {
	at.ErrStatus(w, http.StatusUnauthorized)
}
func (at *Atlas) ErrForbidden(w http.ResponseWriter, r *http.Request) {
	at.ErrStatus(w, http.StatusForbidden)
}

func (at *Atlas) ErrStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
