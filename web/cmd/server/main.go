package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skale-5/skalogram/web"
	"github.com/skale-5/skalogram/web/config"
	"github.com/skale-5/skalogram/web/delivery/http"
	"github.com/skale-5/skalogram/web/pkg/gcs"
	"github.com/skale-5/skalogram/web/pkg/s3"

	"github.com/skale-5/skalogram/web/pkg/postgresql/post"
	"github.com/skale-5/skalogram/web/pkg/redis"

	_ "github.com/lib/pq"
)

func main() {

	isPrintDefaults := flag.Bool("print-defaults", false, "Print default configurations")
	flag.Parse()

	if *isPrintDefaults {
		config.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()

	// DATABASE SERVICE
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Env().Get("PG_HOST"),
		config.Env().Get("PG_PORT"),
		config.Env().Get("PG_USER"),
		config.Env().Get("PG_PASSWORD"),
		config.Env().Get("PG_DBNAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	postDatabaseService := web.NewPostDatabaseService(
		post.New(db),
	)

	// CACHE SERVICE
	redisConn := fmt.Sprintf("%s:%s",
		config.Env().Get("REDIS_HOST"),
		config.Env().Get("REDIS_PORT"),
	)
	postCacheService := web.NewPostCacheService(
		redis.NewClient(redisConn),
	)

	// STORAGE SERVICE
	var postStorageService *web.PostStorageService

	switch config.Env().Get("STORAGE_TYPE") {
	case "gs":
		postStorageService = web.NewPostStorageService(
			gcs.NewClient(ctx),
		)
	case "s3":
		postStorageService = web.NewPostStorageService(
			s3.NewClient(ctx, config.Env().Get("STORAGE_BUCKET_REGION")),
		)
	default:
		log.Fatal("no storage type configured")
	}

	listenAddr := fmt.Sprintf("%s:%s",
		config.Env().Get("LISTEN_ADDR"),
		config.Env().Get("LISTEN_PORT"),
	)
	server := http.NewServer(http.NewServerArgs{
		ListenAddr:          listenAddr,
		PostDatabaseService: postDatabaseService,
		PostCacheService:    postCacheService,
		PostStorageService:  postStorageService,
	})
	server.Run()
}
