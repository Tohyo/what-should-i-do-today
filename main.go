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
	gitlabMRChan := make(chan []gitlab.MR)
	gitlabReviewChan := make(chan []gitlab.MR)

	go func () {
		defer close(todoistChan)

		tasks, err := todoist.GetTodoistTasks(
			"https://api.todoist.com/rest/v2/tasks?project_id=" + os.Getenv("TODOIST_PROJECT_ID"), 
			os.Getenv("TODOIST_API_KEY"),
		)
		if err != nil {
			log.Fatal(err)
		}

		todoistChan <- tasks
	}()

	go func () {
		defer close(gitlabMRChan)

		mrs, err := gitlab.GetMRs(
			os.Getenv("GITLAB_URL") + "/api/v4/merge_requests?state=opened&scope=all&author_id=" + os.Getenv("GITLAB_USER_ID"),
			os.Getenv("GITLAB_PERSONAL_TOKEN"),
		)

		if err != nil {
			log.Fatal(err)
		}

		gitlabMRChan <- mrs
	}()

	go func () {
		defer close(gitlabReviewChan)

		mrs, err := gitlab.GetMRs(
			os.Getenv("GITLAB_URL") + "/api/v4/merge_requests?state=opened&scope=all&reviewer_id=" + os.Getenv("GITLAB_USER_ID"),
			os.Getenv("GITLAB_PERSONAL_TOKEN"),
		)

		if err != nil {
			log.Fatal(err)
		}

		gitlabReviewChan <- mrs
	}()

	fmt.Println("*******************************")
	fmt.Println("*** What should I do today? ***")
	fmt.Println("*******************************")
	fmt.Println("")

	fmt.Println("*** Todoist tasks ***")
	for _, task := range <-todoistChan {
		fmt.Println("# " + task.Content)
	}
	fmt.Println("*********************")
	fmt.Println("")

	fmt.Println("*** Gitlab MRs ***")
	for _, mr := range <-gitlabMRChan {
		fmt.Println(mr.Title + ": " + mr.WebURL)
	}
	fmt.Println("******************")
	fmt.Println("")

	fmt.Println("*** Gitlab Reviews ***")
	for _, mr := range <-gitlabReviewChan {
		fmt.Println(mr.Title + ": " + mr.WebURL)
	}
	fmt.Println("***********************")

}
