package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type uploadResponse struct {
	Response Response `json:"response"`
}
type Story struct {
	ID      int `json:"id"`
	OwnerID int `json:"owner_id"`
}
type Response struct {
	Story Story `json:"story"`
}

func postFile(filename string, targetUrl string) uploadResponse {
	unResp := uploadResponse{}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("photo", filename)
	if err != nil {
		log.Fatal(err)
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Fatal(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&unResp)
	if err != nil {
		log.Fatal(err)
	}
	return unResp
}

func sendStoryToVK(filename, uploadUrl string) uploadResponse {
	storyFile := fmt.Sprintf("tmp/%s.png", filename)

	uploadResp := postFile(storyFile, uploadUrl)

	return uploadResp
}
