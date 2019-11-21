package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (app *Application) RequestHandler(w http.ResponseWriter, r *http.Request) {
	type dataObject struct {
		RollNo        string
		Name          string
		Course        string
		Details       string
		Date          string
		Grade         string
		Success       bool
		Response      bool
		ErrorMessage  string
		TransactionId string `json:"txtid"`
		EditCall      bool
	}

	var data []dataObject

	if r.FormValue("submitted") == "true" {

		RNo := r.FormValue("rollno")
		Course := r.FormValue("course")
		response, err := http.Get("http://localhost:3000/api/powerschool?methodcall=getdetails&rollno=" + url.PathEscape(strings.ToLower(RNo)) + "&course=" + url.PathEscape(Course))

		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain"+err.Error(), 500)
		}
		str, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("Output=" + string(str))

		if response.StatusCode == 200 {

			err = json.Unmarshal([]byte(string(str)), &data)
			if err != nil {
				http.Error(w, "Unable UNMARSHAL"+err.Error(), 500)
				for k := range data {
					data[k].Success = false
					data[k].ErrorMessage = err.Error()
					data[k].EditCall = false
				}
			}
			for k := range data {
				data[k].Success = true
				data[k].Response = true
				data[k].EditCall = false
			}

		} else {
			err = json.Unmarshal([]byte("[{\"course\":\"\",\"date\":\"\",\"details\":\"\",\"name\":\"\",\"rollno\":\"\",\"timestamp\":\"\",\"txtid\":\"\",\"Success\":false,\"Response\":true,\"ErrorMessage\":\"Unable to find the history for this rollno="+RNo+"\"}]"), &data)
			for k := range data {
				data[k].EditCall = false
				data[k].Success = false
				data[k].Response = true
				data[k].ErrorMessage = "Unable to find the RollNo=" + RNo + " for this Course=" + Course + " in blockchain. Please try with different RollNo/Course."
			}

		}

	} else if r.FormValue("submitted_edit") == "true" {

		rollno := r.FormValue("erollno")
		name := r.FormValue("ename")
		course := r.FormValue("ecourse")
		details := r.FormValue("edetails")
		date := r.FormValue("edate")
		grade := r.FormValue("egrade")

		pathURL := "http://localhost:3000/api/powerschool?methodcall=updatedetails&rollno=" + url.PathEscape(rollno) + "&name=" + url.PathEscape(name) + "&&course=" + url.PathEscape(course) + "&details=" + url.PathEscape(details) + "&date=" + url.PathEscape(date) + "&grade=" + url.PathEscape(grade)
		fmt.Println("Requesting Webservice call " + pathURL)
		response, err := http.Get(pathURL)
		fmt.Println("Response code=")
		fmt.Print(response.StatusCode)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data1, _ := ioutil.ReadAll(response.Body)
			str := string(data1)
			if response.StatusCode == 200 {
				err = json.Unmarshal([]byte("[{\"course\":\"\",\"date\":\"\",\"details\":\"\",\"name\":\"\",\"rollno\":\"\",\"timestamp\":\"\",\"txtid\":\"\",\"Success\":false,\"Response\":true}]"), &data)
				for k := range data {
					data[k].TransactionId = str
					data[k].Success = true
					data[k].Response = true
					data[k].ErrorMessage = ""
					data[k].EditCall = true
				}
			} else {
				err = json.Unmarshal([]byte("[{\"course\":\"\",\"date\":\"\",\"details\":\"\",\"name\":\"\",\"rollno\":\"\",\"timestamp\":\"\",\"txtid\":\"\",\"Success\":false,\"Response\":true}]"), &data)
				for k := range data {
					data[k].ErrorMessage = str
					data[k].TransactionId = ""
					data[k].Success = false
					data[k].Response = true
					data[k].EditCall = true
				}
			}
		}

	}

	renderTemplate(w, r, "request.html", data)
}
