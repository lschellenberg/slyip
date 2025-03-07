package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"yip/src/api"
	"yip/src/app"
	"yip/src/config"
)

func main() {
	fmt.Println("starting service")
	c, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(c)

	a, err := app.InitApp(&c)

	if err != nil {
		log.Fatalln(err)
	}

	startService(a)
}

func startService(app *app.App) {
	a := api.NewApi(app)
	showRoutes(&a)
	log.Println("starting service at port:", app.Config.API.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", app.Config.API.Port), a.Router))
}

func showRoutes(api *api.Api) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(api.Router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}
