package yomichan

import (
	"strings"

	"golang.org/x/exp/slices"
)

// Returns text with all katakana characters converted into hiragana.
func katakanaToHiragana(text string) string {
	f := func(x rune) rune {
		if x >= 'ァ' && x <= 'ヶ' || x >= 'ヽ' && x <= 'ヾ' {
			return x - 0x60
		} else {
			return x
		}
	}
	return strings.Map(f, text)
}

// Replace hiragana iteration marks with the appropriate characters.
// E.g. "さゝき" -> "ささき"; "たゞの" -> "ただの"
func replaceIterationMarks(text string) string {
	iterationMarks := []struct {
		char   rune
		offset rune
	}{
		{'ゝ', 0x00},
		{'ゞ', 0x01},
	}
	for _, x := range iterationMarks {
		for strings.IndexRune(text, x.char) > 0 {
			runes := []rune(text)
			idx := slices.Index(runes, x.char)
			runes[idx] = runes[idx-1] + x.offset
			text = string(runes)
		}
	}
	return text
}

// Returns an array of the input text split into segments.
// E.g. "しょくぎょう" -> ["しょ", "く", "ぎょ", "う"]
// Returns nil if no segmentation is possible.
func makeKanaSegments(kana string) (segments []string) {
	hiragana := replaceIterationMarks(katakanaToHiragana(kana))
	kanaRunes := []rune{}
	for _, kanaRune := range hiragana {
		kanaRunes = append(kanaRunes, kanaRune)
	}
	kanaRuneCount := len(kanaRunes)
	for i := 0; i < kanaRuneCount; i++ {
		for j := 0; j < kanaRuneCount-i; j++ {
			segment := string(kanaRunes[i : kanaRuneCount-j])
			if _, ok := kanaSegmentToRomajiList[segment]; ok {
				segments = append(segments, segment)
				i = kanaRuneCount - j - 1
				break
			}
			if j == kanaRuneCount-i-1 {
				return nil
			}
		}
	}
	return segments
}

// Returns a map of ltr substrings of the input text.
// E.g. "nihon" -> ["n", "ni", "nih", "niho", "nihon"]
func makeSubstringMap(text string) map[string]bool {
	substrings := make(map[string]bool)
	for i := 1; i <= len(text); i++ {
		substring := text[:i]
		substrings[substring] = true
	}
	return substrings
}

// Determines if the input text is a valid romaji representation of
// the input kana.
//
// The strategy is to calculate every possible romaji representation
// of a given string of kana and check if the input text is one of
// them. Since the number of combinations grows very large for long
// strings of kana, we need to prune invalid branches from the
// combination tree along the way.
func isTransliteration(text string, kana string) bool {
	romaji := strings.TrimSpace(strings.ToLower(text))
	validSubstrings := makeSubstringMap(romaji)
	kanaSegments := makeKanaSegments(kana)
	possibilities := []string{""}
	for _, segment := range kanaSegments {
		newPossibilities := map[string]bool{}
		for _, x := range possibilities {
			for _, y := range kanaSegmentToRomajiList[segment] {
				z := x + y
				newPossibilities[z] = true
			}
		}
		possibilities = nil
		for z := range newPossibilities {
			if validSubstrings[z] {
				possibilities = append(possibilities, z)
			}
		}
		if possibilities == nil {
			return false
		}
	}
	return slices.Contains(possibilities, romaji)
}

