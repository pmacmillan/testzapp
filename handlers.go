package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func zoomAuthRedirectHandler(w http.ResponseWriter, req *http.Request) {
	var err error

	vars := mux.Vars(req)

	token, err := oauthConfig.Exchange(oauth2.NoContext, vars["code"])

	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	fmt.Fprintf(w, "%s\n", token.AccessToken)

	// request a deep link
	deepLinkReqURL := "https://devep.zoomdev.us/v2/zoomapp/deeplink/"
	deepLinkReqBody := []byte(`{"target": "panel", "action": "go"}`)

	dlr, err := http.NewRequest("POST", deepLinkReqURL, bytes.NewBuffer(deepLinkReqBody))
	dlr.Header.Set("Content-Type", "application/json")
	dlr.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(dlr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "response Status: %s\n", resp.Status)
	fmt.Fprintf(w, "response Headers: %s\n", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "response Body: %s\n", string(body))
}
