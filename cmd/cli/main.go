package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"

	// "fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // sql pool register
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"test.iamgak.net/models"
)

type application struct {
	infoLog         *log.Logger
	errorLog        *log.Logger
	db              *sql.DB
	models          *models.Init
	session         *sessions.CookieStore
	isAuthenticated bool
	user_id         int64
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// // Access the loaded environment variables
	sessionKey := os.Getenv("SESSION_KEY")
	// dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	var store = sessions.NewCookieStore([]byte(sessionKey))
	addr := flag.String("addr", ":"+dbPort, "HTTP network address")
	dsn := flag.String("dsn", fmt.Sprintf("%s:@/%s?parseTime=true", dbUser, dbName), "MySQL data source name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	ctx := context.Background()
	redis_name := "localhost"
	redis_password := ""
	redis_port := 6379
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redis_name, redis_port),
		Password: redis_password, // no password set
		DB:       0,              // use default DB
	})

	err = client.Set(ctx, "foo", "bar111", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("foo", val)
	// And add it to the application dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		db:       db,
		models:   models.Constructor(db, client),
		session:  store,
	}

	// app.SetSession()
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		TLSConfig: tlsConfig,
		// MaxHeaderBytes: 524288, // 0.5MB Max header size per request
		IdleTimeout:  time.Minute, // conncection close after 1 minute it do again handshake or something
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
