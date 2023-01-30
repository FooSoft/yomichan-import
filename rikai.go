package yomichan

import (
	"database/sql"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type rikaiEntry struct {
	kanji string
	kana  string
	entry string
}

func rikaiBuildRules(term *dbTerm) {
	for _, tag := range term.DefinitionTags {
		switch tag {
		case "adj-i", "v1", "vk", "vz":
			term.addRules(tag)
		default:
			if strings.HasPrefix(tag, "v5") {
				term.addRules("v5")
			} else if strings.HasPrefix(tag, "vs") {
				term.addRules("vs")
			}
		}
	}
}

func rikaiBuildScore(term *dbTerm) {
	for _, tag := range term.DefinitionTags {
		switch tag {
		case "news", "ichi", "spec", "gai":
			term.Score++
		case "arch", "iK":
			term.Score--
		case "P":
			term.Score += 5
		}
	}
}

func rikaiExtractTerms(rows *sql.Rows) (dbTermList, error) {
	var terms dbTermList

	dfnExp := regexp.MustCompile(`^(?:＊\(KC\) )?((?:\((?:[\w\-\,\:]*)*\)\s*)*)(.*)$`)
	readExp := regexp.MustCompile(`\[([^\]]+)\]`)
	tagExp := regexp.MustCompile(`[\s\(\),]`)

	var sequence int

	for rows.Next() {
		var (
			kanji, kana, entry *string
		)

		if err := rows.Scan(&kanji, &kana, &entry); err != nil {
			return nil, err
		}

		if entry == nil {
			continue
		}

		segments := strings.Split(*entry, "/")
		indexFirst := 0

		if !strings.HasPrefix(segments[0], "＊") && !strings.HasPrefix(segments[0], "(") {
			expParts := strings.Split(segments[0], " ")
			if len(expParts) > 1 {
				if readMatch := readExp.FindStringSubmatch(expParts[1]); readMatch != nil {
					kana = &readMatch[1]
				}
			}

			indexFirst = 1
			if indexFirst == len(segments) {
				continue
			}
		}

		var term dbTerm
		term.Sequence = sequence
		if kana != nil {
			term.Expression = *kana
			term.Reading = *kana
		}
		if kanji != nil {
			term.Expression = *kanji
		}

		for i := indexFirst; i < len(segments); i++ {
			segment := segments[i]

			if dfnMatch := dfnExp.FindStringSubmatch(segment); dfnMatch != nil {
				for _, tag := range tagExp.Split(dfnMatch[1], -1) {
					if rikaiTagParsed(tag) {
						term.addDefinitionTags(tag)
					}
				}

				if len(dfnMatch[2]) > 0 {
					term.Glossary = append(term.Glossary, dfnMatch[2])
				}
			}
		}

		rikaiBuildRules(&term)
		rikaiBuildScore(&term)

		terms = append(terms, term)

		sequence++
	}

	return terms, nil
}

func rikaiExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	db, err := sql.Open("sqlite3", inputPath)
	if err != nil {
		return err
	}
	defer db.Close()

	dictRows, err := db.Query("SELECT kanji, kana, entry FROM dict")
	if err != nil {
		return err
	}

	terms, err := rikaiExtractTerms(dictRows)
	if err != nil {
		return err
	}

	if title == "" {
		title = "Rikai"
	}

	tags := dbTagList{
		dbTag{Name: "P", Category: "popular", Order: -10},
		dbTag{Name: "exp", Category: "expression", Order: -5},
		dbTag{Name: "id", Category: "expression", Order: -5},
		dbTag{Name: "arch", Category: "archaism", Order: -4},
		dbTag{Name: "iK", Category: "archaism", Order: -4},
	}

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  tags.crush(),
	}

	index := dbIndex{
		Title:     title,
		Revision:  "rikai2",
		Sequenced: true,
	}

	return writeDb(
		outputPath,
		index,
		recordData,
		stride,
		pretty,
	)
}

func rikaiTagParsed(tag string) bool {
	tags := []string{
		"Buddh",
		"MA",
		"X",
		"abbr",
		"adj",
		"adj-f",
		"adj-i",
		"adj-na",
		"adj-no",
		"adj-pn",
		"adj-t",
		"adv",
		"adv-n",
		"adv-to",
		"arch",
		"ateji",
		"aux",
		"aux-adj",
		"aux-v",
		"c",
		"chn",
		"col",
		"comp",
		"conj",
		"ctr",
		"derog",
		"eK",
		"ek",
		"exp",
		"f",
		"fam",
		"fem",
		"food",
		"g",
		"geom",
		"gikun",
		"gram",
		"h",
		"hon",
		"hum",
		"iK",
		"id",
		"ik",
		"int",
		"io",
		"iv",
		"ling",
		"m",
		"m-sl",
		"male",
		"male-sl",
		"math",
		"mil",
		"n",
		"n-adv",
		"n-pref",
		"n-suf",
		"n-t",
		"num",
		"oK",
		"obs",
		"obsc",
		"ok",
		"on-mim",
		"P",
		"p",
		"physics",
		"pn",
		"poet",
		"pol",
		"pr",
		"pref",
		"prt",
		"rare",
		"s",
		"sens",
		"sl",
		"st",
		"suf",
		"u",
		"uK",
		"uk",
		"v1",
		"v2a-s",
		"v4h",
		"v4r",
		"v5",
		"v5aru",
		"v5b",
		"v5g",
		"v5k",
		"v5k-s",
		"v5m",
		"v5n",
		"v5r",
		"v5r-i",
		"v5s",
		"v5t",
		"v5u",
		"v5u-s",
		"v5uru",
		"v5z",
		"vi",
		"vk",
		"vn",
		"vs",
		"vs-c",
		"vs-i",
		"vs-s",
		"vt",
		"vulg",
		"vz",
	}

	for _, curr := range tags {
		if curr == tag {
			return true
		}
	}

	return false
}
