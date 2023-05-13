package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shahinm95/bookings/internal/config"
	"github.com/shahinm95/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig
func TestMain(m *testing.M) {
	//what I'm going to put in session 
	gob.Register(models.Reservation{})
		// set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = false
	
		testApp.Session = session
		app = &testApp // this would give app inside render.go , session it needs in order to execute the test

		os.Exit(m.Run())
}