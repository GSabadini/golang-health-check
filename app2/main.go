package main

import (
	"app2/health"
	"app2/health/check"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	//Connect database
	_, err := newMySQLHandler()
	if err != nil {
		log.Println(err)
	}

	var (
		//Connect cache redis
		_ = newClientRedis()

		//Gorilla router
		router = mux.NewRouter()
	)

	router.HandleFunc("/monitors", monitors())
	log.Println("Start HTTP server :3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
		panic(err)
	}
}

func monitors() http.HandlerFunc {
	redis1 := health.Register(health.Config{
		Name:  "Redis",
		Port:  "6379",
		Host:  "redis",
		Check: check.Redis{DNS: "redis:6379"}.Check,
	})

	mysql := health.Register(health.Config{
		Name:  "MySQL",
		Port:  "3306",
		Host:  "mysql",
		Check: check.MySQL{DNS: "root:dev@tcp(mysql:3306)/app"}.Check,
	})

	app1 := health.Register(health.Config{
		Name:  "App1",
		Port:  "3001",
		Host:  "app1",
		Check: check.WebServer{URL: "http://app1:3001/health"}.Check,
	})

	app3 := health.Register(health.Config{
		Name: "App3",
		Port: "3003",
		Host: "app3",
		Check: func() error {
			return errors.New("FAKE")
		},
	})

	h, err := health.New(redis1, mysql, app1, app3)
	if err != nil {
		log.Fatal(err)
	}

	return h.HandlerFunc
}

func newClientRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println("connected redis", pong, err)

	return client
}

func newMySQLHandler() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:dev@tcp(mysql:3306)/app")
	if err != nil {
		return &sql.DB{}, err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS stats (name VARCHAR(50) PRIMARY KEY, count INTEGER);"); err != nil {
		return &sql.DB{}, err
	}

	log.Println("Database connected")

	return db, nil
}
