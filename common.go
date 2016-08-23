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
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	REF_STEP_COUNT = 1000
)

type termJson struct {
	Refs     int        `json:"refs"`
	Entities [][]string `json:"ents"`
	defs     [][]string
}

type termSource struct {
	Expression string
	Reading    string
	Tags       []string
	Glossary   []string
}

func (s *termSource) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, s.Tags) {
			s.Tags = append(s.Tags, tag)
		}
	}
}

func (s *termSource) addTagsPri(tags ...string) {
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

func buildTermJson(entries []termSource, entities map[string]string) termJson {
	var dict termJson

	for name, value := range entities {
		ent := []string{name, value}
		dict.Entities = append(dict.Entities, ent)
	}

	for _, e := range entries {
		def := []string{e.Expression, e.Reading, strings.Join(e.Tags, " ")}
		def = append(def, e.Glossary...)
		dict.defs = append(dict.defs, def)
	}

	dict.Refs = len(dict.defs) / REF_STEP_COUNT

	return dict
}

func marshalJson(obj interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(obj, "", "    ")
	}

	return json.Marshal(obj)
}

func outputTermJson(outputDir string, entries []termSource, entities map[string]string, pretty bool) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputIndex, err := os.Create(path.Join(outputDir, "index.json"))
	if err != nil {
		return err
	}
	defer outputIndex.Close()

	dict := buildTermJson(entries, entities)

	indexBytes, err := marshalJson(dict, pretty)
	if err != nil {
		return err
	}

	if _, err = outputIndex.Write(indexBytes); err != nil {
		return err
	}

	defCnt := len(dict.defs)

	for i := 0; i < defCnt; i += REF_STEP_COUNT {
		outputRef, err := os.Create(path.Join(outputDir, fmt.Sprintf("ref_%d.json", i/REF_STEP_COUNT)))
		if err != nil {
			return err
		}
		defer outputRef.Close()

		indexSrc := i
		indexDst := i + REF_STEP_COUNT
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

func hasString(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}
