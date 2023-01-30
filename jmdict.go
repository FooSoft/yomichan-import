package yomichan

import (
	"errors"
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

func calculateTermScore(senseNumber int, depth int, headword headword) int {
	const senseWeight int = 1
	const depthWeight int = 100
	const entryPositionWeight int = 10000
	const priorityWeight int = 1000000

	score := 0
	score -= (senseNumber - 1) * senseWeight
	score -= depth * depthWeight
	score -= headword.Index * entryPositionWeight
	score += headword.Score() * priorityWeight

	return score
}

func doDisplaySenseNumberTag(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) bool {
	// Display sense numbers if the entry has more than one sense
	// or if the headword is found in multiple entries.
	hash := headword.Hash()
	if !meta.extraMode {
		return false
	} else if meta.language != "eng" {
		return false
	} else if meta.seqToSenseCount[entry.Sequence] > 1 {
		return true
	} else if len(meta.headwordHashToSeqs[hash]) > 1 {
		return true
	} else {
		return false
	}
}

func jmdictPublicationDate(dictionary jmdict.Jmdict) string {
	if len(dictionary.Entries) == 0 {
		return "unknown"
	}
	dateEntry := dictionary.Entries[len(dictionary.Entries)-1]
	if len(dateEntry.Sense) == 0 || len(dateEntry.Sense[0].Glossary) == 0 {
		return "unknown"
	}
	r := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	jmdictDate := r.FindString(dateEntry.Sense[0].Glossary[0].Content)
	if jmdictDate != "" {
		return jmdictDate
	} else {
		return "unknown"
	}
}

func createFormsTerm(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) (dbTerm, bool) {
	// Don't add "forms" terms to non-English dictionaries.
	// Information would be duplicated if users installed more
	// than one version.
	if meta.language != "eng" || !meta.extraMode {
		return dbTerm{}, false
	}
	// Don't need a "forms" term for entries with one unique
	// headword which does not appear in any other entries.
	if !meta.hasMultipleForms[entry.Sequence] {
		if len(meta.headwordHashToSeqs[headword.Hash()]) == 1 {
			return dbTerm{}, false
		}
	}

	term := baseFormsTerm(entry)
	term.Expression = headword.Expression
	term.Reading = headword.Reading

	term.addTermTags(headword.TermTags...)

	term.addDefinitionTags("forms")
	senseNumber := meta.seqToSenseCount[entry.Sequence] + 1
	entryDepth := meta.entryDepth[entry.Sequence]
	term.Score = calculateTermScore(senseNumber, entryDepth, headword)
	return term, true
}

func createSearchTerm(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) (dbTerm, bool) {
	// Don't add "search" terms to non-English dictionaries.
	// Information would be duplicated if users installed more
	// than one version.
	if meta.language != "eng" {
		return dbTerm{}, false
	}

	term := dbTerm{
		Expression: headword.Expression,
		Sequence:   -entry.Sequence,
	}
	for _, sense := range entry.Sense {
		rules := grammarRules(sense.PartsOfSpeech)
		term.addRules(rules...)
	}
	term.addTermTags(headword.TermTags...)
	term.Score = calculateTermScore(1, 0, headword)

	redirectHeadword := meta.seqToMainHeadword[entry.Sequence]
	expHash := redirectHeadword.ExpHash()
	doDisplayReading := (len(meta.expHashToReadings[expHash]) > 1)

	content := contentSpan(
		contentAttr{fontSize: "130%"},
		"⟶",
		redirectHeadword.ToInternalLink(doDisplayReading),
	)

	term.Glossary = []any{contentStructure(content)}
	return term, true
}

func createSenseTerm(sense jmdict.JmdictSense, senseNumber int, headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) (dbTerm, bool) {
	if sense.RestrictedReadings != nil && !slices.Contains(sense.RestrictedReadings, headword.Reading) {
		return dbTerm{}, false
	}
	if sense.RestrictedKanji != nil && !slices.Contains(sense.RestrictedKanji, headword.Expression) {
		return dbTerm{}, false
	}

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

	entryDepth := meta.entryDepth[entry.Sequence]
	term.Score = calculateTermScore(senseNumber, entryDepth, headword)

	return term, true
}

func extractTerms(headword headword, entry jmdict.JmdictEntry, meta jmdictMetadata) ([]dbTerm, bool) {
	if meta.seqToSenseCount[entry.Sequence] == 0 {
		return nil, false
	}
	if headword.IsSearchOnly {
		if searchTerm, ok := createSearchTerm(headword, entry, meta); ok {
			return []dbTerm{searchTerm}, true
		} else {
			return nil, false
		}
	}
	terms := []dbTerm{}
	senseNumber := 1
	for _, sense := range entry.Sense {
		if !glossaryContainsLanguage(sense.Glossary, meta.language) {
			// Do not increment sense number
			continue
		}
		if senseTerm, ok := createSenseTerm(sense, senseNumber, headword, entry, meta); ok {
			terms = append(terms, senseTerm)
		}
		senseNumber += 1
	}

	if formsTerm, ok := createFormsTerm(headword, entry, meta); ok {
		terms = append(terms, formsTerm)
	}

	return terms, true
}

func jmdExportDb(inputPath string, outputPath string, languageName string, title string, stride int, pretty bool) error {
	if _, ok := langNameToCode[languageName]; !ok {
		return errors.New("Unrecognized language parameter: " + languageName)
	}

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

	return writeDb(
		outputPath,
		index,
		recordData,
		stride,
		pretty,
	)
}
