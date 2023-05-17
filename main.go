package main

import (
	"fmt"
	"os"
	"os/signal"
	"shortner/data"
	"shortner/domain"
	"shortner/presentation"
	"strconv"
	"syscall"

	"github.com/go-redis/redis/v8"
)

var (
	redisAddr     = os.Getenv("REDIS_ADDR")
	redisUsername = os.Getenv("REDIS_USERNAME")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB       = os.Getenv("DB")
)

func main() {
	rdb, err := strconv.Atoi(redisDB)
	if err != nil {
		rdb = 0
	}
	redisOpt := redis.Options{
		Addr:     redisAddr,
		Username: redisUsername,
		Password: redisPassword,
		DB:       rdb,
	}
	store := data.NewRedisData(redisOpt)
	service := domain.NewService(store)
	sigs := make(chan os.Signal, 1)
	errChan := make(chan int, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		presentation.NewAdmin(service).Start(":7000")
		errChan <- 0
	}()
	go func() {
		presentation.NewHttpRedirect(service).Start(":8000")
		errChan <- 0
	}()
	for {
		select {
		case <-sigs:
			fmt.Println("shouting down the application")
			os.Exit(0)

		case <-errChan:
			fmt.Println("there is an error in application")
			os.Exit(1)
		}

	}

}
