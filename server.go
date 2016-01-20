package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

type Context struct {
	Title  string
	City   string
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

	log.Printf(string(req.FormValue("city")))
	resp, err := http.Get("http://api.worldweatheronline.com/free/v2/tz.ashx?key=f787b94970cdef524e0f93891f2e5&q=" + req.FormValue("city") + "&format=xml")

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
	if body2 == "" {
		ahora2 := "Error parsing the city"
		context := Context{Title: ahora2, City: string(req.FormValue("city"))}

		render(w, "result.html", context)

	}

	//s := new(Result)
	//xml.Unmarshal(body, &s)
	log.Printf(body2)

	x, err2 := strconv.ParseFloat(body2, 64)
	if err2 != nil {

	}
	//fmt.Printf("%f\n", x)
	//fmt.Printf("%d\n", int(x))
	desfase := int(x)
	fmt.Printf("%d\n", desfase)

	//hora_Madrid = time.Now().Hour()
	//ahora_ah := time.Now().Date()

	// if start < ahora_ah {
	// 	hora_Madrid := -(desfase) + 6
	// } else {
	// 	hora_Madrid := desfase + 5
	// }
	ahora_3 := 0
	if desfase <= 0 {

		ahora_3 = desfase - 1
		ahora_3 = ahora_3 * -1

	} else {

		ahora_3 = 1 + desfase
		ahora_3 = ahora_3 * -1
	}

	fmt.Printf("%s\n", string(req.FormValue("hour")))
	you := strings.Replace(string(req.FormValue("hour")), ":", ".", -1)
	ahora, _ := strconv.ParseFloat(you, 64)
	ahora_4 := ahora + float64(ahora_3)
	fmt.Printf("%f\n", ahora)
	fmt.Printf("%d\n", ahora_3)
	fmt.Printf("%f\n", ahora_4)
	ahora5 := 0.00
	if ahora_4 >= 12 {

		ahora5 = ahora_4 - 12

	} else {

		if ahora_4 < 0 {
			ahora5 = ahora_4 + 12

		}
	}

	ahora2 := strconv.FormatFloat(ahora5, 'f', 2, 64)

	ahora2 = strings.Replace(ahora2, ".", ":", -1)
	fmt.Printf(ahora2)

	if ahora_4 > 24 {

		ahora2 = ahora2 + " am"

	} else {

		ahora2 = ahora2 + " pm"
	}

	context := Context{Title: ahora2, City: string(req.FormValue("city"))}
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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
