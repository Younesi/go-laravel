package atlas

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (a *Atlas) SessionLoad(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}

func (a *Atlas) NoCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(a.config.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")
	
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteDefaultMode,
		Domain:   a.config.cookie.domain,
	})

	return csrfHandler
}
