package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type RSVP struct {
	Name, Email, Phone string
	WillAttend         bool
}

type FormData struct {
	*RSVP
	Errors []string
}

var responses = make([]*RSVP, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	// TODO: load template here
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded Template", index, name)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

func FormHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, FormData{
			RSVP: &RSVP{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := RSVP{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please Enter your Name")
		}

		if responseData.Email == "" {
			errors = append(errors, "Please Enter a Valid Email ID")
		}

		if responseData.Phone == "" {
			errors = append(errors, "Please Enter a Valid Phone Number")
		}

		if len(errors) > 0 {
			templates["form"].Execute(writer, FormData{
				RSVP: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func main() {
	//fmt.Println("TODO: Add Some Features")
	loadTemplates()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", FormHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
