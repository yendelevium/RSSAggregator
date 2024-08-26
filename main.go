package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// By default, this loads the .env file, but we can also specify
	// This makes it so that everytime we start a session, it pulls the env variable form our .env file into the environment
	// so we don't need to explicitly specify "export PORT=8080" in the cmd or for that case ANY env variable
	godotenv.Load(".env")

	// using this, we can get the value of a variable using a key
	// Here the key is PORT. Getenv means get-environment variable
	// fmt.Println(os.Environ()) This will print ALL environmental variables in key-value pairs
	portString := os.Getenv("PORT")
	if portString == "" {
		// log.Fatal() will exit the program immediately with Error code 1 and a message
		// This will get us an error, if the port isn't there in the envi
		// w/o the package, we will have to use "export PORT=8080" in the cmd everytime we start up a new cmd
		// But with the package, no need to manually set the env varian
		log.Fatal("Error: PORT not found in the environment")
	}
	// This creates a new router object
	router := chi.NewRouter()

	// Adding a cors configuration to our router. This is so that users can make requests from their web browser
	// cors are basically telling our server to send some extra HTTP header in our responses to tell browsers we allow
	// requests frm these these origins, and methods and headers etc etc. This configuration I've written is basically very UNrestrictive
	// We can make it a lot tighter if we wantfor security reasons, but this is fine for now
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Let's connect the HTTP handler we made to a sub-router. We r creating a new sub-router, and adding the handler to that
	v1Router := chi.NewRouter()
	// the .Get(path, handler function) will execute the handlerReadiness function when a GET Req is sent to that path
	// /healthz is a standard path to check if ur server is up and running, and that's the purpose
	// of the handlerReadiness, to check if the server is alive and running
	v1Router.Get("/healthz", handlerReadiness)
	// v1Router.HandleFunc("/healthz", handlerReadiness) is the same as the abv function, but it allows ANY request, not only get

	// Creating am error endpoint so the ppl using the api know what an error will look like
	v1Router.Get("/err", handlerErr)

	// The reason we made a new router, is coz we r gonna mount that to our original router
	// We r nesting a v1 r path will be localhost:8080/v1/healthz
	// Nesting subrouters like this is actually very common practice in web-development, as its very useful
	// Nesting subrouters helps to organize the routing logic into smaller, modular components.
	// Each subrouter can handle a specific feature or resource (e.g., users, products, orders) and can be developed, tested, and maintained independently.
	// Also, as the application grows, managing everything in a single subrouter is annoying, so this is chill
	// This also allows us to reuse these subrouters in bigger applications
	// Subrouters allow you to apply specific middleware to a group of routes.
	router.Mount("/v1", v1Router)

	// Lets connect the router to a HTTP server
	srv := &http.Server{
		// It needs a handler, which is the router
		Handler: router,
		// And an address, which is just :8080
		Addr: ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	// ListenAndServe() blocks. It's now listening to http requests on the server
	// If anything goes wrong, it returns an error, and we will log that error and exit that program
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
