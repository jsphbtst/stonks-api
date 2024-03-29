package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/db"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/middleware"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/routes"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/services"
)

const PORT int = 3000

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file %+v\n", err)
	}

	database := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
	uri := fmt.Sprintf("%s?authToken=%s", database, tursoToken)
	tursoDb, err := db.Init(uri)
	if err != nil {
		log.Fatalf("Failed to connect to Turso: %+v\n", err)
	}

	services.Init(tursoDb)
	log.Println("✅ Successfully connected to Turso")

	router := chi.NewRouter()
	router.Use(middleware.JsonContentTypeHeader)
	router.Use(middleware.RouteRuntimeLogger)

	router.Get("/", routes.Root)

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/companies/{symbol}", routes.GetCompanyById)
	})

	log.Printf("✅ Starting server in port %d\n", PORT)
	addr := fmt.Sprintf(":%d", PORT)
	server := &http.Server{Addr: addr, Handler: router}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start Go API server: ", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to exit ", err)
	}

	log.Println("Server exited")

}