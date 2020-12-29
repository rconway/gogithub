package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main1() {

	client := &http.Client{}

	authURL, err := url.Parse("https://github.com/login/oauth/authorize")
	panicOnError(err)
	params := authURL.Query()
	params.Set("client_id", "6e33c9212621230df631")
	params.Set("fred", "burt &zlarry")
	authURL.RawQuery = params.Encode()
	fmt.Printf("URL = %v\n", authURL.String())
	req, err := http.NewRequest("GET", authURL.String(), nil)
	panicOnError(err)

	resp, err := client.Do(req)
	panicOnError(err)
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}

	err = scanner.Err()
	panicOnError(err)
}

func main2() {
	Server()
}

// MiddlewareConstructor zzz
type MiddlewareConstructor func(http.Handler) http.Handler

// Richard zzz
type Richard struct {
	name string
}

// ServeHTTP zzz
func (rich Richard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("hello %v\n", rich.name)))
}

// RichardMiddleware zzz
func (rich Richard) RichardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rich.ServeHTTP(w, r)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// RichardHandler zzz
func RichardHandler(name string) MiddlewareConstructor {
	rich := Richard{name}
	return func(next http.Handler) http.Handler {
		return rich.RichardMiddleware(next)
	}
}

// LoginHandler zzz
type LoginHandler struct {
	prefix string
	http.ServeMux
}

// NewLoginHandler zzz
func NewLoginHandler(prefix string) *LoginHandler {
	login := &LoginHandler{prefix, *http.NewServeMux()}
	login.Handle(prefix+"/", RichardHandler("burt")(nil))
	login.HandleFunc(prefix+"/check/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			io.WriteString(w, fmt.Sprintf("r.RequestURI = %v\n", r.RequestURI))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%v %v\n", http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		}
	})
	return login
}

func main3() {
	root := http.NewServeMux()
	h1 := RichardHandler("fred")
	h2 := RichardHandler("bob")
	root.Handle("/", h1(h2(nil)))

	login := NewLoginHandler("/login")
	root.Handle("/login/", login)

	// http.ListenAndServe(":8080", h1(h2(nil)))
	http.ListenAndServe(":8080", root)
}

func main() {
	main3()
}
