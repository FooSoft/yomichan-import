package yomichan

import (
	"strings"

	"foosoft.net/projects/jmdict"
	"golang.org/x/exp/slices"
)

type sequence = int

type jmdictMetadata struct {
	language           string
	condensedGlosses   map[senseID]string
	seqToSenseCount    map[sequence]int
	seqToPartsOfSpeech map[sequence][]string
	seqToMainHeadword  map[sequence]headword
	expHashToReadings  map[hash][]string
	headwordHashToSeqs map[hash][]sequence
	references         []string
	referenceToSeq     map[string]sequence
	hashToSearchValues map[hash][]searchValue
	seqToSearchHashes  map[sequence][]searchHash
	entryDepth         map[sequence]int
	hasMultipleForms   map[sequence]bool
	maxSenseCount      int
	extraMode          bool
}

type senseID struct {
	sequence sequence
	number   int
}

func (meta *jmdictMetadata) CalculateEntryDepth(headwords []headword, seq sequence) {
	// This is to ensure that terms are grouped among their
	// entries of origin and displayed in correct sequential order
	maxDepth := 0
	for _, headword := range headwords {
		hash := headword.Hash()
		for _, seq := range meta.headwordHashToSeqs[hash] {
			seqDepth := meta.entryDepth[seq]
			if seqDepth == 0 {
				meta.entryDepth[seq] = 1
				seqDepth = 1
			}
			if maxDepth < seqDepth+1 {
				maxDepth = seqDepth + 1
			}
		}
	}
	meta.entryDepth[seq] = maxDepth
}

func (meta *jmdictMetadata) AddEntry(entry jmdict.JmdictEntry) {
	partsOfSpeech := []string{}
	senseCount := 0
	for _, sense := range entry.Sense {
		// Only English-language senses contain part-of-speech info,
		// but other languages need them for deinflection rules.
		for _, pos := range sense.PartsOfSpeech {
			if !slices.Contains(partsOfSpeech, pos) {
				partsOfSpeech = append(partsOfSpeech, pos)
			}
		}

		if glossaryContainsLanguage(sense.Glossary, meta.language) {
			senseCount += 1
		} else {
			continue
		}

		for _, reference := range sense.References {
			meta.references = append(meta.references, reference)
		}
		for _, antonym := range sense.Antonyms {
			meta.references = append(meta.references, antonym)
		}

		currentSenseID := senseID{entry.Sequence, senseCount}
		glosses := []string{}
		for _, gloss := range sense.Glossary {
			if glossContainsLanguage(gloss, meta.language) && gloss.Type == nil {
				glosses = append(glosses, gloss.Content)
			}
		}
		meta.condensedGlosses[currentSenseID] = strings.Join(glosses, "; ")
	}
	meta.seqToPartsOfSpeech[entry.Sequence] = partsOfSpeech
	meta.seqToSenseCount[entry.Sequence] = senseCount
}

func (meta *jmdictMetadata) AddHeadword(headword headword, seq sequence) {
	if meta.seqToSenseCount[seq] == 0 {
		return
	}

	// main headwords (first ones that are found in entries).
	if _, ok := meta.seqToMainHeadword[seq]; !ok {
		meta.seqToMainHeadword[seq] = headword
	}

	// hash the term pair so we can determine if it's used
	// in more than one JMdict entry later.
	headwordHash := headword.Hash()
	if !slices.Contains(meta.headwordHashToSeqs[headwordHash], seq) {
		meta.headwordHashToSeqs[headwordHash] =
			append(meta.headwordHashToSeqs[headwordHash], seq)
	}

	// hash the expression so that we can determine if we
	// need to disambiguate it by displaying its reading
	// in reference notes later.
	expHash := headword.ExpHash()
	if !slices.Contains(meta.expHashToReadings[expHash], headword.Reading) {
		meta.expHashToReadings[expHash] =
			append(meta.expHashToReadings[expHash], headword.Reading)
	}

	// e.g. for JMdict (English) we expect to end up with
	// seqToHashedHeadwords[1260670] == 【元・もと】、【元・元】、【もと・もと】、【本・もと】、【本・本】、【素・もと】、【素・素】、【基・もと】、【基・基】
	// used for correlating references to sequence numbers later.
	searchHashes := []searchHash{
		searchHash{headwordHash, headword.IsPriority},
		searchHash{expHash, headword.IsPriority},
		searchHash{headword.ReadingHash(), headword.IsPriority},
	}
	for _, x := range searchHashes {
		if !slices.Contains(meta.seqToSearchHashes[seq], x) {
			meta.seqToSearchHashes[seq] = append(meta.seqToSearchHashes[seq], x)
		}
	}
}

func newJmdictMetadata(dictionary jmdict.Jmdict, languageName string) jmdictMetadata {
	meta := jmdictMetadata{
		language:           langNameToCode[languageName],
		seqToSenseCount:    make(map[sequence]int),
		seqToPartsOfSpeech: make(map[sequence][]string),
		condensedGlosses:   make(map[senseID]string),
		seqToMainHeadword:  make(map[sequence]headword),
		expHashToReadings:  make(map[hash][]string),
		seqToSearchHashes:  make(map[sequence][]searchHash),
		headwordHashToSeqs: make(map[hash][]sequence),
		references:         []string{},
		hashToSearchValues: nil,
		referenceToSeq:     nil,
		entryDepth:         make(map[sequence]int),
		hasMultipleForms:   make(map[sequence]bool),
		maxSenseCount:      0,
		extraMode:          languageName == "english_extra",
	}

	for _, entry := range dictionary.Entries {
		meta.AddEntry(entry)
		headwords := extractHeadwords(entry)
		formCount := 0
		for _, headword := range headwords {
			meta.AddHeadword(headword, entry.Sequence)
			if !headword.IsSearchOnly {
				formCount += 1
			}
		}
		meta.CalculateEntryDepth(headwords, entry.Sequence)
		meta.hasMultipleForms[entry.Sequence] = (formCount > 1)
	}

	// this correlation process will be unnecessary once JMdict
	// includes sequence numbers in its cross-reference data
	meta.MakeReferenceToSeqMap()

	for _, senseCount := range meta.seqToSenseCount {
		if meta.maxSenseCount < senseCount {
			meta.maxSenseCount = senseCount
		}
	}

	return meta
}
