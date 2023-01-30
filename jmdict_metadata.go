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

func (meta *jmdictMetadata) CalculateEntryDepth(headwords []headword, entrySequence sequence) {
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
	meta.entryDepth[entrySequence] = maxDepth
}

func (meta *jmdictMetadata) AddHeadword(headword headword, entry jmdict.JmdictEntry) {

	// Determine how many senses are in this entry for this language
	if _, ok := meta.seqToSenseCount[entry.Sequence]; !ok {
		senseCount := 0
		for _, entrySense := range entry.Sense {
			for _, gloss := range entrySense.Glossary {
				if glossContainsLanguage(gloss, meta.language) {
					senseCount += 1
					break
				}
			}
		}
		meta.seqToSenseCount[entry.Sequence] = senseCount
	}

	if meta.seqToSenseCount[entry.Sequence] == 0 {
		return
	}

	// main headwords (first ones that are found in entries).
	if _, ok := meta.seqToMainHeadword[entry.Sequence]; !ok {
		meta.seqToMainHeadword[entry.Sequence] = headword
	}

	// hash the term pair so we can determine if it's used
	// in more than one JMdict entry later.
	headwordHash := headword.Hash()
	if !slices.Contains(meta.headwordHashToSeqs[headwordHash], entry.Sequence) {
		meta.headwordHashToSeqs[headwordHash] = append(meta.headwordHashToSeqs[headwordHash], entry.Sequence)
	}

	// hash the expression so that we can determine if we
	// need to disambiguate it by displaying its reading
	// in reference notes later.
	expHash := headword.ExpHash()
	if !slices.Contains(meta.expHashToReadings[expHash], headword.Reading) {
		meta.expHashToReadings[expHash] = append(meta.expHashToReadings[expHash], headword.Reading)
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
		if !slices.Contains(meta.seqToSearchHashes[entry.Sequence], x) {
			meta.seqToSearchHashes[entry.Sequence] = append(meta.seqToSearchHashes[entry.Sequence], x)
		}
	}

	currentSenseNumber := 1
	for _, entrySense := range entry.Sense {
		if !glossaryContainsLanguage(entrySense.Glossary, meta.language) {
			continue
		}
		if entrySense.RestrictedReadings != nil && !slices.Contains(entrySense.RestrictedReadings, headword.Reading) {
			currentSenseNumber += 1
			continue
		}
		if entrySense.RestrictedKanji != nil && !slices.Contains(entrySense.RestrictedKanji, headword.Expression) {
			currentSenseNumber += 1
			continue
		}

		allReferences := append(entrySense.References, entrySense.Antonyms...)
		for _, reference := range allReferences {
			meta.references = append(meta.references, reference)
		}

		currentSense := senseID{entry.Sequence, currentSenseNumber}
		if meta.condensedGlosses[currentSense] == "" {
			glosses := []string{}
			for _, gloss := range entrySense.Glossary {
				if glossContainsLanguage(gloss, meta.language) && gloss.Type == nil {
					glosses = append(glosses, gloss.Content)
				}
			}
			meta.condensedGlosses[currentSense] = strings.Join(glosses, "; ")
		}
		currentSenseNumber += 1
	}
}

func newJmdictMetadata(dictionary jmdict.Jmdict, languageName string) jmdictMetadata {
	meta := jmdictMetadata{
		language:           langNameToCode[languageName],
		seqToSenseCount:    make(map[sequence]int),
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
		headwords := extractHeadwords(entry)
		formCount := 0
		for _, headword := range headwords {
			meta.AddHeadword(headword, entry)
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
