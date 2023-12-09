package todoist

import (
	"bytes"
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

func CreateTask(url string, token string, content string, projectId string) error {
	values := map[string]string{
		"content": content, 
		"project_id": projectId,
	}
	jsonData, err := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}