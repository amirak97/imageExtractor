package web

import (
	"fmt"
	"html/template"
	"net/http"
)

// HomeHandler serves different pages based on session
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r.Context())
	// check if user has provided a link before (fake check)
	if session.Cookie == "" {
		tmpl := `<html><body>
		<form method="POST" action="/submit">
			<input type="text" name="link" placeholder="Enter your link" />
			<input type="submit" value="Save" />
		</form>
		</body></html>`
		w.Write([]byte(tmpl))
		return
	}

	// otherwise show user page
	tmpl := template.Must(template.New("user").Parse(`
	<html><body>
	<h1>Welcome back!</h1>
	<p>Your session ID: {{.ID}}</p>
	<p>Your cookie: {{.Cookie}}</p>
	</body></html>`))
	tmpl.Execute(w, session)
}

// SubmitHandler saves link for user (just demo)
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r.Context())
	link := r.FormValue("link")
	session.Cookie = link // just save link inside Cookie field
	w.Write([]byte(fmt.Sprintf("Link saved for session %s: %s", session.ID, link)))
}