var kanaSegmentToRomajiList = map[string][]string{
	"ぁ":  []string{"", "a"},
	"ぃ":  []string{"", "i"},
	"ぅ":  []string{"", "u"},
	"ぇ":  []string{"", "e"},
	"ぉ":  []string{"", "o"},
	"ゃ":  []string{"ya"},
	"ゅ":  []string{"yu"},
	"ょ":  []string{"yo"},
	"ゎ":  []string{"wa"},
	"っ":  []string{"", "k", "g", "s", "z", "t", "d", "f", "h", "b", "p", "n", "m", "y", "w", "c"},
	"ー":  []string{"", "a", "i", "u", "e", "o", "-"},
	"あ":  []string{"", "a", "ā", "wa", "wā"},
	"い":  []string{"", "i", "ī", "wi", "wī"},
	"う":  []string{"", "u", "ū", "wu", "wū"},
	"え":  []string{"", "e", "ē", "we", "wē"},
	"お":  []string{"", "o", "ō", "wo", "wō"},
	"ゔ":  []string{"vu", "vū", "bu", "bū"},
	"か":  []string{"ka", "kā"},
	"が":  []string{"ga", "gā"},
	"き":  []string{"ki", "kī"},
	"ぎ":  []string{"gi", "gī"},
	"く":  []string{"ku", "kū"},
	"ぐ":  []string{"gu", "gū"},
	"け":  []string{"ke", "kē"},
	"げ":  []string{"ge", "gē"},
	"こ":  []string{"ko", "kō"},
	"ご":  []string{"go", "gō"},
	"さ":  []string{"sa", "sā"},
	"ざ":  []string{"za", "zā"},
	"し":  []string{"si", "sī", "shi", "shī"},
	"じ":  []string{"zi", "zī", "ji", "jī"},
	"す":  []string{"su", "sū"},
	"ず":  []string{"zu", "zū"},
	"せ":  []string{"se", "sē"},
	"ぜ":  []string{"ze", "zē"},
	"そ":  []string{"so", "sō"},
	"ぞ":  []string{"zo", "zō"},
	"た":  []string{"ta", "tā"},
	"だ":  []string{"da", "dā"},
	"ち":  []string{"ti", "tī", "chi", "chī"},
	"ぢ":  []string{"di", "dī", "dhi", "dhī", "ji", "jī", "dji", "djī", "dzi", "dzī"},
	"つ":  []string{"tu", "tū", "tsu", "tsū"},
	"づ":  []string{"du", "dū", "dzu", "dzū", "zu", "zū"},
	"て":  []string{"te", "tē"},
	"で":  []string{"de", "dē"},
	"と":  []string{"to", "tō"},
	"ど":  []string{"do", "dō"},
	"な":  []string{"na", "nā"},
	"に":  []string{"ni", "nī"},
	"ぬ":  []string{"nu", "nū"},
	"ね":  []string{"ne", "nē"},
	"の":  []string{"no", "nō"},
	"は":  []string{"ha", "hā", "wa", "wā", "a", "ā"},
	"ば":  []string{"ba", "bā"},
	"ぱ":  []string{"pa", "pā"},
	"ひ":  []string{"hi", "hī", "i", "ī"},
	"び":  []string{"bi", "bī"},
	"ぴ":  []string{"pi", "pī"},
	"ふ":  []string{"hu", "hū", "fu", "fū", "u", "ū"},
	"ぶ":  []string{"bu", "bū"},
	"ぷ":  []string{"pu", "pū"},
	"へ":  []string{"he", "hē", "e", "ē"},
	"べ":  []string{"be", "bē"},
	"ぺ":  []string{"pe", "pē"},
	"ほ":  []string{"ho", "hō", "o", "ō"},
	"ぼ":  []string{"bo", "bō"},
	"ぽ":  []string{"po", "pō"},
	"ま":  []string{"ma", "mā"},
	"み":  []string{"mi", "mī"},
	"む":  []string{"mu", "mū"},
	"め":  []string{"me", "mē"},
	"も":  []string{"mo", "mō"},
	"や":  []string{"ya", "yā"},
	"ゆ":  []string{"yu", "yū"},
	"よ":  []string{"yo", "yō"},
	"ら":  []string{"ra", "rā"},
	"り":  []string{"ri", "rī"},
	"る":  []string{"ru", "rū"},
	"れ":  []string{"re", "rē"},
	"ろ":  []string{"ro", "rō"},
	"わ":  []string{"wa", "wā"},
	"ゐ":  []string{"wi", "wī", "i", "ī"},
	"ゑ":  []string{"we", "wē", "e", "ē"},
	"を":  []string{"wo", "wō", "o", "ō"},
	"ん":  []string{"n", "n'", "m"},
	"うぁ": []string{"wa", "wā", "ua", "uā"},
	"うぃ": []string{"wi", "wī", "ui", "uī"},
	"うぇ": []string{"we", "wē", "ue", "uē"},
	"うぉ": []string{"wo", "wō", "uo", "uō"},
	"きゃ": []string{"kya", "kyā"},
	"きゅ": []string{"kyu", "kyū"},
	"きょ": []string{"kyo", "kyō"},
	"ぎゃ": []string{"gya", "gyā"},
	"ぎゅ": []string{"gyu", "gyū"},
	"ぎょ": []string{"gyo", "gyō"},
	"くゎ": []string{"kwa", "kwā"},
	"くゅ": []string{"kyu", "kyū"},
	"しぇ": []string{"she", "shē", "shie", "shiē"},
	"しゃ": []string{"sha", "shā", "sya", "syā"},
	"しゅ": []string{"shu", "shū", "syu", "syū"},
	"しょ": []string{"sho", "shō", "syo", "syō"},
	"じぇ": []string{"je", "jē"},
	"じゃ": []string{"ja", "jā", "jya", "jyā"},
	"じゅ": []string{"ju", "jū", "jyu", "jyū"},
	"じょ": []string{"jo", "jō", "jyo", "jyō"},
	"ちぁ": []string{"cha", "chā", "chia", "chiā"},
	"ちぇ": []string{"che", "chē", "chie", "chiē"},
	"ちゃ": []string{"cha", "chā", "tya", "tyā"},
	"ちゅ": []string{"chu", "chū", "tyu", "tyū"},
	"ちょ": []string{"cho", "chō", "tyo", "tyō"},
	"ぢゃ": []string{"ja", "jā", "jya", "jyā", "dya", "dyā"},
	"ぢゅ": []string{"ju", "jū", "jyu", "jyū", "dyu", "dyū"},
	"ぢょ": []string{"jo", "jō", "jyo", "jyō", "dyo", "dyō"},
	"つぁ": []string{"tsa", "tsā", "tsua", "tsuā"},
	"つぇ": []string{"tse", "tsē", "tsue", "tsuē"},
	"てぃ": []string{"ti", "tī", "tei", "teī"},
	"でぃ": []string{"di", "dī", "dei", "deī"},
	"でゅ": []string{"dyu", "dyū", "deyu", "deyū"},
	"にゃ": []string{"nya", "nyā"},
	"にゅ": []string{"nyu", "nyū"},
	"にょ": []string{"nyo", "nyō"},
	"ひゃ": []string{"hya", "hyā"},
	"ひゅ": []string{"hyu", "hyū"},
	"ひょ": []string{"hyo", "hyō"},
	"びゃ": []string{"bya", "byā"},
	"びゅ": []string{"byu", "byū"},
	"びょ": []string{"byo", "byō"},
	"ぴゃ": []string{"pya", "pyā"},
	"ぴゅ": []string{"pyu", "pyū"},
	"ぴょ": []string{"pyo", "pyō"},
	"ふぁ": []string{"fa", "fā"},
	"ふぃ": []string{"fi", "fī"},
	"ふぇ": []string{"fe", "fē"},
	"ふぉ": []string{"fo", "fō"},
	"みゃ": []string{"mya", "myā"},
	"みゅ": []string{"myu", "myū"},
	"みょ": []string{"myo", "myō"},
	"りゃ": []string{"rya", "ryā"},
	"りゅ": []string{"ryu", "ryū"},
	"りょ": []string{"ryo", "ryō"},
}
