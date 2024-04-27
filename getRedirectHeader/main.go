package main

import (
	"fmt"
	"net/http"
)

func main() {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get("https://distopia.savi2w.workers.dev/")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Status Code: ", resp.Status)

	resp.Body.Close()

	fmt.Println("Body response: ", resp.Header.Get("Distopia"))
}
