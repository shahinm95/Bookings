package render

import (
	"net/http"
	"testing"

	"github.com/shahinm95/bookings/internal/models"
)

func TestRenderTemplate(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td , r)
	if result.Flash != "123" {
		t.Error("flash should be 123 but not found in session")
	}
}

func getSession()(*http.Request, error){
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil , err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil

}