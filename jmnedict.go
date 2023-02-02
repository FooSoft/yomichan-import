package yomichan

import (
	"os"
	"regexp"

	"foosoft.net/projects/jmdict"
)

func jmnedictPublicationDate(dictionary jmdict.Jmnedict) string {
	if len(dictionary.Entries) == 0 {
		return "unknown"
	}
	dateEntry := dictionary.Entries[len(dictionary.Entries)-1]
	if len(dateEntry.Translations) == 0 || len(dateEntry.Translations[0].Translations) == 0 {
		return "unknown"
	}
	r := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	jmnedictDate := r.FindString(dateEntry.Translations[0].Translations[0])
	if jmnedictDate != "" {
		return jmnedictDate
	} else {
		return "unknown"
	}
}

func jmnedictSenseTerm(headword headword, seq sequence, sense jmdict.JmnedictTranslation, senseNumber int) dbTerm {
	term := dbTerm{
		Expression: headword.Expression,
		Reading:    headword.Reading,
		Sequence:   seq,
	}
	for _, gloss := range sense.Translations {
		term.Glossary = append(term.Glossary, gloss)
	}
	term.addDefinitionTags(sense.NameTypes...)
	term.Score = calculateTermScore(senseNumber, 0, headword)
	return term
}

func jmnedictTerms(headword headword, entry jmdict.JmnedictEntry, g genericTermInfo) []dbTerm {
	terms := []dbTerm{}
	for idx, sense := range entry.Translations {
		if g.IsGenericName(headword, sense.Translations) {
			g.AddGlosses(headword.Expression, sense.NameTypes, headword.Reading)
		} else {
			g.AddUsedSequence(entry.Sequence)
			senseTerm := jmnedictSenseTerm(headword, entry.Sequence, sense, idx+1)
			terms = append(terms, senseTerm)
		}
	}
	return terms
}

func jmnedictHeadwords(entry jmdict.JmnedictEntry) (headwords []headword) {
	// Note that JMnedict doesn't (currently) use priority tags,
	// frequency tags, or any sort of reading/kanji restrictions.
	for _, reading := range entry.Readings {
		for _, kanji := range entry.Kanji {
			h := headword{
				Expression: kanji.Expression,
				Reading:    reading.Reading,
			}
			h.Index = len(headwords)
			headwords = append(headwords, h)
		}
	}
	if len(entry.Kanji) == 0 {
		for _, reading := range entry.Readings {
			h := headword{
				Expression: reading.Reading,
				Reading:    reading.Reading,
			}
			h.Index = len(headwords)
			headwords = append(headwords, h)
		}
	}
	return headwords
}

func jmnedictExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dictionary, entities, err := jmdict.LoadJmnedictNoTransform(reader)
	if err != nil {
		return err
	}

	genericTermInfo := newGenericTermInfo()

	terms := dbTermList{}
	for _, entry := range dictionary.Entries {
		headwords := jmnedictHeadwords(entry)
		for _, headword := range headwords {
			newTerms := jmnedictTerms(headword, entry, genericTermInfo)
			terms = append(terms, newTerms...)
		}
	}
	terms = append(terms, genericTermInfo.Terms()...)

	tags := dbTagList{}
	tags = append(tags, entityTags(entities)...)

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  tags.crush(),
	}

	if title == "" {
		title = "JMnedict"
	}
	jmnedictDate := jmnedictPublicationDate(dictionary)

	index := dbIndex{
		Title:       title,
		Revision:    "JMnedict." + jmnedictDate,
		Sequenced:   true,
		Attribution: edrdgAttribution,
	}

	return writeDb(
		outputPath,
		index,
		recordData,
		stride,
		pretty,
	)
}
