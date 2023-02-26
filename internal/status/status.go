package status

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/curusarn/resh/internal/httpclient"
	"github.com/curusarn/resh/internal/msg"
)

func get(port int) (*http.Response, error) {
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	client := httpclient.New()
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while GET'ing daemon /status: %w", err)
	}
	return resp, nil
}

func IsDaemonRunning(port int) (bool, error) {
	resp, err := get(port)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func GetDaemonStatus(port int) (*msg.StatusResponse, error) {
	resp, err := get(port)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsn, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading 'daemon /status' response: %w", err)
	}
	var msgResp msg.StatusResponse
	err = json.Unmarshal(jsn, &msgResp)
	if err != nil {
		return nil, fmt.Errorf("error while decoding 'daemon /status' response: %w", err)
	}
	return &msgResp, nil
}
