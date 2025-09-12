package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type secretLoadingTransport struct {
	secret string
}

func (s secretLoadingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.secret))
	return http.DefaultTransport.RoundTrip(req)
}

type ClientOptions struct {
	AuthServer          string
	RespondUserOverride string
}

func MakeDefaultClientOptions() *ClientOptions {
	return &ClientOptions{}
}

type ClientInfo struct {
	GrowthLevel int
	RespondTo   string
}

type Client struct {
	ClientInfo
	client http.Client
	*ClientOptions
}

func parsedError(message string) error {
	error_regex := regexp.MustCompile("^(\\d{3} [a-zA-Z ]+)\n\n(.*?)($|\n)")
	err_matches := error_regex.FindAllStringSubmatch(message, 1)
	if len(err_matches) == 0 {
		return fmt.Errorf("Unknown error")
	}
	return fmt.Errorf("%s: %s", err_matches[0][1], err_matches[0][2])
}

func assertAndParse[T any](resp *http.Response) (T, error) {
	var ret T
	if resp.StatusCode == 503 {
		return ret, errors.New("API server is down")
	}
	_body, _ := io.ReadAll(resp.Body)
	body := string(_body)
	if resp.StatusCode != 200 {
		return ret, parsedError(body)
	}
	err := json.Unmarshal([]byte(body), &ret)
	if err != nil {
		return ret, errors.New(fmt.Sprintf("Error parsing API response: %v", err))
	}
	return ret, nil
}

func CreateClient(options *ClientOptions) (*Client, error) {
	secret_res, err := http.Get(fmt.Sprintf("%s/api/new", options.AuthServer))
	if err != nil {
		return nil, errors.New("Failed to contact with API")
	}
	_secret_body, _ := io.ReadAll(secret_res.Body)
	secret_body := string(_secret_body)
	if secret_res.StatusCode == 503 {
		return nil, errors.New("API server is down")
	} else if secret_res.StatusCode != 200 {
		return nil, parsedError(secret_body)
	}
	authclient := &Client{
		client: http.Client{Transport: secretLoadingTransport{
			secret: secret_body,
		}},
		ClientOptions: options,
		ClientInfo: ClientInfo{
			RespondTo: options.RespondUserOverride,
		},
	}
	return authclient, nil
}
