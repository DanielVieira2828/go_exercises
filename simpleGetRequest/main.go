package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("O status da requisição foi: ", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp.Body.Close()

	fmt.Println("O corpo da resposta foi: ", string(body))
}
