package helpers

import (
	"encoding/json"
	"io"
	"net/http"
)

func getActionsIP() ([]string, error) {
	res, err := http.Get("https://api.github.com/meta")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var meta GitHubMeta
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return nil, err
	}

	return meta.Actions, nil
}
