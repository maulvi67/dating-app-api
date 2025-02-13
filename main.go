package main

import (
	"context"
	"dating-apps/app"
	"dating-apps/app/api/initialization"
	"dating-apps/app/api/middleware"
	"dating-apps/helper/config"
	"dating-apps/helper/logger"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	cfg := config.Init()
	log := logger.NewLogger(&cfg.ServerConfig.LogConfig)

	//Init DB Connection
	db, err := initialization.InitDatabase(&cfg.DBConfig)
	if err != nil {
		panic(err.Error())
	}

	// pass to infra value
	infra := &app.Infra{Db: &db, Log: log, Config: cfg}

	// Routing initialization
	mux := initialization.InitRouting(infra)
	address := *flag.String("listen", ":"+strconv.Itoa(cfg.ServerConfig.Port), "Listen address.")
	httpServer := http.Server{
		Addr:    address,
		Handler: middleware.ServeHTTP(mux, infra.Log),
	}

	// Setup graceful shutdown
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Info("start graceful shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			panic(err)
		}
		close(idleConnectionsClosed)
	}()

	log.Info(fmt.Sprintf("Listening at port %s", strconv.Itoa(cfg.ServerConfig.Port)))
	if err = httpServer.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	<-idleConnectionsClosed
}
