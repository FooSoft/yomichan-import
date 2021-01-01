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
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	yomichan "github.com/FooSoft/yomichan-import"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] input-path output-path\n", path.Base(os.Args[0]))
	fmt.Fprint(os.Stderr, "https://foosoft.net/projects/yomichan-import/\n\n")
	fmt.Fprint(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	var (
		format   = flag.String("format", yomichan.DefaultFormat, "dictionary format [edict|enamdict|epwing|kanjidic|rikai]")
		language = flag.String("language", yomichan.DefaultLanguage, "dictionary language (if supported)")
		title    = flag.String("title", yomichan.DefaultTitle, "dictionary title")
		stride   = flag.Int("stride", yomichan.DefaultStride, "dictionary bank stride")
		pretty   = flag.Bool("pretty", yomichan.DefaultPretty, "output prettified dictionary JSON")
	)

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	if err := yomichan.ExportDb(flag.Arg(0), flag.Arg(1), *format, *language, *title, *stride, *pretty); err != nil {
		log.Fatal(err)
	}
}
