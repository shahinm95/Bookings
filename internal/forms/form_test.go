package forms

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestForm_New(t *testing.T) {
	var data url.Values
	var testType *Form
	test2 := struct {
		int
		string
	}{}
	returnForm := New(data)
	if reflect.TypeOf(testType) != reflect.TypeOf(returnForm) {
		t.Error("return type is not type of form!")
	}
	if reflect.TypeOf(test2) == reflect.TypeOf(returnForm) {
		t.Error("return type shouldn't be equal to a non form type")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)

	form := New(r.PostForm)

	isTrue := form.Required("a", "b", "c", "d")
	if form.Valid() {
		t.Error("form should't be valid when requried data is not in the form")
	}
	if isTrue {
		t.Error("required should't be true when requried data is not in the form")
	}
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")
	postData.Add("d", "a")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	isTrue = form.Required("a", "b", "c", "d")
	if !form.Valid() {
		t.Error("form should be valid but it is not!")
	}
	if !isTrue {
		t.Error("required should be true but it is not!")
	}
}

func TestForm_Has(t *testing.T) {
	var testfield string
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)
	rt := form.Has(testfield, r)
	if rt {
		t.Error("return type of Has is true for empty string!")
	}

	r = httptest.NewRequest("POST", "/whatever", nil)

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")
	postData.Add("d", "a")

	r.PostForm = postData
	form = New(r.PostForm)
	if !form.Has("a", r) {
		t.Error("form should have data for given field, but it doesn't !")
	}
	if form.Has("g", r) {
		t.Error("form shouldn't have data for given field but it does !")
	}

}

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid form where it should be valid!")
	}
}

func TestMinLength(t *testing.T) {
	var field string
	var length int
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)
	isTrue := form.MinLength(field, length, r)
	if !isTrue {
		t.Error("got false for length of value of empty string")
	}

	postData := url.Values{}
	postData.Add("name", "shahin")
	r.PostForm = postData
	form = New(r.PostForm)
	if !form.MinLength("name", 4, r) {
		t.Error("value of given field has more length than given length , so it sould be true")
	}
	isError := form.Errors.Get("name")
	if isError != "" {
		t.Error("should not have an error, but get one")
	}

	if form.MinLength("name", 7, r) {
		t.Error("value of given field has less length than given length , so it sould be false")
	}
	isError = form.Errors.Get("name")
	if isError == "" {
		t.Error("should have an error, but didn't get one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)

	postData := url.Values{}
	postData.Add("a", "")
	postData.Add("b", "shahin@gmail.com")
	postData.Add("c", "hello")

	r.PostForm = postData

	form := New(r.PostForm)
	if form.IsEmail("a") {
		t.Error("this string is not an email address so it should be false")
	}
	if !form.IsEmail("b") {
		t.Error("this string is an email address so it should be true")
	}
	if form.IsEmail("c") {
		t.Error("this string is not an email address so it should be false")
	}
}
