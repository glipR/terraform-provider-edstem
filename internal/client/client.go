package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	CourseID   string
	Token      string
	HTTPClient *http.Client
}

func NewClient(course_id, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	if course_id != nil {
		c.CourseID = *course_id
	}

	if token != nil {
		c.Token = *token
	}

	return &c, nil
}

func (c *Client) HTTPRequest(path, method string, body bytes.Buffer, boundary *string) (closer io.ReadCloser, err error) {
	fmt.Println("Requesting", c.requestPath((path)), "with method", method)
	fmt.Println(body.String())
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Token", c.Token)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("content-type", "application/json")
	}
	if boundary != nil {
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", *boundary))
		req.Header.Set("content-type", fmt.Sprintf("multipart/form-data; boundary=%s", *boundary))
	}

	resp, err := c.HTTPClient.Do(req)
	fmt.Println(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%s/%s", "https://edstem.org/api", path)
}
