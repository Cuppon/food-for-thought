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
)

const homePage string = "/index.html"
const staticPath string = "web/static"

func main() {
	// whatever directory you provide as a file server, effectively becomes root
	fs := http.FileServer(http.Dir(fmt.Sprintf("./%s/stylesheets", staticPath)))
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", fs))

	// TODO: update w/ appropriate config file information
	db, err := sqlx.Connect("postgres", "user=<fill in> dbname=food_for_thought sslmode=disable")
	if err != nil {
		fmt.Println("uhoh", db) // TODO: handle it!
	}

	var qu = `SELECT jsonb_build_object(
               'attribution', (SELECT jsonb_build_object(
                                  'id', r.attribution_source_id,
                                  'description', s.description,
                                  'location', s.location,
                                  'category', s.category
                               )
                               FROM source AS s
                               WHERE s.id = r.attribution_source_id
               ),
               'components', (SELECT jsonb_agg(jsonb_build_object(
                                'name', component,
                                'ingredient_specifications', ingredient_specifications
                                )) AS components
                               FROM (
                                   SELECT
                                       component,
                                       jsonb_agg(jsonb_build_object(
                                               'note', note,
                                               'component', component,
                                               'ingredient', jsonb_build_object(
                                                       'id', ingredient.id,
                                                       'native_name', ingredient.native_name,
                                                       'english_name', ingredient.english_name,
                                                       'shopping_link', ingredient.shopping_link,
                                                       'translated_name', ingredient.translated_name,
                                                       'english_category', ingredient.english_category
                                                             ),
                                               'amount_quantity', ins.amount_quantity,
                                               'amount_mass', ins.amount_mass,
                                               'preparation_quantity', ins.preparation_quantity,
                                               'preparation_type', ins.preparation_type,
                                               'preparation_length', ins.preparation_length
                                                 )) AS ingredient_specifications
                                   FROM ingredient_specification as ins
                                            JOIN ingredient ON ins.ingredient_id = ingredient.id
                                   GROUP BY component
                               ) AS comps
               ),
               'cuisine', (SELECT jsonb_build_object(
                              'id', r.cuisine_source_id,
                              'description', s.description,
                              'location', s.location,
                              'category', s.category
                           )
                           FROM source AS s
                           WHERE s.id = r.cuisine_source_id
               ),
               'emojis', (SELECT jsonb_agg(jsonb_build_object(
                            'id', rs.id,
                            'description', s.description,
                            'location', s.location,
                            'category', s.category
                          ))
                          FROM source AS s
                          INNER JOIN recipe_source AS rs ON s.id = rs.emoji_source_id
                          GROUP BY rs.recipe_id
               ),
               'instructions', r.instruction,
               'english_name', r.english_name,
               'native_name', r.native_name,
               'notes', r.note
        ) AS recipe
		FROM recipe r
		INNER JOIN ingredient_specification AS isp ON isp.recipe_id = r.id
		WHERE r.id = 1
		GROUP BY r.id;`

	var r recipes.Recipe
	row := db.QueryRow(qu)
	err = row.Scan(&r)
	if err != nil {
		fmt.Println("uhoh", err) // TODO: handle it!
	}

	http.HandleFunc("/", serveTemplate(r))
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
