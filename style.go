package balaboba

import (
	"encoding/json"
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
	Football = Style{
		id:          34,
		title:       "ЧМ 2022",
		description: "Введите заголовок новости про футбол и получите готовый комментарий",
	}

	StylesByID = map[int]Style{
		Standart.id:            Standart,
		ShortStories.id:        ShortStories,
		WikipediaSipmlified.id: WikipediaSipmlified,
		MovieSynopses.id:       MovieSynopses,
		FolkWisdom.id:          FolkWisdom,
		UserManual.id:          UserManual,
		Recipes.id:             Recipes,
		Football.id:            Football,
	}
)

// TODO: remove from public interface
func (s *Style) Set(value string) error {
	id, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("invalid style")
	}

	if style, ok := StylesByID[id]; ok {
		*s = style
	} else {
		*s = Style{id: id}
	}

	return nil
}

func (style *Style) String() string {
	return fmt.Sprintf("%d %s: %s", style.id, style.title, style.description)
}

// UnmarshalJSON is json.Unmarshaler interface implementation.
func (style *Style) UnmarshalJSON(b []byte) error {
	var rep [3]interface{}
	err := json.Unmarshal(b, &rep)
	if err != nil {
		return err
	}

	*style = Style{
		id:          int(rep[0].(float64)),
		title:       rep[1].(string),
		description: rep[2].(string),
	}

	return nil
}
