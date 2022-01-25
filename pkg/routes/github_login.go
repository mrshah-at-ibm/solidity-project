package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mrshah-at-ibm/kaleido-project/pkg/config"
)

func (r *Routes) githubLoginRoute(w http.ResponseWriter, req *http.Request) {

	githubClientID, err := config.GetGithubClientID()
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("Github login not supported. Ref: clientid"))
		return
	}
	githubRedirectURL, err := config.GetGithubRedirectURL()
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("Github login not supported. Ref: redirecturl"))
		return
	}

	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		githubRedirectURL,
	)

	http.Redirect(w, req, redirectURL, 301)
}

func (r *Routes) githubLoginCallbackRoute(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	githubAccessToken, _ := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	loggedinHandler(w, req, githubData)
}

func loggedinHandler(w http.ResponseWriter, req *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	token := config.GenerateToken()
	err := config.SaveToken(token)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("Unable to generate token"))
		return
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, "Token generated: %s", token)
}

func getGithubAccessToken(code string) (string, error) {
	githubClientID, err := config.GetGithubClientID()
	if err != nil {
		return "", err
	}

	githubSecret, err := config.GetGithubClientSecret()
	if err != nil {
		return "", err
	}

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     githubClientID,
		"client_secret": githubSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return "", reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	// Response body converted to stringified JSON
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	err = json.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}
