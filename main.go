package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hackersandslackers/birdteams/api"
	"github.com/kib357/less-go"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Dynamic template values
type SiteData struct {
	Title      string
	TagLine    string
	SiteUrl    string
	ShareImage string
	Background string
	Icon       string
	Videos     []api.YoutubeVideo
	Stream     bool
}

// Render homepage template
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := SiteData{
		Title:      "Bird Teams",
		TagLine:    "Let's go Bird Teams",
		SiteUrl:    "https://birdteams.org/",
		ShareImage: "/static/img/birdteams-share@2x.jpg",
		Background: "/static/img/background@2x.jpg",
		Icon:       "/static/img/favicon.png",
		Videos:     api.GetYoutubeVideos(),
		Stream:     api.GetTwitchStream(),
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
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

// Route declaration
func Router() *mux.Router {
	staticDir := "/static/"
	// Page routes
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	return r
}

// Initiate app
func main() {
	CompileStylesheets()
	api.LoadEnv()
	api.GetTwitchStream()

	router := Router()
	client := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9300",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server now live and listening at %s...", "127.0.0.1:9300")
	log.Fatal(client.ListenAndServe())
}
