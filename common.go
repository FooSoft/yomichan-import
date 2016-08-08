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
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

type vocabJson struct {
	Indices  map[string]string `json:"i"`
	Entities map[string]string `json:"e"`
	Defs     [][]string        `json:"d"`
}

type vocabSource struct {
	Expression string
	Reading    string
	Tags       []string
	Glossary   []string
}

func (s *vocabSource) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, s.Tags) {
			s.Tags = append(s.Tags, tag)
		}
	}
}

func (s *vocabSource) addTagsPri(tags ...string) {
	for _, tag := range tags {
		switch tag {
		case "news1", "ichi1", "spec1", "gai1":
			s.addTags("P")
			fallthrough
		case "news2", "ichi2", "spec2", "gai2":
			s.addTags(tag[:len(tag)-1])
			break
		}
	}
}

func buildVocabJson(entries []vocabSource, entities map[string]string) vocabJson {
	dict := vocabJson{
		Indices:  make(map[string]string),
		Entities: entities,
	}

	for i, e := range entries {
		entry := []string{e.Expression, e.Reading, strings.Join(e.Tags, " ")}
		entry = append(entry, e.Glossary...)
		dict.Defs = append(dict.Defs, entry)

		appendStrIndex(dict.Indices, e.Expression, i)
		if len(e.Reading) > 0 {
			appendStrIndex(dict.Indices, e.Reading, i)
		}
	}

	return dict
}

func outputVocabJson(writer io.Writer, entries []vocabSource, entities map[string]string, pretty bool) error {
	dict := buildVocabJson(entries, entities)

	var (
		bytes []byte
		err   error
	)

	if pretty {
		bytes, err = json.MarshalIndent(dict, "", "    ")
	} else {
		bytes, err = json.Marshal(dict)
	}

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}

func appendStrIndex(indices map[string]string, key string, value int) {
	def, _ := indices[key]
	if len(def) > 0 {
		def += " "
	}

	def += strconv.Itoa(value)
	indices[key] = def
}

func hasString(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}
