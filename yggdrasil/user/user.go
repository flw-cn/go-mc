package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var ServicesURL = "https://api.minecraftservices.com"

var client = http.DefaultClient

type KeyPairResp struct {
	KeyPair struct {
		PrivateKey string `json:"privateKey"`
		PublicKey  string `json:"publicKey"`
	} `json:"keyPair"`
	PublicKeySignature string    `json:"publicKeySignature"`
	ExpiresAt          time.Time `json:"expiresAt"`
	RefreshedAfter     time.Time `json:"refreshedAfter"`
}

func GetOrFetchKeyPair(accessToken string) (KeyPairResp, error) {
	return fetchKeyPair(accessToken) // TODO: cache
}

func fetchKeyPair(accessToken string) (KeyPairResp, error) {
	var keyPairResp KeyPairResp
	err := post("/player/certificates", accessToken, &keyPairResp)
	return keyPairResp, err
}

func post(endpoint string, accessToken string, resp interface{}) error {
	rowResp, err := rawPost(endpoint, accessToken)
	if err != nil {
		return fmt.Errorf("request fail: %v", err)
	}
	defer rowResp.Body.Close()
	err = json.NewDecoder(rowResp.Body).Decode(resp)
	if err != nil {
		return fmt.Errorf("parse resp fail: %v", err)
	}

	return nil
}

func rawPost(endpoint string, accessToken string) (*http.Response, error) {
	PostRequest, err := http.NewRequest(
		http.MethodPost,
		ServicesURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("make request error: %v", err)
	}

	PostRequest.Header.Set("Authorization", "Bearer "+accessToken)
	PostRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Do
	return client.Do(PostRequest)
}
