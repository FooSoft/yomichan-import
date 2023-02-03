package yomichan

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func frequencyTermsExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	return frequencyExportDb(inputPath, outputPath, language, title, stride, pretty, "term_meta")
}

func frequencyKanjiExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	return frequencyExportDb(inputPath, outputPath, language, title, stride, pretty, "kanji_meta")
}

func frequencyExportDb(inputPath, outputPath, language, title string, stride int, pretty bool, key string) error {
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

	index := dbIndex{
		Title:     title,
		Revision:  "frequency1",
		Sequenced: false,
	}

	return writeDb(
		outputPath,
		index,
		recordData,
		stride,
		pretty,
	)
}
