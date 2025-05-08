package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RandomUser struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func GetRandomUsers(count int) ([]RandomUser, error) {
	url := fmt.Sprintf("https://fakerapi.it/api/v1/users?_quantity=%d", count)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch random users: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Data []RandomUser `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return result.Data, nil
}

func GetRandomSentences(count int) ([]string, []string, error) {
	url := fmt.Sprintf("https://fakerapi.it/api/v1/texts?_quantity=%d&_characters=50", count)
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch random sentence: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Data []struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	titles := make([]string, len(result.Data))
	sentences := make([]string, len(result.Data))
	for i, item := range result.Data {
		titles[i] = item.Title
		sentences[i] = item.Content
	}

	if len(sentences) == 0 {
		return nil, nil, fmt.Errorf("no sentences returned from Faker API")
	}

	return titles, sentences, nil
}
