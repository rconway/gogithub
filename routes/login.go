package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Login zzz
func Login(r *gin.RouterGroup) (login *gin.RouterGroup) {
	login = r.Group("login")

	login.GET("", initiateLogin)
	login.GET("code", receiveCodeGrant)

	// Experiments
	login.GET("rac", func(c *gin.Context) {
		c.String(200, "Request=%v\n\n", c.Request)
		c.String(200, "IsAbs=%v\n", c.Request.URL.IsAbs())
		c.String(200, "URL=%v\n", c.Request.URL)
	})

	return
}

// Initiate authentication with github
func initiateLogin(c *gin.Context) {
	// Construct the URL that we want GitHub to call us back on
	suppliedRedirectURI := c.Query("redirect_uri")
	codeRedirectURI, _ := url.Parse("/login/code")
	codeRedirectURI.Scheme = "http"
	codeRedirectURI.Host = c.Request.Host
	{
		params := codeRedirectURI.Query()
		params.Set("redirect_uri", suppliedRedirectURI)
		codeRedirectURI.RawQuery = params.Encode()
	}

	// Construct URL for request to GitHub's auth endpoint
	githubAuthURL, _ := url.Parse("https://github.com/login/oauth/authorize")
	{
		params := githubAuthURL.Query()
		params.Set("scope", "user:email")
		params.Set("client_id", "6e33c9212621230df631")
		params.Set("redirect_uri", codeRedirectURI.String())
		githubAuthURL.RawQuery = params.Encode()
	}

	// Redirect caller to the GitGub auth endpoint
	c.Redirect(303, githubAuthURL.String())
	c.Abort()
}

// Callback endpoint for code grant provided by GitHub
func receiveCodeGrant(c *gin.Context) {
	code := c.Query("code")

	client := &http.Client{}

	// Convert code grant to access token
	accessToken := getAuthTokenFromCodeGrant(client, code)

	// Get user info from GitHub API using access token
	userInfo := getGitHubUserEmails(client, accessToken)

	// Set user info in Cookies and redirect back
	redirectURI := c.Query("redirect_uri")
	c.SetCookie("user-email", (*userInfo)[0].Email, 60, "/", c.Request.Host, false, true)
	c.Redirect(303, redirectURI)
}

func getAuthTokenFromCodeGrant(client *http.Client, code string) (token string) {
	reqURL, _ := url.Parse("https://github.com/login/oauth/access_token")
	params := reqURL.Query()
	params.Set("client_id", "6e33c9212621230df631")
	params.Set("client_secret", "e280813c8d831c9b9e0218f9bd750558a25a4bed")
	params.Set("code", code)
	reqURL.RawQuery = params.Encode()

	req, _ := http.NewRequest("POST", reqURL.String(), nil)
	res, _ := client.Do(req)

	// Get access token from response
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	v, _ := url.ParseQuery(string(data))
	token = v.Get("access_token")
	return
}

// UserEmails JSON data returned by GitHub user/emails API endpoint
type UserEmails []struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func getGitHubUserEmails(client *http.Client, accessToken string) (userEmails *UserEmails) {
	// Request info from API
	reqURL, _ := url.Parse("https://api.github.com/user/emails")
	req, _ := http.NewRequest("GET", reqURL.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// Interpret JSON response
	userEmails = &UserEmails{}
	err := json.Unmarshal(data, userEmails)
	if err != nil {
		userEmails = nil
	}

	return
}
