package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SealTV/govue/server"

	"github.com/SealTV/govue/repository"

	_ "github.com/lib/pq"
)

var (
	host = flag.String("host", "127.0.0.1", "Web server host")
	port = flag.Int("port", 8080, "Web server port")

	dbHost = flag.String("db_host", "localhost", "DB host")
	dbPort = flag.Int("db_port", 5432, "DB port")
	dbName = flag.String("db_name", "gobue_db", "DB name")
	dbUser = flag.String("db_user", "root", "DB user")
	dbPass = flag.String("db_pass", "1234", "DB password")
	dbSsl  = flag.String("db_ssl", "disable", "DB use ssl connnection (disable | enable)")
)

func main() {
	flag.Parse()

	connecionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		*dbHost, *dbPort, *dbUser, *dbPass, *dbName, *dbSsl)
	db, err := sql.Open("postgres", connecionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	go func() {
		cancel := make(chan os.Signal)
		signal.Notify(cancel, syscall.SIGTERM)
		signal.Notify(cancel, syscall.SIGINT)
		sig := <-cancel
		log.Printf("Stop signal: %v", sig)
		os.Exit(0)
	}()

	ar := repository.NewAccountRepository(db)
	webHanler := server.NewHTTPHandler(ar)
	log.Printf("Start listen on: %s", fmt.Sprintf("%s:%d", *host, *port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), webHanler); err != nil {
		log.Fatal(err)
	}
}
