package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/karalef/balaboba"
	"github.com/urfave/cli/v2"
)

var (
	app = cli.App{
		Name: "balaboba",
	}
	style  = flag.Uint("s", 0, "generation style")
	text   = flag.String("t", "", "text to generate")
	styles = flag.Bool("styles", false, "print all available styles")
	help   = flag.Bool("help", false, "print help")
	eng    = flag.Bool("eng", false, "use english client")
)

func init() {
	flag.Parse()
	if *eng {
		client = balaboba.ClientEng
	} else {
		client = balaboba.ClientRus
	}
}

var client *balaboba.Client

func main() {
	if *help {
		fmt.Printf("%s\n\n%s\n\n", client.About(), client.Warn1())
		flag.PrintDefaults()
		return
	}

	if *styles {
		printStyles()
		return
	}

	if *text == "" {
		*text = strings.Join(flag.Args(), " ")
	}
	if *text == "" {
		fmt.Println("write the text to generate")
		return
	}

	fmt.Println("please wait up to 20 seconds")

	r, err := client.Generate(*text, balaboba.Style(*style))
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(r.Text)
}

func printStyles() {
	allStyles := []balaboba.Style{
		balaboba.Standart,
		balaboba.UserManual,
		balaboba.Recipes,
		balaboba.ShortStories,
		balaboba.WikipediaSipmlified,
		balaboba.MovieSynopses,
		balaboba.FolkWisdom,
	}
	fmt.Println("Styles:")
	for _, style := range allStyles {
		str, desc := style.Description(client.Lang())
		fmt.Println(style, "-", str, "-", desc)
	}
}
