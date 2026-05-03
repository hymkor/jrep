package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hymkor/jegan/types"
	"github.com/hymkor/jegan/unjson"
)

func main1(target *types.JsonPath, name string, r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		err := unjson.Unmarshal(br, func(L types.Line) error {
			if target.Equals(L.Path()) {
				fmt.Printf("%s=%#v\n",
					L.Path().String(),
					types.Unwrap(L.Data()))
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

func mains(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%s JSONPATH Files...", os.Args[0])
	}
	target, err := types.ParseJson(args[0])
	if err != nil {
		return err
	}
	if len(args) < 2 {
		return main1(target, "<STDIN>", os.Stdin)
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
			err1 := main1(target, fn, fd)
			err2 := fd.Close()
			if err1 != nil || err2 != nil {
				return fmt.Errorf("%s: %w", fn, errors.Join(err1, err2))
			}
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
