package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/peteraba/xemmet/counter"
)

type Mode string

const (
	ModeHTML Mode = "html"
	ModeXML  Mode = "xml"
	ModeHTMX Mode = "htmx"
)

const (
	defaultIndentation = "    "
)

func main() {
	app := &cli.App{
		Name:  "xemmet",
		Usage: "An Emmet.HTML rewrite in GO",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "mode",
				Value: string(ModeHTML),
				Usage: "Output mode (html, xml, htmx)",
			},
			&cli.StringFlag{
				Name:  "indentation",
				Value: defaultIndentation,
				Usage: "Indentation to apply (not multiline if empty)",
			},
			&cli.IntFlag{
				Name:  "depth",
				Value: 0,
				Usage: "Initial indentation level to use",
			},
			&cli.BoolFlag{
				Name:  "inline",
				Value: false,
				Usage: "Enable debug mode",
			},
			&cli.StringFlag{
				Name:  "tabStop",
				Value: "",
				Usage: "Unique set of characters to surround variable names used for tabs stops (if empty, then tab stops will not be added)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			var (
				str         = cCtx.Args().First()
				mode        = Mode(cCtx.String("mode"))
				indentation = cCtx.String("indentation")
				depth       = cCtx.Int("depth")
				multiline   = !cCtx.Bool("inline")
				tabStop     = cCtx.String("tabStop")
			)

			got, err := Xemmet(mode, str, indentation, depth, multiline, tabStop)
			if err != nil {
				return err
			}

			fmt.Print(got) // nolint: forbidigo

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const (
	ErrTokenizingMsg = "error tokenizing string"
)

func Xemmet(mode Mode, str string, indentation string, depth int, multiline bool, tabStopWrapper string) (string, error) {
	l := NewLexer(mode)

	tokens, _, err := l.Tokenize([]rune(str), false)
	if err != nil {
		return "", errors.Wrap(err, ErrTokenizingMsg)
	}

	s := NewSnippeter(mode)
	tokens = s.Walk(tokens)

	elemList := Build(tokens, 1, 1)

	counter.ResetGlobalTabStopCounter()

	return strings.Trim(elemList.HTML(mode, indentation, depth, multiline, tabStopWrapper), "\n\t\r "), nil
}
