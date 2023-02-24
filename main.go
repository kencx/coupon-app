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

var hostDB = "postgres"
var hostRedis = "redis"
var port = "8080"

const (
	portDB   = 5432
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
	return nil
}

func (a *App) Close() error {
	if err := a.db.db.Close(); err != nil {
		return err
	}
	log.Println("Database connection closed")

	if err := a.db.redis.Close(); err != nil {
		return err
	}
	log.Println("Redis connection closed")

	tc, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	return a.server.Shutdown(tc)
}

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	if os.Getenv("POSTGRES_HOST") != "" {
		hostDB = os.Getenv("POSTGRES_HOST")
	}
	if os.Getenv("REDIS_HOST") != "" {
		hostRedis = os.Getenv("REDIS_HOST")
	}
}

func main() {
	dbDSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, hostDB, portDB, dbname)
	cacheDSN := fmt.Sprintf("%s:6379", hostRedis)

	database, err := Open(dbDSN, cacheDSN)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection successfully established!")
	log.Println("Redis connection successfully established")

	r := mux.NewRouter()
	app := &App{
		db:     database,
		server: &http.Server{Handler: r},
	}

	api := r.PathPrefix("/api").Subrouter()
	api.Use(corsMiddleware)
	api.HandleFunc("/coupon/{id:[0-9]+}", app.getCouponHandler).Methods(http.MethodGet)
	api.HandleFunc("/coupons", app.getCouponsHandler).Methods(http.MethodGet)
	api.HandleFunc("/coupons", app.addCouponHandler).Methods(http.MethodPost)
	api.HandleFunc("/{id:[0-9]+}/redeem", app.redeemCouponHandler).Methods(http.MethodPost)
	api.HandleFunc("/{id:[0-9]+}", app.deleteCouponHandler).Methods(http.MethodDelete)

	go app.Run(fmt.Sprintf(":%s", port))
	log.Printf("Server started at %s", port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Received signal %s, shutting down...", s.String())

	app.Close()
	log.Println("Application gracefully stopped")
}
