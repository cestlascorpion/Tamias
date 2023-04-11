package core

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

func Parse(file []byte) (int, int, string, error) {
	img, name, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return 0, 0, "", err
	}
	return img.Bounds().Dx(), img.Bounds().Dy(), name, nil
}

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %s", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func Upload(url string, headers map[string]string, fields map[string]string, name string, content []byte, timeout time.Duration) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range fields {
		_ = writer.WriteField(key, value)
	}

	part, err := writer.CreateFormFile("file", name)
	_, err = io.Copy(part, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Charset", "UTF-8")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
