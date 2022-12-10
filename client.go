package balaboba

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/multierr"
)

// Lang represents balaboba language.
type Lang int

// available languages.
const (
	Rus Lang = iota
	Eng
)

const apiurl = "https://yandex.ru/lab/api/yalm/"

type ClientConfig struct {
	// Lang is language used for text generation.
	// If not specified, Russian is used.
	Lang Lang
	// HTTP is http.Client api to make requests.
	// If not specified, http.DefaultClient is used.
	HTTP *http.Client
}

// Client is Yandex Balaboba service client
type Client struct {
	httpClient *http.Client
	lang       Lang
}

// New makes new balaboba api client.
func New(config ClientConfig) *Client {
	if config.HTTP == nil {
		config.HTTP = http.DefaultClient
	}
	if config.Lang != Rus && config.Lang != Eng {
		panic("invalid language")
	}

	return &Client{
		httpClient: config.HTTP,
		lang:       config.Lang,
	}
}

type responseBase struct {
	ErrorCode int `json:"error"`
}

func (r *responseBase) Error() error {
	if r.ErrorCode != 0 {
		return fmt.Errorf("balaboba error, code: %d", r.ErrorCode)
	}
	return nil
}

// do query balaboba API by endpoint with optional request body and unmarshal result json to reqsponse struct
func (c *Client) do(ctx context.Context, endpoint string, request map[string]any, response interface{}) error {
	return c.request(ctx, apiurl+endpoint, request, response)
}

// do query with optional request body and unmarshal result json to reqsponse struct
func (c *Client) request(ctx context.Context, url string, request map[string]any, response interface{}) error {
	if response == nil {
		return errors.New("destination must not be nil")
	}

	method := http.MethodGet
	var body io.Reader

	if request != nil {
		bodyBytes, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("marshaling request failed: %w", err)
		}

		body = bytes.NewBuffer(bodyBytes)
		method = http.MethodPost
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("request create failed: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("balaboba: response status %s", resp.Status)
	}

	responseDecoder := json.NewDecoder(resp.Body)
	if err := responseDecoder.Decode(response); err != nil {
		raw, err2 := io.ReadAll(io.MultiReader(
			responseDecoder.Buffered(),
			resp.Body,
		))
		multierr.AppendInto(&err, err2)

		return fmt.Errorf("could not parse response, raw=%q: %w", raw, err)
	}

	return nil
}

// Response contains generated text.
type Response struct {
	// Text generated
	Text     string
	BadQuery bool
}

// Generate generates text with passed parameters
func (c *Client) Generate(ctx context.Context, query string, style Style) (*Response, error) {
	var resp struct {
		responseBase
		Query string `json:"query"`
		Text  string `json:"text"`
		// BadQuery is really a boolean: 0 or 1
		BadQuery int `json:"bad_query"`
		// IsCached is really a boolean: 0 or 1
		IsCached  int    `json:"is_cached"`
		Intro     int    `json:"intro"`
		Signature string `json:"signature"`
	}

	if err := c.do(ctx, "text3", map[string]any{
		"query":  query,
		"intro":  style.id,
		"filter": 1,
	}, &resp); err != nil {
		return nil, err
	}

	return &Response{
		BadQuery: resp.BadQuery != 0,
		Text:     resp.Text,
	}, nil
}

// Styles gets list of available generating styles
func (c *Client) Styles(ctx context.Context) ([]Style, error) {
	endpoint := "intros"
	if c.lang == Eng {
		endpoint = "intros_eng"
	}

	var resp struct {
		responseBase
		Styles []Style `json:"intros"`
	}
	return resp.Styles, c.do(ctx, endpoint, nil, &resp)
}
