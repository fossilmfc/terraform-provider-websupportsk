package websupportsk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	mediaType = "application/json"
)

type Client struct {
	baseURL     string
	userAgent   string
	apiKey      string
	secret      string
	httpClient  *http.Client
	headers     map[string]string
	requestLock sync.Mutex

	Dns DnsClient
}

type ErrorResponse struct {
	Status   string    `json:"status"`
	Messages []string  `json:"messages"`
	Errors   ErrorData `json:"errors"`
	Code     int       `json:"code"`
	Message  string    `json:"message"`
}

type ErrorData struct {
	Content []string `json:"content"`
	Name    []string `json:"name"`
}

func NewClient(apiKey string, secret string, baseUrl string) *Client {
	c := &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		secret:     secret,
		baseURL:    baseUrl,
	}
	c.Dns = DnsClient{client: c}
	return c
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now().UTC()
	formattedTime := timeNow.Format(time.RFC3339)

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Date", formattedTime)
	if body != nil {
		req.Header.Set("Content-Type", mediaType)
	}

	// Authorisation section
	hmacSecret := hmac.New(sha1.New, []byte(c.secret))
	hmacSecret.Write([]byte(fmt.Sprintf("%s %s %d", method, path, timeNow.UnixNano()/int64(time.Second))))
	req.SetBasicAuth(c.apiKey, hex.EncodeToString(hmacSecret.Sum(nil)))
	req.WithContext(ctx)
	return req, nil
}

func (c *Client) Do(req *http.Request, output interface{}) (*http.Response, error) {
	var body []byte
	var err error
	if req.ContentLength > 0 {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			req.Body.Close()
			return nil, err
		}
		req.Body.Close()
	}

	for {
		if req.ContentLength > 0 {
			req.Body = ioutil.NopCloser(bytes.NewReader(body))
		}

		c.requestLock.Lock()
		resp, err := c.httpClient.Do(req)
		c.requestLock.Unlock()
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return resp, err
		}
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))

		if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
			err = errorFromResponse(resp, body)
			if err == nil {
				err = fmt.Errorf("websupport: server responded with status code %d", resp.StatusCode)
			}
			return resp, err
		}
		if output != nil {
			if w, ok := output.(io.Writer); ok {
				_, err = io.Copy(w, bytes.NewReader(body))
			} else {
				err = json.Unmarshal(body, &output)
			}
		}

		return resp, err
	}
}

func errorFromResponse(resp *http.Response, body []byte) error {
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return errors.New("FATAL ERROR: Unable to parse returned error message - incorrect Content-Type")
	}

	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return errors.New("FATAL ERROR: Unable to parse returned error message - Failed to unmarshal response JSON")
	}

	if errorResponse.Code != 0 && errorResponse.Message != "" {
		return errors.New(fmt.Sprintf("ERROR DATA (Status code = %d): %s \n", errorResponse.Code, errorResponse.Message))
	}

	return errors.New(fmt.Sprintf("ERROR DATA (Status code = %d): %s \n", resp.StatusCode, strings.Join(append(errorResponse.Errors.Name, errorResponse.Errors.Content...), " | ")))
}
