package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"
	reader := bufio.NewReader(os.Stdin)

	// Prompt for login credentials
	fmt.Println("Введите логин")
	login, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	login = strings.TrimSpace(login)

	fmt.Println("Введите пароль")
	password, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	password = strings.TrimSpace(password)

	// Obtain token using credentials and persist it
	if login != "" && password != "" {
		if token, err := obtainTokenFromCreds(login, password); err == nil && token != "" {
			_ = os.WriteFile("jwt_token.txt", []byte(token), 0600)
		}
	}

	// Prompt for long URL
	fmt.Println("Введите длинный URL")
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")

	response, err := sendRequest(endpoint, long)
	if err != nil {
		panic(err)
	}

	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

// sendRequest sends a POST request to the given endpoint with the provided URL
func sendRequest(endpoint, url string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(url))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Attach Authorization header if token is available; otherwise try to obtain it via auth service
	token := getAuthToken()
	if token == "" {
		if t, err := obtainToken(); err == nil && t != "" {
			token = t
			// persist for next runs
			_ = os.WriteFile("jwt_token.txt", []byte(token), 0600)
		}
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
		// Also set cookie for middlewares that expect it
		request.AddCookie(&http.Cookie{Name: "token", Value: token, Path: "/"})
	}

	return client.Do(request)
}

// getAuthToken returns token from AUTH_TOKEN env or jwt_token.txt if present
func getAuthToken() string {
	if t, ok := os.LookupEnv("AUTH_TOKEN"); ok && strings.TrimSpace(t) != "" {
		return strings.TrimSpace(t)
	}
	// Fallback to local file
	if data, err := os.ReadFile("jwt_token.txt"); err == nil {
		return strings.TrimSpace(string(data))
	}
	return ""
}

// obtainToken tries to log in using login_data.json against the auth service
func obtainToken() (string, error) {
	// Read credentials
	type creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	data, err := os.ReadFile("login_data.json")
	if err != nil {
		return "", err
	}
	c := creds{}
	if err := json.Unmarshal(data, &c); err != nil {
		return "", err
	}
	if strings.TrimSpace(c.Login) == "" || strings.TrimSpace(c.Password) == "" {
		return "", fmt.Errorf("empty login or password in login_data.json")
	}

	// Determine auth server URL
	authURL := os.Getenv("AUTH_SERVER_URL")
	if strings.TrimSpace(authURL) == "" {
		authURL = "http://localhost:8082"
	}

	// Perform login
	bodyBytes, _ := json.Marshal(c)
	req, err := http.NewRequest(http.MethodPost, authURL+"/api/auth/login", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", strings.TrimSpace(string(b)))
	}
	var respJSON struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return "", err
	}
	if strings.TrimSpace(respJSON.Token) == "" {
		return "", fmt.Errorf("empty token in response")
	}
	return respJSON.Token, nil
}

// obtainTokenFromCreds logs in to auth service with provided login/password
func obtainTokenFromCreds(login, password string) (string, error) {
	type creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	c := creds{Login: strings.TrimSpace(login), Password: strings.TrimSpace(password)}
	if c.Login == "" || c.Password == "" {
		return "", fmt.Errorf("empty credentials")
	}

	authURL := os.Getenv("AUTH_SERVER_URL")
	if strings.TrimSpace(authURL) == "" {
		authURL = "http://localhost:8082"
	}

	bodyBytes, _ := json.Marshal(c)
	req, err := http.NewRequest(http.MethodPost, authURL+"/api/auth/login", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", strings.TrimSpace(string(b)))
	}
	var respJSON struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return "", err
	}
	if strings.TrimSpace(respJSON.Token) == "" {
		return "", fmt.Errorf("empty token in response")
	}
	return respJSON.Token, nil
}
