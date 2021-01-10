package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bhongy/rediret-api-clean-architecture-golang/api"
	mr "github.com/bhongy/rediret-api-clean-architecture-golang/repository/mongodb"
	rr "github.com/bhongy/rediret-api-clean-architecture-golang/repository/redis"
	"github.com/bhongy/rediret-api-clean-architecture-golang/shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// repo <- service -> serializer -> http

func main() {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		port := httpPort()
		fmt.Printf("Listening on port %s\n", port)
		errs <- http.ListenAndServe(port, r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		log.Printf("Creating a new redis repo (redisURL=%q)...", redisURL)
		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully created redis repo")
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		log.Printf("Creating a new mongo repo (mongoURL=%q)...", mongoURL)
		repo, err := mr.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully created mongo repo")
		return repo
	}
	return nil
}
