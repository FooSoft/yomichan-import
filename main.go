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
	flagPrettyJson = 1 << iota
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s file_format input_file output_file\n\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func outputJson(fileFormat, inputPath, outputDir string, flags int) error {
	handlers := map[string]func(string, io.Reader, int) error{
		"meta":     outputMetaJson,
		"edict":    outputEdictJson,
		"enamdict": outputJmnedictJson,
		"kanjidic": outputKanjidicJson,
	}

	handler, ok := handlers[fileFormat]
	if !ok {
		return errors.New("unrecognized file format")
	}

	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	return handler(outputDir, input, flags)
}

func main() {
	prettyJson := flag.Bool("prettyJson", false, "output prettified json")

	flag.Usage = usage
	flag.Parse()

	var flags int
	if *prettyJson {
		flags |= flagPrettyJson
	}

	if flag.NArg() == 3 {
		if err := outputJson(flag.Arg(0), flag.Arg(1), flag.Arg(2), flags); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
		os.Exit(2)
	}
}
