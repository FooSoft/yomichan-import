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
)

type termMetaEntry struct {
	Text string `json:"text"`
	Tags string `json:"tags"`
}

type termMetaIndex struct {
	Version  int               `json:"version"`
	Banks    int               `json:"banks"`
	Entities map[string]string `json:"entities"`
	Entries  []termMetaEntry   `json:"entries"`
}

func newTermMetaIndex(entries []termMetaEntry, entities map[string]string) termMetaIndex {
	return termMetaIndex{
		Version:  DB_VERSION,
		Banks:    bankCount(len(entries)),
		Entities: entities,
		Entries:  entries,
	}
}

func (index *termMetaIndex) output(dir string, pretty bool) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	bytes, err := marshalJson(index, pretty)
	if err != nil {
		return err
	}

	fp, err := os.Create(path.Join(dir, "index.json"))
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err = fp.Write(bytes); err != nil {
		return err
	}

	count := len(index.Entries)

	for i := 0; i < count; i += BANK_STRIDE {
		indexSrc := i
		indexDst := i + BANK_STRIDE
		if indexDst > count {
			indexDst = count
		}

		bytes, err := marshalJson(index.Entries[indexSrc:indexDst], pretty)
		if err != nil {
			return err
		}

		fp, err := os.Create(path.Join(dir, fmt.Sprintf("bank_%d.json", i/BANK_STRIDE+1)))
		if err != nil {
			return err
		}
		defer fp.Close()

		if _, err = fp.Write(bytes); err != nil {
			return err
		}
	}

	return nil
}

func outputTermMetaJson(dir string, reader io.Reader, flags int) error {
	// dict, entities, err := jmdict.LoadJmdictNoTransform(reader)
	// if err != nil {
	// 	return err
	// }

	// meta := make(map[string][]string)
	// for _, entry := range dict.Entries {
	// }

	// var entries []termMetaEntry
	// for _, entry := range dict.Entries {
	// 	// defs = append(defs, convertEdictEntry(e)...)
	// }

	// index := newTermMetaIndex(entries, entities)
	// return index.output(dir, flags&flagPrettyJson == flagPrettyJson)

	return nil
}
