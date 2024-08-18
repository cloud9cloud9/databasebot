package main

import (
	wp "db/workerpool"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	users := make([]wp.User, 20)
	generateUsers(users)
	userChan := make(chan wp.User, len(users))

	workerPool := wp.NewWorkerPool(4, users, userChan)
	workerPool.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func generateUsers(users []wp.User) {
	for i := 0; i < 20; i++ {
		users[i] = wp.User{
			Name:  "name" + strconv.Itoa(i),
			Email: "email" + strconv.Itoa(i) + "@gmail.com",
		}
	}
}
