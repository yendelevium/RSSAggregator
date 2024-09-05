package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Yendelevium/RSSAggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// This is kindof a weird way abt how go handles databases
	// We have to import a database drive in our program, but we don't have to call
	// anything from it. It's written in the sqlc docs. So yeah just, go get, go tidy and go vendor this url below
	// adding the _ as the alias name to indicate we won't be using it
	_ "github.com/lib/pq"
)

// This struct holds a connectionn to a database
// The database.Queries type was actually created by sqlc in the database folder
type apiConfig struct {
	DB *database.Queries
}

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

	// We are importing our db connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Error: DB_URL not found in the environment")
	}

	// The go standard library has a build in sql package
	// We connect to the database using sql.Open("driver name", connectionstring)
	// This returns a new connection and an error
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	// Lets create an apiConfig that we can pass into our handlers so that they have access to our database
	apiCfg := apiConfig{
		// This database.New() function is also made by sqlc
		// Just pass the connection to this function aand u get
		// a pointer to the database.Queries type. See the sqlc files to get more understanding
		// The DB is a database.Queries, but our connection is an sql.DB, so we r gonna convert it
		// to a database.Queries using the database.New() function
		DB: database.New(conn),
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

	// This is a POST request to /users, and we r calling the handlerCreateUse METHOD on apiCfg
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	// Hooking up a handler to get the user
	// To hookup the handlerGetUser, since its function signature is no longer a http.HandlerFunc
	// We have to call the middlewareAuth function, which will return a http.HandlerFunc so the
	// chi router can work with it
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	// While creating the feed, not only do u have to pass the name and url as http JSON Body,
	// U also have to pass the Authorization header, as u need that fr authing the user whos creating the feed
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	// This let's any user to get all of the feeds in our database
	// This is not an authenticated endpoint, so no need fr the Auth header, or to call the middleware func
	// As the function is already a http.HandlerFunc
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

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
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
