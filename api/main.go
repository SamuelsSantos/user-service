package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/samuelssantos/user-service/domain/entity/user"
	"github.com/samuelssantos/user-service/pkg/password"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/samuelssantos/user-service/api/handler"
	"github.com/samuelssantos/user-service/api/middleware"
	"github.com/samuelssantos/user-service/config"
	"github.com/samuelssantos/user-service/pkg/metric"
)

func main() {

	cfg := config.NewConfig()

	db, err := sql.Open(cfg.Db.Driver, cfg.Db.ToURL())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	userRepo := user.NewSQLRepoRepository(db)
	userManager := user.NewManager(userRepo, password.NewService())

	metricService, err := metric.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.HandlerFunc(middleware.Metrics(metricService)),
		negroni.NewLogger(),
	)

	//user
	handler.MakeUserHandlers(r, *n, userManager)

	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + cfg.Server.Port,
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
