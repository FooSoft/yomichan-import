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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	flagPretty = 1 << iota
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [edict|enamdict|kanjidic|epwing] input-path [output-dir]\n\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func exportDb(inputPath, outputDir, format, title string, pretty bool) error {
	handlers := map[string]func(string, string, io.Reader, bool) error{
		"edict":    jmdictExportDb,
		"enamdict": jmnedictExportDb,
		"kanjidic": kanjidicExportDb,
		"epwing":   epwingExportDb,
	}

	handler, ok := handlers[format]
	if !ok {
		return errors.New("unrecognized dictionray format")
	}

	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	return handler(outputDir, title, input, pretty)
}

func serveDb(serveDir string, port int) error {
	log.Printf("starting HTTP server on port %d...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(serveDir)))
}

func main() {
	var (
		serve  = flag.Bool("serve", false, "serve JSON over HTTP")
		port   = flag.Int("port", 9876, "port to serve JSON on")
		pretty = flag.Bool("pretty", false, "output prettified JSON")
		title  = flag.String("title", "", "dictionary title")
	)

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 && flag.NArg() != 3 {
		usage()
		os.Exit(2)
	}

	var (
		format    = flag.Arg(0)
		inputPath = flag.Arg(1)
		outputDir string
	)

	if flag.NArg() == 3 {
		outputDir = flag.Arg(3)
	} else {
		var err error
		outputDir, err = ioutil.TempDir("", "yomichan_tmp_")
		if err != nil {
			log.Fatal(err)
		}
	}

	if *title == "" {
		t := time.Now()
		*title = fmt.Sprintf("%s-%s", format, t.Format("20060102150405"))
	}

	if err := exportDb(inputPath, outputDir, format, *title, *pretty); err != nil {
		log.Fatal(err)
	}

	if *serve {
		if err := serveDb(outputDir, *port); err != nil {
			log.Fatal(err)
		}
	}
}
