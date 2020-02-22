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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type epwingEntry struct {
	Heading string `json:"heading"`
	Text    string `json:"text"`
}

type epwingSubbook struct {
	Title     string        `json:"title"`
	Copyright string        `json:"copyright"`
	Entries   []epwingEntry `json:"entries"`
}

type epwingBook struct {
	CharCode string          `json:"charCode"`
	DiscCode string          `json:"discCode"`
	Subbooks []epwingSubbook `json:"subbooks"`
}

type epwingExtractor interface {
	extractTerms(entry epwingEntry, sequence int) []dbTerm
	extractKanji(entry epwingEntry) []dbKanji
	getFontNarrow() map[int]string
	getFontWide() map[int]string
	getRevision() string
}

func epwingExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	var toolExec bool
	if stat.IsDir() {
		toolExec = true
	} else if filepath.Base(inputPath) == "CATALOGS" {
		inputPath = filepath.Dir(inputPath)
		toolExec = true
	}

	var data []byte
	if toolExec {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}

		toolPath := filepath.Join("bin", runtime.GOOS, "zero-epwing")
		if runtime.GOOS == "windows" {
			toolPath += ".exe"
		}

		toolPath = filepath.Join(filepath.Dir(ex), toolPath)

		if _, err = os.Stat(toolPath); err != nil {
			return fmt.Errorf("failed to find zero-epwing in '%s'", toolPath)
		}

		cmd := exec.Command(toolPath, "--entries", inputPath)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		log.Printf("invoking zero-epwing from '%s'...\n", toolPath)
		if err := cmd.Start(); err != nil {
			return err
		}

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				log.Printf("\t > %s\n", scanner.Text())
			}
		}()

		if data, err = ioutil.ReadAll(stdout); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		log.Println("completed zero-epwing processing")
	} else {
		data, err = ioutil.ReadFile(inputPath)
	}

	if err != nil {
		return err
	}

	var book epwingBook
	if err := json.Unmarshal(data, &book); err != nil {
		return err
	}

	translateExp := regexp.MustCompile(`{{([nw])_(\d+)}}`)
	epwingExtractors := map[string]epwingExtractor{
		"三省堂　スーパー大辞林":    makeDaijirinExtractor(),
		"大辞泉":            makeDaijisenExtractor(),
		"明鏡国語辞典":         makeMeikyouExtractor(),
		"故事ことわざの辞典":      makeKotowazaExtractor(),
		"研究社　新和英大辞典　第５版": makeWadaiExtractor(),
		"広辞苑第六版":         makeKoujienExtractor(),
		"付属資料":           makeKoujienExtractor(),
		"学研国語大辞典":        makeGakkenExtractor(),
		"古語辞典":           makeGakkenExtractor(),
		"故事ことわざ辞典":       makeGakkenExtractor(),
		"学研漢和大字典":        makeGakkenExtractor(),
	}

	var (
		terms     dbTermList
		kanji     dbKanjiList
		revisions []string
		titles    []string
	)

	log.Println("formatting dictionary data...")

	var sequence int

	for _, subbook := range book.Subbooks {
		if extractor, ok := epwingExtractors[subbook.Title]; ok {
			fontNarrow := extractor.getFontNarrow()
			fontWide := extractor.getFontWide()

			translate := func(str string) string {
				for _, matches := range translateExp.FindAllStringSubmatch(str, -1) {
					var font map[int]string
					if matches[1] == "n" {
						font = fontNarrow
					} else {
						font = fontWide
					}

					code, _ := strconv.Atoi(matches[2])
					replacement, ok := font[code]
					if !ok {
						replacement = "�"
					}

					str = strings.Replace(str, matches[0], replacement, -1)
				}

				return str
			}

			for _, entry := range subbook.Entries {
				entry.Heading = translate(entry.Heading)
				entry.Text = translate(entry.Text)

				terms = append(terms, extractor.extractTerms(entry, sequence)...)
				kanji = append(kanji, extractor.extractKanji(entry)...)

				sequence++
			}

			revisions = append(revisions, extractor.getRevision())
			titles = append(titles, subbook.Title)
		} else {
			return fmt.Errorf("failed to find compatible extractor for '%s'", subbook.Title)
		}
	}

	if title == "" {
		title = strings.Join(titles, ", ")
	}

	recordData := map[string]dbRecordList{
		"kanji": kanji.crush(),
		"term":  terms.crush(),
	}

	return writeDb(
		outputPath,
		title,
		strings.Join(revisions, ";"),
		true,
		recordData,
		stride,
		pretty,
	)
}
