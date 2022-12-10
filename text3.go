package balaboba

import (
	"context"
)

// Response contains generated text.
type Response struct {
	// Contains the query + generated continuation.
	//
	// If BadQuery is true it contains the
	// bad query text in the Client language.
	Text     string
	BadQuery bool

	raw response
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

// Generate generates text with passed parameters.
func (c *Client) Generate(query string, style Style, filter ...bool) (*Response, error) {
	return c.GenerateContext(context.Background(), query, style, filter...)
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

	if resp.raw.BadQuery != 0 {
		return &Response{
			BadQuery: true,
		}, nil
	}

	return &Response{
		raw:  resp.raw,
		Text: resp.raw.Query + resp.raw.Text,
	}, nil
}
