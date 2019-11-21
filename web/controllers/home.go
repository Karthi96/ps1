package controllers

import (
	"net/http"
)

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, r, "home.html", "")
}
