package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hymkor/struct2flag"

	"github.com/hymkor/jegan/types"
	"github.com/hymkor/jegan/unjson"
)

type Application struct {
	ValueOnly bool `flag:"value-only,Output only the value, omitting the key and surrounding syntax."`
	NoComma   bool `flag:"no-comma,Suppress the trailing comma, if present."`
	NoNewLine bool `flag:"no-newline,Do not append a newline after each input."`

	target *types.JsonPath
}

func (app *Application) Process(name string, r io.Reader) error {
	const OFF = 99999
	found := false
	if !app.NoNewLine {
		defer func() {
			if found {
				fmt.Println()
			}
		}()
	}

	br := bufio.NewReader(r)
	nest := OFF
	for {
		err := unjson.Unmarshal(br, func(L types.Line) error {
			n := L.Nest()
			if n == nest {
				if app.NoComma {
					L.DumpWithoutComma(os.Stdout)
				} else {
					L.Dump(os.Stdout)
				}
				nest = OFF
			} else if n > nest {
				L.Dump(os.Stdout)
			}
			if app.target.Equals(L.Path()) {
				found = true
				var v interface {
					Dump(io.Writer)
					DumpWithoutComma(io.Writer)
				} = L
				if app.ValueOnly {
					if p, ok := L.(*types.Pair); ok {
						v = &p.Item
					}
				}
				if app.NoComma {
					v.DumpWithoutComma(os.Stdout)
				} else {
					v.Dump(os.Stdout)
				}
				val := types.Unwrap(L.Data())
				if _, ok := val.(types.Mark); ok {
					nest = n
				}
			}
			return nil
		})
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

func (app *Application) Run(args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("%s JSONPATH Files...", os.Args[0])
	}
	app.target, err = types.ParseJson(args[0])
	if err != nil {
		return
	}
	if len(args) < 2 {
		return app.Process("<STDIN>", os.Stdin)
	}
	for _, arg1 := range args[1:] {
		filenames, err := filepath.Glob(arg1)
		if err != nil || len(filenames) <= 0 {
			filenames = []string{arg1}
		}
		for _, fn := range filenames {
			fd, err := os.Open(fn)
			if err != nil {
				return err
			}
			err1 := app.Process(fn, fd)
			err2 := fd.Close()
			if err1 != nil || err2 != nil {
				return fmt.Errorf("%s: %w", fn, errors.Join(err1, err2))
			}
		}
	}
	return nil
}

func main() {
	app := new(Application)
	struct2flag.BindDefault(app)
	flag.Parse()

	if err := app.Run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
