package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
)

type keys struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	hashFunc   hash.Hash
}

var keysValue keys

func main() {
	http.HandleFunc("/encrypt", handleEncrypt)
	http.HandleFunc("/decrypt", handleDecrypt)

	GenerateRSAKeyPair()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GenerateRSAKeyPair() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	publicKey := &privateKey.PublicKey

	keysValue = keys{
		publicKey:  publicKey,
		privateKey: privateKey,
		hashFunc:   sha256.New(),
	}
}

func handleEncrypt(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(w, req)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msgEncryptedInBase64, err := Encrypt(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, msgEncryptedInBase64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Encrypted response body: ", msgEncryptedInBase64)
}

func handleDecrypt(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(w, req)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msgDecrypted, err := Decrypt(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, msgDecrypted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Decrypted response body: ", msgDecrypted)
}

func Encrypt(msg string) (string, error) {
	if keysValue.publicKey == nil {
		return "", errors.New("public key is nil")
	}

	msgBytesEncrypted, err := rsa.EncryptOAEP(keysValue.hashFunc, rand.Reader, keysValue.publicKey, []byte(msg), nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(msgBytesEncrypted), nil
}

func Decrypt(msg string) (string, error) {
	if keysValue.privateKey == nil {
		return "", errors.New("private key is nil")
	}

	msgEncryptedBytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	msgDecryptedBytes, err := rsa.DecryptOAEP(keysValue.hashFunc, rand.Reader, keysValue.privateKey, msgEncryptedBytes, nil)
	if err != nil {
		return "", err
	}

	return string(msgDecryptedBytes), nil
}
