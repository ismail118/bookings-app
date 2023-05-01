package render

import (
	"github.com/ismail118/bookings-app/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	rest := AddDefaultData(&td, r)

	if rest.Flash != "123" {
		t.Errorf("want %s but got %s", "123", rest.Flash)
	}
}

func TestTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var w myWriter

	err = Template(&w, r, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error(err)
	}

	err = Template(&w, r, "non-exists.page.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("template non-exists.page.gohtml should not exists")
	}
}

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	if tc == nil {
		t.Error("Should not nil")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
