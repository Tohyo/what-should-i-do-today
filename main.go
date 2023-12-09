package main

import (
	"fmt"
	"log"
	"os"
	"whatShouldIDoToday/gitlab"
	"whatShouldIDoToday/todoist"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	todoistChan := make(chan []todoist.Task)
	gitlabChan := make(chan []gitlab.MR)

	go func () {
		defer close(todoistChan)

		tasks, err := todoist.GetTodoistTasks(
			"https://api.todoist.com/rest/v2/tasks", 
			os.Getenv("TODOIST_API_KEY"),
		)
		if err != nil {
			log.Fatal(err)
		}

		todoistChan <- tasks
	}()

	go func () {
		defer close(gitlabChan)

		mrs, err := gitlab.GetOpenedMRs(
			os.Getenv("GITLAB_URL") + "/api/v4/merge_requests?state=opened&scope=all&author_id=" + os.Getenv("GITLAB_USER_ID"),
			os.Getenv("GITLAB_PERSONAL_TOKEN"),
		)

		if err != nil {
			log.Fatal(err)
		}

		gitlabChan <- mrs
	}()

	var todoistTasks = <-todoistChan
	var gitlabMRs = <-gitlabChan

	for _, task := range todoistTasks {
		fmt.Println(task.Content)
	}

	for _, mr := range gitlabMRs {
		fmt.Println(mr.Title + ": " + mr.WebURL)
	}
}
