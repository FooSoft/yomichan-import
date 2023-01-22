package yomichan

import (
	"fmt"
	"strconv"
	"strings"
)

/*
 * In the future, JMdict will be updated to include sequence numbers
 * with each cross reference. At that time, most of the functions and
 * types defined in this file will become unnecessary.  see:
 * https://www.edrdg.org/jmdict_edict_list/2022/msg00008.html
 */

type searchValue struct {
	sequence   sequence
	index      int
	isPriority bool
}

type searchHash struct {
	hash       hash
	isPriority bool
}

func parseReference(reference string) (headword, int, bool) {
	// Reference strings in JMDict currently consist of 3 parts at
	// most, separated by ・ characters. The latter two parts are
	// optional.  When the sense number is not specified, it is
	// implied to be the first sense.
	var h headword
	var senseNumber int
	ok := true
	refParts := strings.Split(reference, "・")
	if len(refParts) == 1 {
		// (Kanji) or (Reading)
		h = headword{Expression: refParts[0], Reading: refParts[0]}
		senseNumber = 1
	} else if len(refParts) == 2 {
		// [Kanji + (Reading or Sense)] or (Reading + Sense)
		val, err := strconv.Atoi(refParts[1])
		if err == nil {
			h = headword{Expression: refParts[0], Reading: refParts[0]}
			senseNumber = val
		} else {
			h = headword{Expression: refParts[0], Reading: refParts[1]}
			senseNumber = 1
		}
	} else if len(refParts) == 3 {
		// Expression + Reading + Sense
		h = headword{Expression: refParts[0], Reading: refParts[1]}
		val, err := strconv.Atoi(strings.TrimSpace(refParts[2]))
		if err == nil {
			senseNumber = val
		} else {
			errortext := "Unexpected format (3rd part not integer) for x-ref \"" + reference + "\""
			fmt.Println(errortext)
			ok = false
		}
	} else {
		errortext := "Unexpected format for x-ref \"" + reference + "\""
		fmt.Println(errortext)
		ok = false
	}
	return h, senseNumber, ok
}

func (meta *jmdictMetadata) MakeReferenceToSeqMap() {

	meta.referenceToSeq = make(map[string]sequence)
	meta.MakeHashToSearchValuesMap()

	for _, reference := range meta.references {
		if meta.referenceToSeq[reference] != 0 {
			continue
		}
		seq := meta.FindBestSequence(reference)
		if seq != 0 {
			meta.referenceToSeq[reference] = seq
		} else {
			fmt.Println("Unable to convert reference to sequence number: `" + reference + "`")
		}
	}
}

func (meta *jmdictMetadata) MakeHashToSearchValuesMap() {
	meta.hashToSearchValues = make(map[hash][]searchValue)
	for seq, searchHashes := range meta.seqToSearchHashes {
		for score, searchHash := range searchHashes {
			searchValue := searchValue{
				sequence:   seq,
				index:      score,
				isPriority: searchHash.isPriority,
			}
			meta.hashToSearchValues[searchHash.hash] =
				append(meta.hashToSearchValues[searchHash.hash], searchValue)
		}
	}
}

/*
 * Generally, correspondence is determined by the order in which term
 * pairs are extracted from each JMdict entry. Take for example the
 * JMdict entry for ご本, which contains a reference to 本 (without a
 * reading specified). To correlate this reference with a sequence
 * number, our program searches each entry for the hash of【本・本】.
 * There are two entries in which it is found in JMdict (English):
 *
 * sequence 1260670: 【元・もと】、【元・元】、【もと・もと】、【本・もと】、【本・本】、【素・もと】、【素・素】、【基・もと】、【基・基】
 * sequence 1522150: 【本・ほん】、【本・本】、【ほん・ほん】
 *
 * Because 【本・本】 is closer to the beginning of the array in the
 * latter (i.e., has the lowest index), sequence number 1522150 is
 * returned.
 *
 * In situations in which multiple sequences are found with the same
 * score, the entry with a priority tag ("news1", "ichi1", "spec1",
 * "spec2", "gai1") is given preference. This mostly affects
 * katakana-only loanwords like ラグ.
 *
 * To improve accuracy, this method also checks to see if the
 * reference's specified sense number really exists in the
 * corresponding entry. For example, sequence 1582850 【如何で・いかんで】
 * has a reference to sense #2 of いかん (no kanji specified), which
 * could belong to 13 different sequences. However, sequences 1582850
 * and 2829697 are the only 2 of those 13 which contain more than one
 * sense. Incidentally, sequence 1582850 is the correct match.
 *
 * All else being equal, the entry with the smallest sequence number
 * is chosen. References in the JMdict file are currently ambiguous,
 * and getting this perfect won't be possible until sequence numbers
 * are explictly identified in these references.  See:
 * https://github.com/JMdictProject/JMdictIssues/issues/61
 */
func (meta *jmdictMetadata) FindBestSequence(reference string) sequence {
	bestSeq := 0
	lowestIndex := 100000
	bestIsPriority := false
	headword, senseNumber, ok := parseReference(reference)
	if !ok {
		return bestSeq
	}
	hash := headword.Hash()
	for _, seqScore := range meta.hashToSearchValues[hash] {
		if meta.seqToSenseCount[seqScore.sequence] < senseNumber {
			// entry must contain the specified sense
			continue
		} else if lowestIndex < seqScore.index {
			// lower indices are better
			continue
		} else if (lowestIndex == seqScore.index) && (bestIsPriority && !seqScore.isPriority) {
			// if scores match, check priority
			continue
		} else if (lowestIndex == seqScore.index) && (bestIsPriority == seqScore.isPriority) && (bestSeq < seqScore.sequence) {
			// if scores and priority match, check sequence number.
			// lower sequence numbers are better
			continue
		} else {
			lowestIndex = seqScore.index
			bestSeq = seqScore.sequence
			bestIsPriority = seqScore.isPriority
		}
	}
	return bestSeq
}
