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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	DEFAULT_STRIDE = 10000
	DEFAULT_PORT   = 9876
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] input-path [output-dir]\n", path.Base(os.Args[0]))
	fmt.Fprint(os.Stderr, "https://foosoft.net/projects/yomichan-import/\n\n")
	fmt.Fprint(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func exportDb(inputPath, outputDir, format, language, title string, stride int, pretty bool) error {
	handlers := map[string]func(string, string, string, string, int, bool) error{
		"edict":    jmdictExportDb,
		"enamdict": jmnedictExportDb,
		"kanjidic": kanjidicExportDb,
		"epwing":   epwingExportDb,
	}

	handler, ok := handlers[format]
	if !ok {
		return errors.New("unrecognized dictionray format")
	}

	log.Printf("converting '%s' to '%s' in '%s' format...", inputPath, outputDir, format)
	return handler(inputPath, outputDir, language, title, stride, pretty)
}

func serveDb(serveDir string, port int) error {
	log.Printf("starting dictionary server on port %d...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(serveDir)))
}

func makeTmpDir() (string, error) {
	return ioutil.TempDir("", "yomichan_tmp_")
}

func main() {
	var (
		format   = flag.String("format", "", "dictionary format [edict|enamdict|kanjidic|epwing]")
		language = flag.String("language", "english", "dictionary language (if applicable)")
		title    = flag.String("title", "", "dictionary title")
		port     = flag.Int("port", DEFAULT_PORT, "port to serve dictionary JSON on")
		stride   = flag.Int("stride", DEFAULT_STRIDE, "dictionary bank stride")
		pretty   = flag.Bool("pretty", false, "output prettified dictionary JSON")
		serve    = flag.Bool("serve", false, "serve dictionary JSON for extension")
	)

	flag.Usage = usage
	flag.Parse()

	var (
		inputPath string
		outputDir string
	)

	if flag.NArg() == 0 {
		if err := gui(); err == nil {
			return
		} else {
			usage()
			os.Exit(2)
		}
	} else {
		inputPath = flag.Arg(0)
		if flag.NArg() > 1 {
			outputDir = flag.Arg(1)
		}
	}

	if _, err := os.Stat(inputPath); err != nil {
		log.Fatalf("dictionary path '%s' does not exist", inputPath)
	}

	if *format == "" {
		if *format = detectFormat(inputPath); *format == "" {
			log.Fatal("failed to detect dictionary format")
		}
	}

	if outputDir == "" {
		var err error
		if outputDir, err = makeTmpDir(); err != nil {
			log.Fatal(err)
		}

		*serve = true
	}

	if err := exportDb(inputPath, outputDir, *format, *language, *title, *stride, *pretty); err != nil {
		log.Fatal(err)
	}

	if *serve {
		if err := serveDb(outputDir, *port); err != nil {
			log.Fatal(err)
		}
	}
}
