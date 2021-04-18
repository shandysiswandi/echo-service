package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/shandysiswandi/echo-service/internal/config"
	"github.com/shandysiswandi/echo-service/internal/infrastructure/app"
	"github.com/shandysiswandi/echo-service/internal/infrastructure/mongodb"
)

func main() {
	/********** ********** ********** ********** **********/
	/* load .env
	/********** ********** ********** ********** **********/
	if err := godotenv.Load(); err != nil {
		println("error loading .env file")
	}

	/********** ********** ********** ********** **********/
	/* create config instance & context
	/********** ********** ********** ********** **********/
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cfg := config.New()

	/********** ********** ********** ********** **********/
	/* create mongo db & connect
	/********** ********** ********** ********** **********/
	mongoDB := mongodb.New(cfg, ctx)
	if err := mongoDB.Connect(); err != nil {
		println("error connect mongo db")
	}

	/********** ********** ********** ********** **********/
	/* call freamework echo engine
	/********** ********** ********** ********** **********/
	e := app.Injection(cfg)

	/********** ********** ********** ********** **********/
	/* run server on goroutine
	/********** ********** ********** ********** **********/
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			println(err.Error())
			os.Exit(1)
		}
	}()

	/********** ********** ********** ********** **********/
	/* make channel to receive signal
	/********** ********** ********** ********** **********/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	/********** ********** ********** ********** **********/
	/* defer resources
	/********** ********** ********** ********** **********/
	defer mongoDB.Close() // close mongo connectionx
	defer cancel()        // main context background

	/********** ********** ********** ********** **********/
	/* shutdown server
	/********** ********** ********** ********** **********/
	println("💥 shutdown server ...")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
