package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"deftask/internal/repo"
	"deftask/internal/server"
	"deftask/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var address = flag.String("address", "localhost:3000", "address of server")
var dbDSN = flag.String(
	"db", "host=localhost port=6432 user=postgres password=password database=be",
	"address of db",
)

func main() {
	flag.Parse()

	ctx := context.Background()

	appStopSignal := make(chan os.Signal, 1)
	signal.Notify(appStopSignal, os.Interrupt, syscall.SIGTERM)

	dbConn, err := sqlx.Connect("pgx", *dbDSN)
	if err != nil {
		panic(err)
	}
	fmt.Println("db connected")

	rp := repo.New(dbConn)
	svc := service.New(rp)
	srv := server.New(*address, svc)

	go srv.Run()

	<-appStopSignal

	srv.Shutdown(ctx)
}
