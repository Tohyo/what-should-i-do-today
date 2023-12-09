package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Task struct {
	Id string `json:"id"`
	Content string `json:"content"`
	Url string `json:"url"`
}

func GetTodoistTasks(url string, token string) ([]Task, error) {
	var tasks []Task

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &tasks); err != nil {
		fmt.Println("Can not unmarshal JSON" + err.Error())
	}

	return tasks, nil
}