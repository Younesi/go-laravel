package session_test

import (
	"reflect"
	"testing"

	"github.com/alexedwards/scs/v2"

	sessionPackage "github.com/younesi/celeritas/session"
)

func TestSession_InitSession(t *testing.T) {
	session := &sessionPackage.Session{
		SessionType:    "cookie",
		CookieDomain:   "test",
		CookieName:     "test-name",
		CookieLifetime: "100",
		CookiePersist:  "true",
		CookieSecure:   "false",
	}

	var sm *scs.SessionManager

	ses := session.InitSession()

	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(ses)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		sessKind = rv.Kind()
		sessType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("invalid type or kind; kind:", rv.Kind(), "type:", rv.Type())
	}

	if sessKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned testing cookie session. Expected", reflect.ValueOf(sm).Kind(), "and got", sessKind)
	}

	if sessType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type returned testing cookie session. Expected", reflect.ValueOf(sm).Type(), "and got", sessType)
	}
}
