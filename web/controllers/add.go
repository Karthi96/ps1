package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (app *Application) AddHandler(w http.ResponseWriter, r *http.Request) {

	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
		ErrorMessage  string
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
		ErrorMessage:  "",
	}
	if r.FormValue("submitted") == "true" {
		rollno := r.FormValue("rollno")
		name := r.FormValue("name")
		course := r.FormValue("course")
		details := r.FormValue("details")
		date := r.FormValue("date")
		grade := r.FormValue("grade")

		var args []string
		args = append(args, "createStudent")
		args = append(args, rollno)
		args = append(args, name)
		args = append(args, course)
		args = append(args, details)
		args = append(args, date)
		args = append(args, grade)

		pathURL := "http://localhost:3000/api/powerschool?methodcall=putdetails&rollno=" + url.PathEscape(rollno) + "&name=" + url.PathEscape(name) + "&&course=" + url.PathEscape(course) + "&details=" + url.PathEscape(details) + "&date=" + url.PathEscape(date) + "&grade=" + url.PathEscape(grade)
		fmt.Println("Requesting Webservice call " + pathURL)
		response, err := http.Get(pathURL)
		fmt.Println("Response code=")
		fmt.Print(response.StatusCode)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data1, _ := ioutil.ReadAll(response.Body)
			str := string(data1)
			if strings.Contains(str, "This rollno already exists") {
				data.ErrorMessage = "This Rollno: " + rollno + " is already doing this Course: " + course + ". PLease try with different Rollno/Course."
				data.TransactionId = ""
				data.Success = false
				data.Response = true
			} else if response.StatusCode == 200 {
				data.TransactionId = str
				data.Success = true
				data.Response = true
				data.ErrorMessage = ""
			} else {
				data.ErrorMessage = str
				data.TransactionId = ""
				data.Success = false
				data.Response = true
			}
		}

	}
	renderTemplate(w, r, "add.html", data)
}
