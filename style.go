package balaboba

import (
	"errors"
	"fmt"
	"strconv"
)

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
	return fmt.Sprintf("%2d %-24s: %s", style.id, style.title, style.description)
}
