package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/karalef/balaboba"
	"github.com/urfave/cli/v2"
)

var (
	app = cli.App{
		Name:  "balaboba",
		Usage: "generate text using yandex's balaboba neural network",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "eng",
				DefaultText: "use english language",
				Value:       false,
			},
			&cli.GenericFlag{
				Name:    "style",
				Aliases: []string{"s"},
				Usage:   "generation style",
				Value:   &balaboba.Standart,
			},
		},
		UsageText: `Нейросеть не знает, что говорит, и может сказать всякое — если что, не обижайтесь.
Распространяя получившиеся тексты, помните об ответственности. (18+)

Генератор может выдавать очень странные тексты.
Пожалуйста, будьте разумны, распространяя их.
Подумайте, не будет ли текст обидным для кого-то и не станет ли его публикация нарушением закона.

Балабоба не принимает запросы на острые темы, например, про политику или религию.
Люди могут слишком серьёзно отнестись к сгенерированным текстам.

Вероятность того, что запрос задаёт одну из острых тем, определяет нейросеть, обученная на оценках случайных людей.
Но она может перестараться или, наоборот, что-то пропустить.`,
		Commands: []*cli.Command{{
			Name:  "styles",
			Usage: "print all available styles",
			Action: func(ctx *cli.Context) error {
				// TODO: change to stylesByID
				allStyles := []balaboba.Style{
					balaboba.Standart,
					balaboba.ShortStories,
					balaboba.WikipediaSipmlified,
					balaboba.MovieSynopses,
					balaboba.FolkWisdom,
					balaboba.UserManual,
					balaboba.Recipes,
				}

				fmt.Println("Styles:")
				for _, style := range allStyles {
					fmt.Println(style.String())
				}

				return nil
			},
		}, {
			Name:  "generate",
			Usage: "generate text",
			Action: func(ctx *cli.Context) error {
				args := ctx.Args().Slice()
				eng := ctx.Bool("eng")
				style := *ctx.Generic("style").(*balaboba.Style)

				text := strings.Join(args, " ")
				if text == "" {
					return errors.New("write the text to generate")
				}

				lang := balaboba.Rus
				if eng {
					lang = balaboba.Eng
				}

				client := balaboba.New(balaboba.ClientConfig{
					Lang: lang,
				})

				r, err := client.Generate(ctx.Context, text, style)
				if err != nil {
					return err
				}

				fmt.Println(r.Text())

				return nil
			},
		}},
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
