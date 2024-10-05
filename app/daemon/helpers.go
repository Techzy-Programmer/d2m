package daemon

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/vars"
)

type ghMeta struct {
	Actions []string `json:"actions"`
}

func scheduleGHActionIPFetch() {
	tries := 0

	for {
		ips, err := getActionsIP()
		if err != nil {
			time.Sleep(time.Duration(tries+1) * 10 * time.Second) // Linear backoff

			if tries < 3 {
				tries++
				continue
			}

			// Network error
			paint.Error("Failed to fetch GitHub Actions IPs:", err)
		}

		tries = 0
		vars.GHActionIps = ips
		time.Sleep(6 * time.Hour)
	}
}

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

	var meta ghMeta
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return nil, err
	}

	return meta.Actions, nil
}
