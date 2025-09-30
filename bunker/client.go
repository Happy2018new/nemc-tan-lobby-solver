package bunker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	AuthServer string
	FBToken    string
}

func NewClient(authServer string, token string) *Client {
	return &Client{
		AuthServer: authServer,
		FBToken:    token,
	}
}

func parseHttpResponse[T any](resp *http.Response) (result T, err error) {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("parseHttpResponse: The status code of http request is not 200 (code = %d)", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("parseHttpResponse: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return result, fmt.Errorf("parseHttpResponse: %v", err)
	}

	return result, nil
}
