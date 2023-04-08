package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)

	f := New(r.Form)
	isValid := f.Valid()
	if !isValid {
		t.Error("got invalid when should valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)

	f := New(r.Form)

	f.Required("a", "b")

	isValid := f.Valid()
	if isValid {
		t.Error("got valid when should invalid")
	}

	r = httptest.NewRequest(http.MethodPost, "/whatever", nil)

	f = New(r.Form)

	formData := url.Values{}
	formData.Add("a", "a")
	formData.Add("b", "b")
	f.Data = formData

	f.Required("a", "b")

	isValid = f.Valid()
	if !isValid {
		t.Error("got invalid when should valid")
	}
}

func TestForm_Has(t *testing.T) {
	formData := url.Values{}
	formData.Add("a", "a")

	f := New(formData)

	has := f.Has("a")
	if !has {
		t.Error("got invalid when should valid")
	}

	formData = url.Values{}
	f = New(formData)

	has = f.Has("a")
	if has {
		t.Error("got valid when should invalid")
	}
}

func TestForm_MinLength(t *testing.T) {
	formData := url.Values{}
	formData.Add("a", "a")

	f := New(formData)

	isValid := f.MinLength("a", 1)
	if !isValid {
		t.Error("got invalid when should valid")
	}

	formData = url.Values{}
	f = New(formData)

	isValid = f.MinLength("a", 1)
	if isValid {
		t.Error("got valid when should invalid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	formData := url.Values{}
	formData.Add("a", "a@email.com")

	f := New(formData)

	f.IsEmail("a")

	isValid := f.Valid()
	if !isValid {
		t.Error("got invalid when should valid")
	}

	formData = url.Values{}
	formData.Add("a", "a")

	f = New(formData)

	f.IsEmail("a")

	isValid = f.Valid()
	if isValid {
		t.Error("got valid when should invalid")
	}
}
