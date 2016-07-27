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
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s file_format input_file output_file\n\n", path.Base(os.Args[0]))
}

func process(fileFormat, inputFile, outputFile string) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	switch strings.ToLower(fileFormat) {
	case "edict":
		return processEdict(input, output)
	default:
		return errors.New("unrecognized file format")
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 3 {
		if err := process(args[0], args[1], args[2]); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
		os.Exit(2)
	}
}
