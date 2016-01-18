package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

type Context struct {
	Title  string
	Static string
}

type time_zone struct {
	type2 string
	query string
}

type request struct {
	localtime string
	utcOffset string
}

type Result struct {
	request   request
	Time_zone time_zone
}

func Home(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	render(w, "index.html", context)
}

func About(w http.ResponseWriter, req *http.Request) {

	resp, err := http.Get("http://api.worldweatheronline.com/free/v2/tz.ashx?key=f787b94970cdef524e0f93891f2e5&q=New-York&format=xml")

	if err != nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//log.Printf("body %v", string(body))
	log.Printf(string(body))
	body2 := ""
	for i := 0; i < len(string(body)); i++ {

		if string(body[i]) == "e" && string(body[i+1]) == "t" {
			//log.Printf(string(body[i]))
			for j := i + 3; j < i+15; j++ {
				//log.Printf(string(body[j]))
				if string(body[j]) == "<" {

					break

				} else {
					body2 = body2 + string(body[j])

				}
			}
			break
		}

	}
	//s := new(Result)
	//xml.Unmarshal(body, &s)
	log.Printf(body2)

	x, err2 := strconv.ParseFloat(body2, 64)
	if err2 != nil {

	}
	//fmt.Printf("%f\n", x)
	//fmt.Printf("%d\n", int(x))
	ahora := time.Now().Hour() + int(x)
	//fmt.Printf("%d\n", ahora)
	ahora2 := strconv.Itoa(ahora)
	//fmt.Printf(ahora2)
	context := Context{Title: ahora2}
	render(w, "result.html", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/result", About)
	http.HandleFunc(STATIC_URL, StaticHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
