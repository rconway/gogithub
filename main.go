package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main2() {

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

func main() {
	Server()
}
