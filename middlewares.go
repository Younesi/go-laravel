package atlas

import "net/http"

func (a *Atlas) SessionLoad(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}
