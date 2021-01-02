/*
 * Copyright (c) 2017-2021 Alex Yatskov <alex@foosoft.net>
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

package yomichan

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const frequencyRevision = "frequency1"

func frequencyTermsExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	return frequncyExportDb(inputPath, outputPath, language, title, stride, pretty, "term_meta")
}

func frequencyKanjiExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	return frequncyExportDb(inputPath, outputPath, language, title, stride, pretty, "kanji_meta")
}

func frequncyExportDb(inputPath, outputPath, language, title string, stride int, pretty bool, key string) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	var frequencies dbMetaList
	for scanner := bufio.NewScanner(reader); scanner.Scan(); {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		expression := parts[0]
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			expression = parts[1]
			count, err = strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
		}

		frequencies = append(frequencies, dbMeta{expression, "freq", count})
	}

	if title == "" {
		title = "Frequency"
	}

	recordData := map[string]dbRecordList{
		key: frequencies.crush(),
	}

	return writeDb(
		outputPath,
		title,
		frequencyRevision,
		false,
		recordData,
		stride,
		pretty,
	)
}
