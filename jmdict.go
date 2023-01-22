package yomichan

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"foosoft.net/projects/jmdict"
	"golang.org/x/exp/slices"
)

func grammarRules(partsOfSpeech []string) []string {
	rules := []string{}
	for _, partOfSpeech := range partsOfSpeech {
		switch partOfSpeech {
		case "adj-i", "vk", "vz":
			rules = append(rules, partOfSpeech)
		default:
			if strings.HasPrefix(partOfSpeech, "v5") {
				rules = append(rules, "v5")
			} else if strings.HasPrefix(partOfSpeech, "v1") {
				rules = append(rules, "v1")
			} else if strings.HasPrefix(partOfSpeech, "vs-") {
				rules = append(rules, "vs")
			}
		}
	}
	return rules
}

func calculateTermScore(senseNumber int, headword headword) int {
	const senseWeight int = 1
	const entryPositionWeight int = 100
	const priorityWeight int = 10000

	score := 0
	score -= (senseNumber - 1) * senseWeight
	score -= headword.Index * entryPositionWeight
	score += headword.Score() * priorityWeight

	return score
}

func doDisplaySenseNumberTag(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) bool {
	// Display sense numbers if the entry has more than one sense
	// or if the headword is found in multiple entries.
	hash := headword.Hash()
	if meta.seqToSenseCount[entry.Sequence] > 1 {
		return true
	} else if len(meta.headwordHashToSeqs[hash]) > 1 {
		return true
	} else {
		return false
	}
}

func jmdictPublicationDate(dictionary jmdict.Jmdict) string {
	dateEntry := dictionary.Entries[len(dictionary.Entries)-1]
	r := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	jmdictDate := r.FindString(dateEntry.Sense[0].Glossary[0].Content)
	return jmdictDate
}

func createFormsTerm(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) dbTerm {
	term := baseFormsTerm(entry)
	term.Expression = headword.Expression
	term.Reading = headword.Reading

	term.addTermTags(headword.TermTags...)

	term.addDefinitionTags("forms")
	senseNumber := meta.seqToSenseCount[entry.Sequence] + 1
	term.Score = calculateTermScore(senseNumber, headword)
	return term
}

func createSearchTerm(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) dbTerm {
	term := dbTerm{
		Expression: headword.Expression,
		Sequence:   -entry.Sequence,
	}
	for _, sense := range entry.Sense {
		rules := grammarRules(sense.PartsOfSpeech)
		term.addRules(rules...)
	}
	term.addTermTags(headword.TermTags...)
	term.Score = calculateTermScore(0, headword)

	redirectHeadword := meta.seqToMainHeadword[entry.Sequence]
	expHash := redirectHeadword.ExpHash()
	doDisplayReading := (len(meta.expHashToReadings[expHash]) > 1)

	content := contentSpan(
		contentAttr{fontSize: "130%"},
		"⟶",
		redirectHeadword.ToInternalLink(doDisplayReading),
	)

	term.Glossary = []any{contentStructure(content)}
	return term
}

func createSenseTerm(sense jmdict.JmdictSense, senseNumber int, headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) dbTerm {
	term := dbTerm{
		Expression: headword.Expression,
		Reading:    headword.Reading,
		Sequence:   entry.Sequence,
	}

	term.Glossary = createGlossary(sense, meta)

	term.addTermTags(headword.TermTags...)

	if doDisplaySenseNumberTag(headword, entry, meta) {
		senseNumberTag := strconv.Itoa(senseNumber)
		term.addDefinitionTags(senseNumberTag)
	}
	term.addDefinitionTags(sense.PartsOfSpeech...)
	term.addDefinitionTags(sense.Fields...)
	term.addDefinitionTags(sense.Misc...)
	term.addDefinitionTags(sense.Dialects...)

	rules := grammarRules(sense.PartsOfSpeech)
	term.addRules(rules...)

	term.Score = calculateTermScore(senseNumber, headword)

	return term
}

func extractTerms(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) ([]dbTerm, bool) {
	if meta.seqToSenseCount[entry.Sequence] == 0 {
		return nil, false
	}
	if headword.IsSearchOnly {
		searchTerm := createSearchTerm(headword, entry, meta)
		return []dbTerm{searchTerm}, true
	}
	terms := []dbTerm{}
	senseNumber := 1
	for _, sense := range entry.Sense {
		if !glossaryContainsLanguage(sense.Glossary, meta.language) {
			continue
		}
		if sense.RestrictedReadings != nil && !slices.Contains(sense.RestrictedReadings, headword.Reading) {
			senseNumber += 1
			continue
		}
		if sense.RestrictedKanji != nil && !slices.Contains(sense.RestrictedKanji, headword.Expression) {
			senseNumber += 1
			continue
		}
		senseTerm := createSenseTerm(sense, senseNumber, headword, entry, meta)
		senseNumber += 1
		terms = append(terms, senseTerm)
	}

	if meta.hasMultipleForms[entry.Sequence] {
		formsTerm := createFormsTerm(headword, entry, meta)
		terms = append(terms, formsTerm)
	}
	return terms, true
}

func jmdExportDb(inputPath string, outputPath string, languageName string, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dictionary, entities, err := jmdict.LoadJmdictNoTransform(reader)
	if err != nil {
		return err
	}

	meta := newJmdictMetadata(dictionary, languageName)

	terms := dbTermList{}
	for _, entry := range dictionary.Entries {
		headwords := extractHeadwords(entry)
		for _, headword := range headwords {
			if newTerms, ok := extractTerms(headword, entry, meta); ok {
				terms = append(terms, newTerms...)
			}
		}
	}

	tags := dbTagList{}
	tags = append(tags, entityTags(entities)...)
	tags = append(tags, senseNumberTags(meta.maxSenseCount)...)
	tags = append(tags, newsFrequencyTags()...)
	tags = append(tags, customDbTags()...)

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  tags.crush(),
	}

	if title == "" {
		title = "JMdict"
	}
	jmdictDate := jmdictPublicationDate(dictionary)

	index := dbIndex{
		Title:       title,
		Revision:    "JMdict." + jmdictDate,
		Sequenced:   true,
		Attribution: edrdgAttribution,
	}
	index.setDefaults()

	return writeDb(
		outputPath,
		index,
		recordData,
		stride,
		pretty,
	)
}
