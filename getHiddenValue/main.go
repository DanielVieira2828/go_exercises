package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {

	url := "https://distopia-a1e2.savi2w.workers.dev/"
	Agent := "DanielVieirass"

	tlsCustomConfig := tls.Config{
		MaxVersion: tls.VersionTLS12,
		MinVersion: tls.VersionTLS12,
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsCustomConfig,
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Erro ao criar a solicitação:", err)
		return
	}

	req.Header.Set("User-Agent", Agent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status da resposta:", resp.Status)
}
