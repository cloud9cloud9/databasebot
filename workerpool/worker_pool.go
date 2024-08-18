package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type WorkerPool struct {
	Users       []User
	Concurrency int
	userChan    chan User
	wg          *sync.WaitGroup
	F           *Folder
}

func (t *User) processWriteToFile(folder *Folder) {
	fmt.Printf("processing task with name %s and email %s\n", t.Name, t.Email)

	err := folder.WriteUserToFile(*t)
	if err != nil {
		fmt.Printf("error writing user to file: %v\n", err)
	}

	time.Sleep(time.Second * 2)
}

func (wp *WorkerPool) worker() {
	for user := range wp.userChan {
		task := User{Name: user.Name, Email: user.Email}
		task.processWriteToFile(wp.F)
		wp.wg.Done()
	}
}

func (wp *WorkerPool) processReadFromFile() {
	for {
		users, err := wp.F.ReadFromFile()
		if err != nil {
			fmt.Printf("error reading user from file: %v\n", err)
		} else {
			fmt.Println("Users read from file:")
			for _, user := range users {
				fmt.Printf("Name: %s, Email: %s\n", user.Name, user.Email)
			}
		}
		time.Sleep(time.Second * 6)
	}
}

func (wp *WorkerPool) Run() {
	wp.userChan = make(chan User, len(wp.Users))
	for i := 0; i < wp.Concurrency; i++ {
		go wp.worker()
	}

	go wp.processReadFromFile()
	wp.wg.Add(len(wp.Users))
	for _, users := range wp.Users {
		wp.userChan <- users
	}
	close(wp.userChan)
	wp.wg.Wait()
}

func NewWorkerPool(concurrency int, users []User, uc chan User) *WorkerPool {
	return &WorkerPool{
		Users:       users,
		Concurrency: concurrency,
		userChan:    uc,
		wg:          &sync.WaitGroup{},
		F:           &Folder{&sync.RWMutex{}, "users.json"},
	}
}
