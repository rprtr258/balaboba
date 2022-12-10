package balaboba

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/multierr"
)

// Lang represents balaboba language.
type Lang uint8

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

func (c *Client) do(ctx context.Context, endpoint string, request map[string]any, response interface{}) error {
	return c.request(ctx, apiurl+endpoint, request, response)
}

func (c *Client) request(ctx context.Context, url string, request map[string]any, response interface{}) error {
	if response == nil {
		panic("destination must not be nil")
	}

	method := http.MethodGet
	var body io.Reader

	if request != nil {
		var w *io.PipeWriter
		body, w = io.Pipe()
		go func() { w.CloseWithError(json.NewEncoder(w).Encode(request)) }()
		method = http.MethodPost
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("balaboba: response status %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(response); err != nil {
		raw, err2 := io.ReadAll(io.MultiReader(dec.Buffered(), resp.Body))
		multierr.AppendInto(&err, err2)

		return fmt.Errorf("could not parse response, raw=%q: %w", raw, err)
	}

	return nil
}

// Response contains generated text.
type Response struct {
	BadQuery bool
	raw      generateResponse
}

// Text generated plus query
func (resp *Response) Text() string {
	return resp.raw.Query + resp.raw.Text
}

type generateResponse struct {
	responseBase
	Query     string `json:"query"`
	Text      string `json:"text"`
	BadQuery  uint8  `json:"bad_query"`
	IsCached  uint8  `json:"is_cached"`
	Intro     int    `json:"intro"`
	Signature string `json:"signature"`
}

// Generate generates text with passed parameters.
// It uses the context for the request.
func (c *Client) Generate(ctx context.Context, query string, style Style) (*Response, error) {
	var resp Response
	err := c.do(ctx, "text3", map[string]any{
		"query":  query,
		"intro":  style.id,
		"filter": 1,
	}, &resp.raw)
	if err != nil {
		return nil, err
	}

	return &Response{
		raw:      resp.raw,
		BadQuery: resp.raw.BadQuery != 0,
	}, nil
}
