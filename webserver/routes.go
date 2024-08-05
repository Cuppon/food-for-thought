package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cuppon/foodpls/recipes"
)

func RedirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func SetNextRecipeHandler(recipeConfig *recipes.RecipeConfig) Route {
	return func(mux *http.ServeMux) {
		endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var nextRecipe recipes.Recipe
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&nextRecipe); err != nil {
				panic(err) // TODO: handle it!
			}

			recipeConfig.SetNextRecipe(nextRecipe)
		})

		endpoint = AddMiddleware(endpoint, AuthorizeMiddleware, ValidateJSONMiddleware)

		mux.Handle("/update-next-recipe", endpoint)
	}
}

func StaticFilesHandler(staticFilesPath string) Route {
	return func(mux *http.ServeMux) {
		// whatever directory you provide as a file server, effectively becomes root
		// TODO: make this configurable
		fs := http.FileServer(http.Dir(fmt.Sprintf("./%s/serveable", staticFilesPath)))
		mux.Handle("/serveable/", serveStatic(staticFilesPath, fs))
	}
}

func serveStatic(staticFilesPath string, fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		if !strings.HasPrefix(path, "/serveable/") {
			http.NotFound(w, r)
			return
		}

		fp := filepath.Join(staticFilesPath, path)
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

		http.StripPrefix("/serveable/", fileServer).ServeHTTP(w, r)
	}
}

type TemplateConfig struct {
	HomePage   string
	StaticPath string
}

func TemplateHandler(templateConfig TemplateConfig, recipeConfig *recipes.RecipeConfig) Route {
	return func(mux *http.ServeMux) {
		// TODO: this will eventually be replaced with cache-related config, and serveTemplate will pull from a cache
		mux.HandleFunc("/", serveTemplate(templateConfig, recipeConfig))
	}
}

func serveTemplate(tc TemplateConfig, rc *recipes.RecipeConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: put back in after done testing
		//if filepath.Clean(r.URL.Path) != tc.HomePage {
		//	log.Println("Attempted access to non-home page: " + r.URL.Path) // TODO: log where people are making requests to
		//	http.NotFound(w, r)
		//	return
		//}

		// layout.html is applied to every web page: it's a template for the overall structure of HTML pages
		lp := filepath.Join(tc.StaticPath, "layout.html")
		// path to the specific file the user is requesting
		fp := filepath.Join(tc.StaticPath, filepath.Clean(r.URL.Path)+".html")

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

		// TODO: cache the parsed templates
		tmpl, err := template.ParseFiles(lp, fp)
		if err != nil {
			log.Println("Could not parse template file: ", err.Error()) // TODO: hook up proper logging
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "layout", rc.DailyRecipe)
		if err != nil {
			log.Println("Could not execute template: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
