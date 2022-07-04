package yomichan

import (
	"os"

	"foosoft.net/projects/jmdict"
)

const jmnedictRevision = "jmnedict1"

func jmnedictBuildTagMeta(entities map[string]string) dbTagList {
	var tags dbTagList

	for name, value := range entities {
		tag := dbTag{Name: name, Notes: value}

		switch name {
		case "company", "fem", "given", "masc", "organization", "person", "place", "product", "station", "surname", "unclass", "work":
			tag.Category = "name"
			tag.Order = 4
		}

		tags = append(tags, tag)
	}

	return tags
}

func jmnedictExtractTerms(enamdictEntry jmdict.JmnedictEntry) []dbTerm {
	var terms []dbTerm

	convert := func(reading jmdict.JmnedictReading, kanji *jmdict.JmnedictKanji) {
		if kanji != nil && hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var term dbTerm
		term.Sequence = enamdictEntry.Sequence
		term.addTermTags(reading.Information...)

		if kanji == nil {
			term.Expression = reading.Reading
		} else {
			term.Expression = kanji.Expression
			term.Reading = reading.Reading
			term.addTermTags(kanji.Information...)

			for _, priority := range kanji.Priorities {
				if hasString(priority, reading.Priorities) {
					term.addTermTags(priority)
				}
			}
		}

		for _, trans := range enamdictEntry.Translations {
			term.Glossary = append(term.Glossary, trans.Translations...)
			term.addDefinitionTags(trans.NameTypes...)
		}

		terms = append(terms, term)
	}

	if len(enamdictEntry.Kanji) > 0 {
		for _, kanji := range enamdictEntry.Kanji {
			for _, reading := range enamdictEntry.Readings {
				convert(reading, &kanji)
			}
		}
	} else {
		for _, reading := range enamdictEntry.Readings {
			convert(reading, nil)
		}
	}

	return terms
}

func jmnedictExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dict, entities, err := jmdict.LoadJmnedictNoTransform(reader)
	if err != nil {
		return err
	}

	var terms dbTermList
	for _, entry := range dict.Entries {
		terms = append(terms, jmnedictExtractTerms(entry)...)
	}

	if title == "" {
		title = "JMnedict"
	}

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  jmnedictBuildTagMeta(entities).crush(),
	}

	return writeDb(
		outputPath,
		title,
		jmnedictRevision,
		true,
		recordData,
		stride,
		pretty,
	)
}
