package main

import (
	"crud/config"
	"crud/controllers"
	"crud/env"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/pdrum/swagger-automation/api"

	_ "github.com/pdrum/swagger-automation/docs"
)

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
	})
}

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
	port := env.GetValue("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	err := config.InitMongoDB()
	if err != nil {
		log.Fatal("error while initializing mongoDB client")
		panic(err)
	}
	config.InitRedisConfig()
	log.Println("Starting the HTTP server on port 8090")

	router := mux.NewRouter().StrictSlash(true)
	initaliseHandlers(router)
	log.Fatal(http.ListenAndServe(":"+port, RequestLogger(router)))
}

func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/getcovidcases", controllers.GetCovidCasesByLocation).Methods("GET")
}
