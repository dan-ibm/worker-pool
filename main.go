package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()
	// if workers count is more than 10, time will not decrease,
	// because every process of file saving is sleeping for a second
	const usersCount, workersCount = 100, 10
	var wg sync.WaitGroup
	defer wg.Wait()
	usersChan := make(chan User, usersCount)

	// worker pool
	for i := 0; i < workersCount; i++ {
		go func() {
			for uCh := range usersChan {
				saveUserInfo(uCh)
				wg.Done()
			}
		}()
	}

	generateUsers(usersCount, usersChan, &wg)

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(user User) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func generateUsers(usersCount int, usersChan chan<- User, wg *sync.WaitGroup) {
	for i := 0; i < usersCount; i++ {
		wg.Add(1)
		// filling of users channel
		usersChan <- User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}
	// closing channel to indicate stopping of value sending process
	close(usersChan)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
