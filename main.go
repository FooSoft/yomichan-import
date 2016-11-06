/*
 * Copyright (c) 2016 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const (
	flagPretty = 1 << iota
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [edict|enamdict|kanjidic] input output\n\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func exportDb(inputPath, outputDir, format, title string, flags int) error {
	handlers := map[string]func(string, string, io.Reader, int) error{
		"edict":    exportJmdictDb,
		"enamdict": exportJmnedictDb,
		"kanjidic": exportKanjidicDb,
	}

	handler, ok := handlers[format]
	if !ok {
		return errors.New("unrecognized file format")
	}

	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	return handler(outputDir, title, input, flags)
}

func main() {
	pretty := flag.Bool("pretty", false, "output prettified json")
	format := flag.String("format", "", "dictionary format")
	title := flag.String("title", "", "dictionary title")

	flag.Usage = usage
	flag.Parse()

	var flags int
	if *pretty {
		flags |= flagPretty
	}

	if flag.NArg() == 2 && len(*format) > 0 && len(*title) > 0 {
		if err := exportDb(flag.Arg(0), flag.Arg(1), *format, *title, flags); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
		os.Exit(2)
	}
}
