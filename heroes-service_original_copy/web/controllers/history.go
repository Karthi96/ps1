package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (app *Application) HistoryHandler(w http.ResponseWriter, r *http.Request) {

	type dataObject struct {
		RollNo        string
		Name          string
		Course        string
		Details       string
		Date          string
		Grade         string
		TransactionID string `json:"txtid"`
		LastModified  string `json:"timestamp"`
		Success       bool
		Response      bool
		ErrorMessage  string
	}
	var data []dataObject
	if r.FormValue("submitted") == "true" {

		RNo := r.FormValue("rollno")
		Course := r.FormValue("course")
		fmt.Println("Requesting Webservie call")
		response, err := http.Get("http://localhost:3000/api/powerschool?methodcall=historydetails&rollno=" + url.PathEscape(strings.ToLower(RNo)) + "&course=" + url.PathEscape(Course))

		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		str, _ := ioutil.ReadAll(response.Body)
		if response.StatusCode == 200 {

			err = json.Unmarshal([]byte(string(str)), &data)
			if err != nil {
				http.Error(w, "Unable UNMARSHAL"+err.Error(), 500)
			}
			fmt.Println(data)
			for k := range data {
				data[k].Success = true
				data[k].Response = true
			}

		} else if response.StatusCode == 500 {

			err = json.Unmarshal([]byte("[{\"course\":\"\",\"date\":\"\",\"details\":\"\",\"name\":\"\",\"rollno\":\"\",\"timestamp\":\"\",\"txtid\":\"\",\"Success\":false,\"Response\":true,\"ErrorMessage\":\"Unable to find the trace history for this RollNo="+RNo+" for this Course="+Course+". Please try with different Rollno/Course\"}]"), &data)
			//data[0].ErrorMessage = string(str)
			//fmt.Println("Error" + err.Error())
		}

	}
	renderTemplate(w, r, "history.html", data)
}
