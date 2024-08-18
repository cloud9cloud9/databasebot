package workerpool

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Folder struct {
	Mu   *sync.RWMutex
	Path string
}

func (f *Folder) WriteUserToFile(user User) error {
	f.Mu.Lock()
	defer f.Mu.Unlock()

	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func (f *Folder) ReadFromFile() ([]User, error) {
	var users []User

	file, err := os.Open(f.Path)
	if err != nil {
		return users, fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var user User
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			return users, fmt.Errorf("failed to unmarshal user: %v", err)
		}
		users = append(users, user)
	}

	if err := scanner.Err(); err != nil {
		return users, fmt.Errorf("error reading file: %v", err)
	}

	return users, nil
}

//func (f *Folder) GetUserByEmail(email string) (User, error) {
//	users, err := f.ReadFromFile()
//	if err != nil {
//		return User{}, err
//	}
//	for _, user := range users {
//		if user.Email == email {
//			return user, nil
//		}
//	}
//	return User{}, fmt.Errorf("user with email %s not found", email)
//}
