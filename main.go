package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var host = "postgres"

const (
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"

	closeTimeout = 5 * time.Second
)

type App struct {
	db     *DB
	server *http.Server
}

func (a *App) Run(port string) error {
	a.server.Addr = port

	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}

	log.Printf("Server started at %s", port)
	return nil
}

func (a *App) Close() error {
	if err := a.db.db.Close(); err != nil {
		return err
	}
	log.Println("Database connection closed")

	tc, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	return a.server.Shutdown(tc)
}

func init() {
	if os.Getenv("POSTGRES_HOST") != "" {
		host = os.Getenv("POSTGRES_HOST")
	}
}

func main() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
	database, err := Open(dsn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection successfully established!")

	r := mux.NewRouter()
	app := &App{
		db:     database,
		server: &http.Server{Handler: r},
	}

	api := r.PathPrefix("/api").Subrouter()
	api.Use(corsMiddleware)
	api.HandleFunc("/coupons", app.getCouponsHandler).Methods(http.MethodGet)
	api.HandleFunc("/coupons", app.addCouponHandler).Methods(http.MethodPost)
	api.HandleFunc("/{id:[0-9]+}/redeem", app.redeemCouponHandler).Methods(http.MethodPost)
	api.HandleFunc("/{id:[0-9]+}", app.deleteCouponHandler).Methods(http.MethodDelete)

	go app.Run(":8080")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Received signal %s, shutting down...", s.String())

	app.Close()
	log.Println("Application gracefully stopped")
}
