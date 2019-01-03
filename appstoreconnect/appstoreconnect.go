package appstoreconnect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/godrei/go-appstoreconnect/debug"
	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.appstoreconnect.apple.com/"
)

type service struct {
	client *Client
}

// Client communicate with the Apple API
type Client struct {
	keyID             string
	issuerID          string
	privateKeyContent []byte

	token       *jwt.Token
	signedToken string

	client  *http.Client
	BaseURL *url.URL

	common     service // Reuse a single struct instead of allocating one for each service on the heap.
	TestFlight *TestFlightService
}

// NewClient creates a new client
func NewClient(keyID, issuerID, privateKeyPath string) (*Client, error) {
	httpClient := http.DefaultClient
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	privateKeyContent, err := fileutil.ReadBytesFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	c := &Client{
		keyID:             keyID,
		issuerID:          issuerID,
		privateKeyContent: privateKeyContent,

		client:  httpClient,
		BaseURL: baseURL,
	}
	c.common.client = c
	c.TestFlight = (*TestFlightService)(&c.common)

	return c, nil
}

// ensureSignedToken makes sure that the JWT auth token is not expired
// and return a signed key
func (c *Client) ensureSignedToken() (string, error) {
	if c.token != nil {
		claim := c.token.Claims.(claims)
		expiration := time.Unix(int64(claim.Expiration), 0)
		// token is marked valid if it will not expire in the upcoming 10 sec
		if expiration.After(time.Now().Add(10 * time.Second)) {
			return c.signedToken, nil
		}
	}

	c.token = createToken(c.keyID, c.issuerID)
	var err error
	if c.signedToken, err = signToken(c.token, c.privateKeyContent); err != nil {
		return "", err
	}
	return c.signedToken, nil
}

// NewRequest creates a new http.Request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	urlStr = "v1/" + urlStr
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	signedToken, err := c.ensureSignedToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+signedToken)

	log.Debugf("NewRequest:")
	debug.PrintRequest(req)

	return req, nil
}

// ErrorResponseError ...
type ErrorResponseError struct {
	Code   string      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	ID     string      `json:"id,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Source interface{} `json:"source,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Response *http.Response
	Errors   []ErrorResponseError `json:"errors,omitempty"`
}

// Error ...
func (r ErrorResponse) Error() string {
	m := fmt.Sprintf("%v %v: %d\n", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode)
	var s string

	for _, err := range r.Errors {
		m += s + fmt.Sprintf("- %v %v", err.Title, err.Detail)
		s = "\n"
	}

	return m
}

func checkResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Debugf("do:")
	debug.PrintResponse(resp)

	if err := checkResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}

	return resp, err
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
