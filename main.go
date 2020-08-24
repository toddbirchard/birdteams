package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kib357/less-go"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Dynamic template values
type SiteMetaData struct {
	Title      string
	TagLine    string
	SiteUrl    string
	ShareImage string
	Background string
	Icon       string
}

// Compile and minify .LESS files
func CompileStylesheets() {
	staticFolder := "./static/styles/%s"
	err := less.RenderFile(
		fmt.Sprintf(staticFolder, "style.less"),
		fmt.Sprintf(staticFolder, "style.css"),
		map[string]interface{}{"compress": true})
	if err != nil {
		log.Fatal(err)
	}
}

// Render homepage template
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := SiteMetaData{
		Title:      "Bird Teams",
		TagLine:    "Let's go Bird Teams",
		SiteUrl:    "https://birdteams.org/",
		ShareImage: "/static/img/birdteams-share@2x.jpg",
		Background: "/static/img/background@2x.jpg",
		Icon:       "/static/img/favicon.png",
	}
	_ = tmpl.Execute(w, data)
}

// Route declaration
func Router() *mux.Router {
	staticDir := "/static/"
	// Page routes
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	return r
}

// Initiate web server
func main() {
	CompileStylesheets()
	GetTwitchToken()

	router := Router()
	client := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9300",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(client.ListenAndServe())
}
