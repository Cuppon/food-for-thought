package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Cuppon/foodpls/recipes"
	"github.com/Cuppon/foodpls/webserver"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// TODO: update w/ appropriate config file information
	db, err := sqlx.Connect("postgres", "user=<fill in> dbname=food_for_thought sslmode=disable")
	if err != nil {
		panic(err) // TODO: handle it!
	}

	appConf := recipes.AppConfig{
		Storage: &recipes.PG{Conn: db},
	}
	fmt.Println(appConf) // TODO: actually use this via injection to endpoint to set next recipe

	recipeConf := &recipes.RecipeConfig{
		TickerDuration: time.Hour,        // TODO: pull this from config file
		DailyRecipe:    recipes.Recipe{}, // TODO: update with an actual recipe
		NextRecipe:     recipes.Recipe{}, // TODO: to be updated via endpoint
	}

	go func() {
		recipeConf.ScheduleDailyRecipe(appConf)
	}()

	// TODO: pull this from config file
	templateConf := webserver.TemplateConfig{
		HomePage:   "/index.html",
		StaticPath: "web/static",
	}

	srv := webserver.NewServer(
		webserver.StaticFilesHandler(templateConf.StaticPath),
		webserver.TemplateHandler(templateConf, recipeConf),
		webserver.SetNextRecipeHandler(recipeConf),
	)
	httpsServer := &http.Server{
		Handler: srv,
	}
	if err = httpsServer.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatal(err) // TODO: hook up proper logging
	}
}
