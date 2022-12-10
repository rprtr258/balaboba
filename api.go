package balaboba

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Lang represents balaboba language.
type Lang uint8

// available languages.
const (
	Rus Lang = iota
	Eng

	apiurl = "https://yandex.ru/lab/api/yalm/"
)

var (
	// ClientRus is default russian client.
	ClientRus = NewWithTimeout(Rus, 20*time.Second)

	// ClientEng is default english client.
	ClientEng = NewWithTimeout(Eng, 20*time.Second)
)

// New makes new balaboba api client.
func New(lang Lang, client ...*http.Client) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
		lang:       lang,
	}
	if len(client) > 0 && client[0] != nil {
		c.httpClient = client[0]
	}
	return c
}

// NewWithTimeout makes new balaboba api client with resuests timeout.
func NewWithTimeout(lang Lang, timeout time.Duration) *Client {
	return New(lang, &http.Client{Timeout: timeout})
}

// Client is Yandex Balaboba client.
type Client struct {
	httpClient *http.Client
	lang       Lang
}

type responseBase struct {
	Error int `json:"error"`
}

func (r responseBase) err() int { return r.Error }

type errorable interface{ err() int }

func (c *Client) do(endpoint string, data interface{}, dst errorable) error {
	return c.doContext(context.Background(), endpoint, data, dst)
}

func (c *Client) doContext(ctx context.Context, endpoint string, data interface{}, dst errorable) error {
	err := c.request(ctx, apiurl+endpoint, data, dst)
	if err != nil {
		return err
	}
	if c := dst.err(); c != 0 {
		err = fmt.Errorf("balaboba: error code %d", c)
	}
	return err
}

func (c *Client) request(ctx context.Context, url string, data, dst interface{}) error {
	method := http.MethodGet
	var body io.Reader

	if data != nil {
		var w *io.PipeWriter
		body, w = io.Pipe()
		go func() { w.CloseWithError(json.NewEncoder(w).Encode(data)) }()
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

	if dst == nil {
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(dst); err != nil {
		raw, err2 := io.ReadAll(io.MultiReader(dec.Buffered(), resp.Body))
		if err2 != nil {
			return err2
		}

		return fmt.Errorf("response: %s, error: %w", string(raw), err)
	}

	return nil
}

// Response contains generated text.
type Response struct {
	BadQuery bool
	raw      response
}

// Text generated plus query
func (resp *Response) Text() string {
	return resp.raw.Query + resp.raw.Text
}

type response struct {
	responseBase
	Query     string `json:"query"`
	Text      string `json:"text"`
	BadQuery  uint8  `json:"bad_query"`
	IsCached  uint8  `json:"is_cached"`
	Intro     int    `json:"intro"`
	Signature string `json:"signature"`
}

// GenerateContext generates text with passed parameters.
// It uses the context for the request.
func (c *Client) GenerateContext(ctx context.Context, query string, style Style, filter ...bool) (*Response, error) {
	f := 0
	if len(filter) > 0 && filter[0] {
		f = 1
	}

	var resp Response
	err := c.doContext(ctx, "text3", map[string]interface{}{
		"query":  query,
		"intro":  style,
		"filter": f,
	}, &resp.raw)
	if err != nil {
		return nil, err
	}

	return &Response{
		raw:      resp.raw,
		BadQuery: resp.raw.BadQuery != 0,
	}, nil
}

// Generate generates text with passed parameters.
func (c *Client) Generate(query string, style Style, filter ...bool) (*Response, error) {
	return c.GenerateContext(context.Background(), query, style, filter...)
}
