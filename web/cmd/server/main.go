package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/skale-5/skalogram/web"
	"github.com/skale-5/skalogram/web/delivery/http"
	"github.com/skale-5/skalogram/web/pkg/gcs"
	"github.com/skale-5/skalogram/web/pkg/postgresql/post"
	"github.com/skale-5/skalogram/web/pkg/redis"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "skalogram"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	postCacheService := web.NewPostCacheService(
		redis.NewClient("127.0.0.1:6379"),
	)

	postDatabaseService := web.NewPostDatabaseService(
		post.New(db),
	)

	postStorageService := web.NewPostStorageService(
		gcs.NewClient(context.Background()),
	)

	server := http.NewServer(http.NewServerArgs{
		ListenAddr:          "0.0.0.0:8080",
		PostDatabaseService: postDatabaseService,
		PostCacheService:    postCacheService,
		PostStorageService:  postStorageService,
	})
	server.Run()
}
