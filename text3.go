package balaboba

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

// Style of generating text.
type Style struct {
	id          int
	title       string
	description string
}

var (
	Standart = Style{
		id:          0,
		title:       "Без стиля",
		description: "Напишите что-нибудь и получите продолжение от Балабобы",
	}
	ShortStories = Style{
		id:          6,
		title:       "Короткие истории",
		description: "Начните писать историю, а Балабобы продолжит — иногда страшно, но чаще смешно",
	}
	WikipediaSipmlified = Style{
		id:          8,
		title:       "Короче, Википедия",
		description: "Напишите какое-нибудь слово, а Балабоба даст этому определение",
	}
	MovieSynopses = Style{
		id:          9,
		title:       "Синопсисы фильмов",
		description: "Напишите название фильма (существующего или нет), а Балабоба расскажет вам, о чем он",
	}
	FolkWisdom = Style{
		id:          11,
		title:       "Народные мудрости",
		description: "Напишите что-нибудь и получите народную мудрость",
	}
	UserManual = Style{
		id:          24,
		title:       "Инструкции по применению",
		description: "Перечислите несколько предметов, а Балабоба придумает, как их использовать",
	}
	Recipes = Style{
		id:          25,
		title:       "Рецепты",
		description: "Перечислите съедобные ингредиенты, а Балабоба придумает рецепт с ними",
	}

	stylesByID = map[int]Style{
		Standart.id:            Standart,
		ShortStories.id:        ShortStories,
		WikipediaSipmlified.id: WikipediaSipmlified,
		MovieSynopses.id:       MovieSynopses,
		FolkWisdom.id:          FolkWisdom,
		UserManual.id:          UserManual,
		Recipes.id:             Recipes,
	}
)

func (s *Style) Set(value string) error {
	id, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("invalid style")
	}

	style, ok := stylesByID[id]
	if !ok {
		return errors.New("invalid style")
	}

	*s = style
	return nil
}

func (style *Style) String() string {
	return fmt.Sprintf("%2d %12s: %s", style.id, style.title, style.description)
}
