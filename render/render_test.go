package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
}{
	{
		name:          "go_page",
		renderer:      "go",
		template:      "home",
		errorExpected: false,
	},
	{
		name:          "jet_page",
		renderer:      "go",
		template:      "home",
		errorExpected: false,
	},
	{
		name:          "go_page_none_existing_test",
		renderer:      "go",
		template:      "",
		errorExpected: true,
	},
	{
		name:          "go_page",
		renderer:      "jet",
		template:      "",
		errorExpected: true,
	},
	{
		name:          "go_page",
		renderer:      "not-supported-renderer",
		template:      "",
		errorExpected: true,
	},
}

func TestRender_Page(t *testing.T) {
	for _, e := range pageData {
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()

		testRenderer.Renderer = e.renderer
		testRenderer.RootPath = "./testdata"

		err = testRenderer.Page(w, r, e.template, nil, nil)
		if e.errorExpected {
			if err == nil {
				t.Errorf("%s expected error, but got no error", e.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s not expected error, but got :  %s", e.name, err.Error())
			}
		}
	}
}

func TestRender_GoPage(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.GoPage(w, r, "home", nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}

	err = testRenderer.GoPage(w, r, "no-file", nil)
	if err == nil {
		t.Error("Expected error rendering no existing page")
	}
}

func TestRender_JetPage(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	testRenderer.Renderer = "jet"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.JetPage(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}

	err = testRenderer.JetPage(w, r, "no-file", nil, nil)
	if err == nil {
		t.Error("Expected error rendering no existing jet page")
	}
}
