package acceptance_tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getJWT(ctx context.Context, hostName string) (string, error) {
	url := hostName + "/sessions"

	authToken := createBasicAuthToken(Username, Password)

	client := &http.Client{}

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", "Basic "+authToken)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	res := struct {
		Jwt string `json:"jwt"`
	}{}
	if err = json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if res.Jwt == "" {
		const msg = "no JWT token"
		log.Print(msg, string(body))
		return "", ErrNoJWT
	}

	return res.Jwt, nil
}

func createBasicAuthToken(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
