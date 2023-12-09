package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MR struct {
	Title string `json:"title"`
	WebURL string `json:"web_url"`
}

func GetMRs(url string, token string) ([]MR, error) {
	var mrs []MR

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &mrs); err != nil {
		fmt.Println("Can not unmarshal JSON" + err.Error())
	}

	return mrs, nil
}