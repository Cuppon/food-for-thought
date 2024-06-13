package main

import (
	"fmt"
	"github.com/Cuppon/foodpls/recipes"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const homePage string = "/index.html"
const staticPath string = "web/static"

func main() {
	// TODO: update w/ appropriate config file information
	db, err := sqlx.Connect("postgres", "user=<fill in> dbname=food_for_thought sslmode=disable")
	if err != nil {
		panic(err) // TODO: handle it!
	}

	appConf := recipes.Config{
		Storage: &recipes.PG{Conn: db},
	}
	fmt.Println(appConf) // TODO: actually use this via injection to endpoint to set next recipe

	scheduleConf := &recipes.ScheduleConfig{
		TickerDuration: time.Hour,         // TODO: pull this from config file
		DailyRecipe:    &recipes.Recipe{}, // TODO: update with an actual recipe
		NextRecipe:     &recipes.Recipe{}, // TODO: to be updated via endpoint
	}

	go func() {
		scheduleConf.ScheduleDailyRecipe()
	}()

	// whatever directory you provide as a file server, effectively becomes root
	fs := http.FileServer(http.Dir(fmt.Sprintf("./%s/stylesheets", staticPath)))
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", fs))
	http.HandleFunc("/", serveTemplate(*scheduleConf.DailyRecipe)) // TODO: double check pointer semantics
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err) // TODO: hook up proper logging
	}
}

func serveTemplate(recipe recipes.Recipe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if filepath.Clean(r.URL.Path) != homePage {
			log.Println("Attempted access to non-home page: " + r.URL.Path) // TODO: log where people are making requests to
			http.NotFound(w, r)
			return
		}

		// layout.html is applied to every web page: it's a template for the overall structure of HTML pages
		lp := filepath.Join(staticPath, "layout.html")
		// path to the specific file the user is requesting
		fp := filepath.Join(staticPath, filepath.Clean(r.URL.Path))

		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
		}

		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp, fp)
		if err != nil {
			log.Println("Could not parse template file: ", err.Error()) // TODO: hook up proper logging
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "layout", recipe)
		if err != nil {
			log.Println("Could not execute template: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
