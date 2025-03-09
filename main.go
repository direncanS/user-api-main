package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lambda-lama/user-api/handlers"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

const (
	localServerPort = ":8080"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Not found", r.RequestURI)
		http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	})

	r.HandleFunc("/user/images", handlers.GetDataTopics).Methods("GET")
	r.HandleFunc("/user/images/{topic}", handlers.GetByTopic).Methods("GET")
	r.HandleFunc("/user/videos", handlers.GetVideosByTopic).Methods("GET")
	r.HandleFunc("/user/videos/{folder}", handlers.GetVideoFromFolder).Methods("GET")

	return r
}

func setupHttpRouter() *http.Server {
	return &http.Server{
		Addr:    localServerPort,
		Handler: setupRouter(),
	}
}

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	fmt.Println("Started")
	if runtime_api, _ := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); runtime_api != "" {
		log.Println("Starting up in Lambda Runtime")
		adapter := gorillamux.NewV2(setupRouter())
		lambda.Start(adapter.ProxyWithContext)
	} else {
		log.Println("Starting up on local")
		srv := setupHttpRouter()
		_ = srv.ListenAndServe()
	}
}
