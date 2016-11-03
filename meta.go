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
	"fmt"
	"io"
	"os"
	"path"

	"github.com/FooSoft/jmdict"
)

func outputMetaIndex(outputDir string, entries []termSource, entities map[string]string, pretty bool) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputIndex, err := os.Create(path.Join(outputDir, "index.json"))
	if err != nil {
		return err
	}
	defer outputIndex.Close()

	dict := buildTermIndex(entries, entities)
	indexBytes, err := marshalJson(dict, pretty)
	if err != nil {
		return err
	}

	if _, err = outputIndex.Write(indexBytes); err != nil {
		return err
	}

	defCnt := len(dict.defs)
	for i := 0; i < defCnt; i += BANK_STRIDE {
		outputRef, err := os.Create(path.Join(outputDir, fmt.Sprintf("bank_%d.json", i/BANK_STRIDE+1)))
		if err != nil {
			return err
		}
		defer outputRef.Close()

		indexSrc := i
		indexDst := i + BANK_STRIDE
		if indexDst > defCnt {
			indexDst = defCnt
		}

		refBytes, err := marshalJson(dict.defs[indexSrc:indexDst], pretty)
		if err != nil {
			return err
		}

		if _, err = outputRef.Write(refBytes); err != nil {
			return err
		}
	}

	return nil
}

func outputMetaJson(outputDir string, reader io.Reader, flags int) error {
	dict, entities, err := jmdict.LoadJmdictNoTransform(reader)
	if err != nil {
		return err
	}

	var entries []termSource
	for _, e := range dict.Entries {
		entries = append(entries, convertEdictEntry(e)...)
	}

	return outputMetaIndex(outputDir, entries, entities, flags&flagPrettyJson == flagPrettyJson)
}
