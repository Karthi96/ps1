package controllers

import (
	"fmt"
	"net/http"
)

type Book struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var books []Book

func (app *Application) PowerSchoolAPIHandler(w http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()
	MethodCall := keys.Get("methodcall")
	if len(MethodCall) > 0 {
		if MethodCall == "getdetails" {
			RollNo := keys.Get("rollno")
			Course := keys.Get("course")
			if len(RollNo) > 0 && len(Course) > 0 {
				var args []string
				args = append(args, "readStudent")
				args = append(args, RollNo)
				args = append(args, Course)
				txid, err := app.Fabric.GetStudents(args)
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				fmt.Println(txid)
				fmt.Fprintf(w, txid)
				w.Header().Set("Content-Type", "application/json")
				//json.NewEncoder(w).Encode(txid)
			} else {
				http.Error(w, "Page not found or check the parameters", 500)
			}

		} else if MethodCall == "putdetails" {
			RollNo := keys.Get("rollno")
			Name := keys.Get("name")
			Course := keys.Get("course")
			Details := keys.Get("details")
			Date := keys.Get("date")
			Grade := keys.Get("grade")
			if len(RollNo) > 0 && len(Name) > 0 && len(Course) > 0 && len(Details) > 0 && len(Date) > 0 && len(Grade) > 0 {

				var args []string
				args = append(args, "createStudent")
				args = append(args, RollNo)
				args = append(args, Name)
				args = append(args, Course)
				args = append(args, Details)
				args = append(args, Date)
				args = append(args, Grade)
				txid, err := app.Fabric.InvokeSchool(args)
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				fmt.Println(txid)
				fmt.Fprintf(w, txid)
				w.Header().Set("Content-Type", "application/json")
				//json.NewEncoder(w).Encode(txid)
			} else {
				http.Error(w, "Page not found or check the parameters", 500)
			}
		} else if MethodCall == "updatedetails" {

			RollNo := keys.Get("rollno")
			Name := keys.Get("name")
			Course := keys.Get("course")
			Details := keys.Get("details")
			Date := keys.Get("date")
			Grade := keys.Get("grade")
			if len(RollNo) > 0 && len(Name) > 0 && len(Course) > 0 && len(Details) > 0 && len(Date) > 0 && len(Grade) > 0 {

				var args []string
				args = append(args, "updateStudent")
				args = append(args, RollNo)
				args = append(args, Name)
				args = append(args, Course)
				args = append(args, Details)
				args = append(args, Date)
				args = append(args, Grade)
				txid, err := app.Fabric.UpdateStudent(args)
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				fmt.Println(txid)
				fmt.Fprintf(w, txid)
				w.Header().Set("Content-Type", "application/json")
				//json.NewEncoder(w).Encode(txid)
			}
		} else if MethodCall == "historydetails" {
			RollNo := keys.Get("rollno")
			Course := keys.Get("course")
			if len(RollNo) > 0 && len(Course) > 0 {
				var args []string
				args = append(args, "getHistoryForStudent")
				args = append(args, RollNo)
				args = append(args, Course)
				txid, err := app.Fabric.HistoryStudent(args)
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				fmt.Println(txid)
				fmt.Fprintf(w, txid)
				w.Header().Set("Content-Type", "application/json")
			}

		} else {
			http.Error(w, "Page not found or check the parameters", 500)
		}

	} else {
		http.Error(w, "Page not found", 500)
	}
}
