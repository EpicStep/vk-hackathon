package main

import (
	"bytes"
	"encoding/json"
	"flag"
	v1 "github.com/EpicStep/vk-hackathon/pkg/api/v1"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	addr := flag.String("addr", "http://localhost:8181", "set server url")
	flag.Parse()

	serverURL, err := url.Parse(*addr)
	if err != nil {
		log.Fatalf("url.Parse failed with error: %s\n", err.Error())
	}

	if strings.HasSuffix(serverURL.Path, "/") {
		log.Fatalf("URL shouldn't have a trailing slash, but %q does", serverURL)
	}

	hash, err := TestUpload(serverURL)
	if err != nil {
		log.Fatalf("TestUpload failed with error: %s\n", err.Error())
	}

	err = TestGetImage(serverURL, hash)
	if err != nil {
		log.Fatalf("TestGetImage failed with error: %s\n", err.Error())
	}

	log.Println("All tests passed 👺.")
}

func TestUpload(serverURL *url.URL) (string, error) {
	var response v1.UploadResponse

	fileDir, _ := os.Getwd()
	fileName := "test.jpg"
	filePath := path.Join(fileDir, "/_examples/", fileName)

	file, _ := os.Open(filePath)
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	part.Write(content)
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, serverURL.String()+"/image", &body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

func TestGetImage(serverURL *url.URL, hash string) error {
	req, err := http.NewRequest(http.MethodGet, serverURL.String()+"/image", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("id", hash)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
