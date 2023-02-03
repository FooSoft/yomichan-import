package yomichan

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func senseNumberTags(maxSenseCount int) []dbTag {
	tags := []dbTag{}
	for i := 1; i <= maxSenseCount; i++ {
		tag := dbTag{
			Name:  strconv.Itoa(i),
			Order: -10, // these tags will appear on the left side
			Notes: "JMdict Sense #" + strconv.Itoa(i),
		}
		tags = append(tags, tag)
	}
	return tags
}

func newsFrequencyTags() []dbTag {
	// 24,000 ranks divided into 24 tags, news1k ... news24k
	tags := []dbTag{}
	for i := 1; i <= 24; i++ {
		tagName := "news" + strconv.Itoa(i) + "k"
		var startRank string
		if i == 1 {
			startRank = "1"
		} else {
			// technically should be ",001", but that looks odd
			startRank = strconv.Itoa(i-1) + ",000"
		}
		endRank := strconv.Itoa(i) + ",000"
		tag := dbTag{
			Name:     tagName,
			Order:    -2,
			Score:    0,
			Category: "frequent",
			Notes:    "ranked between the top " + startRank + " and " + endRank + " words in a frequency analysis of the Mainichi Shimbun (1990s)",
		}
		tags = append(tags, tag)
	}
	return tags
}

func entityTags(entities map[string]string) []dbTag {
	tags := knownEntityTags()
	for name, notes := range entities {
		idx := slices.IndexFunc(tags, func(t dbTag) bool { return t.Name == name })
		if idx != -1 {
			tags[idx].Notes = notes
		} else {
			fmt.Println("Unknown tag type \"" + name + "\": " + notes)
			unknownTag := dbTag{Name: name, Notes: notes}
			tags = append(tags, unknownTag)
		}
	}
	return tags
}

func customDbTags() []dbTag {
	return []dbTag{
		dbTag{Name: priorityTagName, Order: -10, Score: 10, Category: "popular", Notes: "high priority term"},
		dbTag{Name: rareKanjiTagName, Order: 0, Score: -5, Category: "archaism", Notes: "rarely-used kanji form of this expression"},
		dbTag{Name: irregularTagName, Order: 0, Score: -5, Category: "archaism", Notes: "irregular form of this expression"},
		dbTag{Name: outdatedTagName, Order: 0, Score: -5, Category: "archaism", Notes: "outdated form of this expression"},
		dbTag{Name: "ichi", Order: -2, Score: 0, Category: "frequent", Notes: "included in Ichimango Goi Bunruishuu (１万語語彙分類集)"},
		dbTag{Name: "spec", Order: -2, Score: 0, Category: "frequent", Notes: "specified as common by JMdict editors"},
		dbTag{Name: "gai", Order: -2, Score: 0, Category: "frequent", Notes: "common loanword (gairaigo・外来語)"},
		dbTag{Name: "forms", Order: 0, Score: 0, Category: "", Notes: "other surface forms and readings"},
	}
}

func knownEntityTags() []dbTag {
	return []dbTag{
		// see: https://www.edrdg.org/jmdictdb/cgi-bin/edhelp.py?svc=jmdict&sid=#kwabbr
		// additional descriptions at the beginning of the JMdict file

		// <re_inf> reading info
		dbTag{Name: "gikun", Order: 0, Score: 0, Category: ""}, // gikun (meaning as reading) or jukujikun (special kanji reading)
		dbTag{Name: "ik", Order: 0, Score: -5, Category: ""},   // word containing irregular kana usage
		dbTag{Name: "ok", Order: 0, Score: -5, Category: ""},   // out-dated or obsolete kana usage
		dbTag{Name: "sk", Order: 0, Score: -5, Category: ""},   // search-only kana form

		// <ke_inf> kanji info
		/* kanji info also has a "ik" entity that would go here if not already for the re_inf tag */
		dbTag{Name: "ateji", Order: 0, Score: 0, Category: ""}, // ateji (phonetic) reading
		dbTag{Name: "iK", Order: 0, Score: -5, Category: ""},   // word containing irregular kanji usage
		dbTag{Name: "io", Order: 0, Score: -5, Category: ""},   // irregular okurigana usage
		dbTag{Name: "oK", Order: 0, Score: -5, Category: ""},   // word containing out-dated kanji or kanji usage
		dbTag{Name: "rK", Order: 0, Score: -5, Category: ""},   // rarely-used kanji form
		dbTag{Name: "sK", Order: 0, Score: -5, Category: ""},   // search-only kanji form

		// <misc> miscellaneous sense info
		dbTag{Name: "abbr", Order: 0, Score: 0, Category: ""},              // abbreviation
		dbTag{Name: "arch", Order: -4, Score: 0, Category: "archaism"},     // archaism
		dbTag{Name: "char", Order: 4, Score: 0, Category: "name"},          // character
		dbTag{Name: "chn", Order: 0, Score: 0, Category: ""},               // children's language
		dbTag{Name: "col", Order: 0, Score: 0, Category: ""},               // colloquialism
		dbTag{Name: "company", Order: 4, Score: 0, Category: "name"},       // company name
		dbTag{Name: "creat", Order: 4, Score: 0, Category: "name"},         // creature
		dbTag{Name: "dated", Order: -4, Score: 0, Category: "archaism"},    // dated term
		dbTag{Name: "dei", Order: 4, Score: 0, Category: "name"},           // deity
		dbTag{Name: "derog", Order: 0, Score: 0, Category: ""},             // derogatory
		dbTag{Name: "doc", Order: 4, Score: 0, Category: "name"},           // document
		dbTag{Name: "euph", Order: 0, Score: 0, Category: ""},              // euphemistic
		dbTag{Name: "ev", Order: 4, Score: 0, Category: "name"},            // event
		dbTag{Name: "fam", Order: 0, Score: 0, Category: ""},               // familiar language
		dbTag{Name: "fem", Order: 4, Score: 0, Category: "name"},           // female term, language, or name
		dbTag{Name: "fict", Order: 4, Score: 0, Category: "name"},          // fiction
		dbTag{Name: "form", Order: 0, Score: 0, Category: ""},              // formal or literary term
		dbTag{Name: "given", Order: 4, Score: 0, Category: "name"},         // given name or forename, gender not specified
		dbTag{Name: "group", Order: 4, Score: 0, Category: "name"},         // group
		dbTag{Name: "hist", Order: 0, Score: 0, Category: ""},              // historical term
		dbTag{Name: "hon", Order: 0, Score: 0, Category: ""},               // honorific or respectful (sonkeigo) language
		dbTag{Name: "hum", Order: 0, Score: 0, Category: ""},               // humble (kenjougo) language
		dbTag{Name: "id", Order: -5, Score: 0, Category: "expression"},     // idiomatic expression
		dbTag{Name: "joc", Order: 0, Score: 0, Category: ""},               // jocular, humorous term
		dbTag{Name: "leg", Order: 4, Score: 0, Category: "name"},           // legend
		dbTag{Name: "m-sl", Order: 0, Score: 0, Category: ""},              // manga slang
		dbTag{Name: "male", Order: 4, Score: 0, Category: "name"},          // male term, language, or name
		dbTag{Name: "masc", Order: 4, Score: 0, Category: "name"},          // male term, language, or name
		dbTag{Name: "myth", Order: 4, Score: 0, Category: "name"},          // mythology
		dbTag{Name: "net-sl", Order: 0, Score: 0, Category: ""},            // Internet slang
		dbTag{Name: "obj", Order: 4, Score: 0, Category: "name"},           // object
		dbTag{Name: "obs", Order: -4, Score: 0, Category: "archaism"},      // obsolete term
		dbTag{Name: "on-mim", Order: 0, Score: 0, Category: ""},            // onomatopoeic or mimetic word
		dbTag{Name: "organization", Order: 4, Score: 0, Category: "name"},  // organization name
		dbTag{Name: "oth", Order: 4, Score: 0, Category: "name"},           // other
		dbTag{Name: "person", Order: 4, Score: 0, Category: "name"},        // full name of a particular person
		dbTag{Name: "place", Order: 4, Score: 0, Category: "name"},         // place name
		dbTag{Name: "poet", Order: 0, Score: 0, Category: ""},              // poetical term
		dbTag{Name: "pol", Order: 0, Score: 0, Category: ""},               // polite (teineigo) language
		dbTag{Name: "product", Order: 4, Score: 0, Category: "name"},       // product name
		dbTag{Name: "proverb", Order: 0, Score: 0, Category: "expression"}, // proverb
		dbTag{Name: "quote", Order: 0, Score: 0, Category: "expression"},   // quotation
		dbTag{Name: "rare", Order: -4, Score: 0, Category: "archaism"},     // rare
		dbTag{Name: "relig", Order: 4, Score: 0, Category: "name"},         // religion
		dbTag{Name: "sens", Order: 0, Score: 0, Category: ""},              // sensitive
		dbTag{Name: "serv", Order: 4, Score: 0, Category: "name"},          // service
		dbTag{Name: "ship", Order: 4, Score: 0, Category: "name"},          // ship name
		dbTag{Name: "sl", Order: 0, Score: 0, Category: ""},                // slang
		dbTag{Name: "station", Order: 4, Score: 0, Category: "name"},       // railway station
		dbTag{Name: "surname", Order: 4, Score: 0, Category: "name"},       // family or surname
		dbTag{Name: "uk", Order: 0, Score: 0, Category: ""},                // word usually written using kana alone
		dbTag{Name: "unclass", Order: 4, Score: 0, Category: "name"},       // unclassified name
		dbTag{Name: "vulg", Order: 0, Score: 0, Category: ""},              // vulgar expression or word
		dbTag{Name: "work", Order: 4, Score: 0, Category: "name"},          // work of art, literature, music, etc. name
		dbTag{Name: "X", Order: 0, Score: 0, Category: ""},                 // rude or X-rated term (not displayed in educational software)
		dbTag{Name: "yoji", Order: 0, Score: 0, Category: ""},              // yojijukugo

		// <pos> part-of-speech info
		dbTag{Name: "adj-f", Order: -3, Score: 0, Category: "partOfSpeech"},     // noun or verb acting prenominally
		dbTag{Name: "adj-i", Order: -3, Score: 0, Category: "partOfSpeech"},     // adjective (keiyoushi)
		dbTag{Name: "adj-ix", Order: -3, Score: 0, Category: "partOfSpeech"},    // adjective (keiyoushi) - yoi/ii class
		dbTag{Name: "adj-kari", Order: -3, Score: 0, Category: "partOfSpeech"},  // 'kari' adjective (archaic)
		dbTag{Name: "adj-ku", Order: -3, Score: 0, Category: "partOfSpeech"},    // 'ku' adjective (archaic)
		dbTag{Name: "adj-na", Order: -3, Score: 0, Category: "partOfSpeech"},    // adjectival nouns or quasi-adjectives (keiyodoshi)
		dbTag{Name: "adj-nari", Order: -3, Score: 0, Category: "partOfSpeech"},  // archaic/formal form of na-adjective
		dbTag{Name: "adj-no", Order: -3, Score: 0, Category: "partOfSpeech"},    // nouns which may take the genitive case particle 'no'
		dbTag{Name: "adj-pn", Order: -3, Score: 0, Category: "partOfSpeech"},    // pre-noun adjectival (rentaishi)
		dbTag{Name: "adj-shiku", Order: -3, Score: 0, Category: "partOfSpeech"}, // 'shiku' adjective (archaic)
		dbTag{Name: "adj-t", Order: -3, Score: 0, Category: "partOfSpeech"},     // 'taru' adjective
		dbTag{Name: "adv", Order: -3, Score: 0, Category: "partOfSpeech"},       // adverb (fukushi)
		dbTag{Name: "adv-to", Order: -3, Score: 0, Category: "partOfSpeech"},    // adverb taking the 'to' particle
		dbTag{Name: "aux", Order: -3, Score: 0, Category: "partOfSpeech"},       // auxiliary
		dbTag{Name: "aux-adj", Order: -3, Score: 0, Category: "partOfSpeech"},   // auxiliary adjective
		dbTag{Name: "aux-v", Order: -3, Score: 0, Category: "partOfSpeech"},     // auxiliary verb
		dbTag{Name: "conj", Order: -3, Score: 0, Category: "partOfSpeech"},      // conjunction
		dbTag{Name: "cop", Order: -3, Score: 0, Category: "partOfSpeech"},       // copula
		dbTag{Name: "ctr", Order: -3, Score: 0, Category: "partOfSpeech"},       // counter
		dbTag{Name: "exp", Order: -5, Score: 0, Category: "expression"},         // expressions (phrases, clauses, etc.)
		dbTag{Name: "int", Order: -3, Score: 0, Category: "partOfSpeech"},       // interjection (kandoushi)
		dbTag{Name: "n", Order: -3, Score: 0, Category: "partOfSpeech"},         // noun (common) (futsuumeishi)
		dbTag{Name: "n-adv", Order: -3, Score: 0, Category: "partOfSpeech"},     // adverbial noun (fukushitekimeishi)
		dbTag{Name: "n-pr", Order: -3, Score: 0, Category: "partOfSpeech"},      // proper noun
		dbTag{Name: "n-pref", Order: -3, Score: 0, Category: "partOfSpeech"},    // noun, used as a prefix
		dbTag{Name: "n-suf", Order: -3, Score: 0, Category: "partOfSpeech"},     // noun, used as a suffix
		dbTag{Name: "n-t", Order: -3, Score: 0, Category: "partOfSpeech"},       // noun (temporal) (jisoumeishi)
		dbTag{Name: "num", Order: -3, Score: 0, Category: "partOfSpeech"},       // numeric
		dbTag{Name: "pn", Order: -3, Score: 0, Category: "partOfSpeech"},        // pronoun
		dbTag{Name: "pref", Order: -3, Score: 0, Category: "partOfSpeech"},      // prefix
		dbTag{Name: "prt", Order: -3, Score: 0, Category: "partOfSpeech"},       // particle
		dbTag{Name: "suf", Order: -3, Score: 0, Category: "partOfSpeech"},       // suffix
		dbTag{Name: "unc", Order: -3, Score: 0, Category: "partOfSpeech"},       // unclassified
		dbTag{Name: "v-unspec", Order: -3, Score: 0, Category: "partOfSpeech"},  // verb unspecified
		dbTag{Name: "v1", Order: -3, Score: 0, Category: "partOfSpeech"},        // Ichidan verb
		dbTag{Name: "v1-s", Order: -3, Score: 0, Category: "partOfSpeech"},      // Ichidan verb - kureru special class
		dbTag{Name: "v2a-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb with 'u' ending (archaic)
		dbTag{Name: "v2b-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'bu' ending (archaic)
		dbTag{Name: "v2b-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'bu' ending (archaic)
		dbTag{Name: "v2d-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'dzu' ending (archaic)
		dbTag{Name: "v2d-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'dzu' ending (archaic)
		dbTag{Name: "v2g-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'gu' ending (archaic)
		dbTag{Name: "v2g-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'gu' ending (archaic)
		dbTag{Name: "v2h-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'hu/fu' ending (archaic)
		dbTag{Name: "v2h-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'hu/fu' ending (archaic)
		dbTag{Name: "v2k-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'ku' ending (archaic)
		dbTag{Name: "v2k-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'ku' ending (archaic)
		dbTag{Name: "v2m-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'mu' ending (archaic)
		dbTag{Name: "v2m-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'mu' ending (archaic)
		dbTag{Name: "v2n-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'nu' ending (archaic)
		dbTag{Name: "v2r-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'ru' ending (archaic)
		dbTag{Name: "v2r-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'ru' ending (archaic)
		dbTag{Name: "v2s-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'su' ending (archaic)
		dbTag{Name: "v2t-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'tsu' ending (archaic)
		dbTag{Name: "v2t-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'tsu' ending (archaic)
		dbTag{Name: "v2w-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'u' ending and 'we' conjugation (archaic)
		dbTag{Name: "v2y-k", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (upper class) with 'yu' ending (archaic)
		dbTag{Name: "v2y-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'yu' ending (archaic)
		dbTag{Name: "v2z-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Nidan verb (lower class) with 'zu' ending (archaic)
		dbTag{Name: "v4b", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'bu' ending (archaic)
		dbTag{Name: "v4g", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'gu' ending (archaic)
		dbTag{Name: "v4h", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'hu/fu' ending (archaic)
		dbTag{Name: "v4k", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'ku' ending (archaic)
		dbTag{Name: "v4m", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'mu' ending (archaic)
		dbTag{Name: "v4n", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'nu' ending (archaic)
		dbTag{Name: "v4r", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'ru' ending (archaic)
		dbTag{Name: "v4s", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'su' ending (archaic)
		dbTag{Name: "v4t", Order: -3, Score: 0, Category: "partOfSpeech"},       // Yodan verb with 'tsu' ending (archaic)
		dbTag{Name: "v5aru", Order: -3, Score: 0, Category: "partOfSpeech"},     // Godan verb - -aru special class
		dbTag{Name: "v5b", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'bu' ending
		dbTag{Name: "v5g", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'gu' ending
		dbTag{Name: "v5k", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'ku' ending
		dbTag{Name: "v5k-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Godan verb - Iku/Yuku special class
		dbTag{Name: "v5m", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'mu' ending
		dbTag{Name: "v5n", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'nu' ending
		dbTag{Name: "v5r", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'ru' ending
		dbTag{Name: "v5r-i", Order: -3, Score: 0, Category: "partOfSpeech"},     // Godan verb with 'ru' ending (irregular verb)
		dbTag{Name: "v5s", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'su' ending
		dbTag{Name: "v5t", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'tsu' ending
		dbTag{Name: "v5u", Order: -3, Score: 0, Category: "partOfSpeech"},       // Godan verb with 'u' ending
		dbTag{Name: "v5u-s", Order: -3, Score: 0, Category: "partOfSpeech"},     // Godan verb with 'u' ending (special class)
		dbTag{Name: "v5uru", Order: -3, Score: 0, Category: "partOfSpeech"},     // Godan verb - Uru old class verb (old form of Eru)
		dbTag{Name: "vi", Order: -3, Score: 0, Category: "partOfSpeech"},        // intransitive verb
		dbTag{Name: "vk", Order: -3, Score: 0, Category: "partOfSpeech"},        // Kuru verb - special class
		dbTag{Name: "vn", Order: -3, Score: 0, Category: "partOfSpeech"},        // irregular nu verb
		dbTag{Name: "vr", Order: -3, Score: 0, Category: "partOfSpeech"},        // irregular ru verb, plain form ends with -ri
		dbTag{Name: "vs", Order: -3, Score: 0, Category: "partOfSpeech"},        // noun or participle which takes the aux. verb suru
		dbTag{Name: "vs-c", Order: -3, Score: 0, Category: "partOfSpeech"},      // su verb - precursor to the modern suru
		dbTag{Name: "vs-i", Order: -3, Score: 0, Category: "partOfSpeech"},      // suru verb - included
		dbTag{Name: "vs-s", Order: -3, Score: 0, Category: "partOfSpeech"},      // suru verb - special class
		dbTag{Name: "vt", Order: -3, Score: 0, Category: "partOfSpeech"},        // transitive verb
		dbTag{Name: "vz", Order: -3, Score: 0, Category: "partOfSpeech"},        // Ichidan verb - zuru verb (alternative form of -jiru verbs)

		// <field> usage domain
		dbTag{Name: "agric", Order: 0, Score: 0, Category: ""},    // agriculture
		dbTag{Name: "anat", Order: 0, Score: 0, Category: ""},     // anatomy
		dbTag{Name: "archeol", Order: 0, Score: 0, Category: ""},  // archeology
		dbTag{Name: "archit", Order: 0, Score: 0, Category: ""},   // architecture
		dbTag{Name: "art", Order: 0, Score: 0, Category: ""},      // art, aesthetics
		dbTag{Name: "astron", Order: 0, Score: 0, Category: ""},   // astronomy
		dbTag{Name: "audvid", Order: 0, Score: 0, Category: ""},   // audiovisual
		dbTag{Name: "aviat", Order: 0, Score: 0, Category: ""},    // aviation
		dbTag{Name: "baseb", Order: 0, Score: 0, Category: ""},    // baseball
		dbTag{Name: "biochem", Order: 0, Score: 0, Category: ""},  // biochemistry
		dbTag{Name: "biol", Order: 0, Score: 0, Category: ""},     // biology
		dbTag{Name: "bot", Order: 0, Score: 0, Category: ""},      // botany
		dbTag{Name: "Buddh", Order: 0, Score: 0, Category: ""},    // Buddhism
		dbTag{Name: "bus", Order: 0, Score: 0, Category: ""},      // business
		dbTag{Name: "cards", Order: 0, Score: 0, Category: ""},    // card games
		dbTag{Name: "chem", Order: 0, Score: 0, Category: ""},     // chemistry
		dbTag{Name: "Christn", Order: 0, Score: 0, Category: ""},  // Christianity
		dbTag{Name: "cloth", Order: 0, Score: 0, Category: ""},    // clothing
		dbTag{Name: "comp", Order: 0, Score: 0, Category: ""},     // computing
		dbTag{Name: "cryst", Order: 0, Score: 0, Category: ""},    // crystallography
		dbTag{Name: "dent", Order: 0, Score: 0, Category: ""},     // dentistry
		dbTag{Name: "ecol", Order: 0, Score: 0, Category: ""},     // ecology
		dbTag{Name: "econ", Order: 0, Score: 0, Category: ""},     // economics
		dbTag{Name: "elec", Order: 0, Score: 0, Category: ""},     // electricity, elec. eng.
		dbTag{Name: "electr", Order: 0, Score: 0, Category: ""},   // electronics
		dbTag{Name: "embryo", Order: 0, Score: 0, Category: ""},   // embryology
		dbTag{Name: "engr", Order: 0, Score: 0, Category: ""},     // engineering
		dbTag{Name: "ent", Order: 0, Score: 0, Category: ""},      // entomology
		dbTag{Name: "film", Order: 0, Score: 0, Category: ""},     // film
		dbTag{Name: "finc", Order: 0, Score: 0, Category: ""},     // finance
		dbTag{Name: "fish", Order: 0, Score: 0, Category: ""},     // fishing
		dbTag{Name: "food", Order: 0, Score: 0, Category: ""},     // food, cooking
		dbTag{Name: "gardn", Order: 0, Score: 0, Category: ""},    // gardening, horticulture
		dbTag{Name: "genet", Order: 0, Score: 0, Category: ""},    // genetics
		dbTag{Name: "geogr", Order: 0, Score: 0, Category: ""},    // geography
		dbTag{Name: "geol", Order: 0, Score: 0, Category: ""},     // geology
		dbTag{Name: "geom", Order: 0, Score: 0, Category: ""},     // geometry
		dbTag{Name: "go", Order: 0, Score: 0, Category: ""},       // go (game)
		dbTag{Name: "golf", Order: 0, Score: 0, Category: ""},     // golf
		dbTag{Name: "gramm", Order: 0, Score: 0, Category: ""},    // grammar
		dbTag{Name: "grmyth", Order: 0, Score: 0, Category: ""},   // Greek mythology
		dbTag{Name: "hanaf", Order: 0, Score: 0, Category: ""},    // hanafuda
		dbTag{Name: "horse", Order: 0, Score: 0, Category: ""},    // horse racing
		dbTag{Name: "kabuki", Order: 0, Score: 0, Category: ""},   // kabuki
		dbTag{Name: "law", Order: 0, Score: 0, Category: ""},      // law
		dbTag{Name: "ling", Order: 0, Score: 0, Category: ""},     // linguistics
		dbTag{Name: "logic", Order: 0, Score: 0, Category: ""},    // logic
		dbTag{Name: "MA", Order: 0, Score: 0, Category: ""},       // martial arts
		dbTag{Name: "mahj", Order: 0, Score: 0, Category: ""},     // mahjong
		dbTag{Name: "manga", Order: 0, Score: 0, Category: ""},    // manga
		dbTag{Name: "math", Order: 0, Score: 0, Category: ""},     // mathematics
		dbTag{Name: "mech", Order: 0, Score: 0, Category: ""},     // mechanical engineering
		dbTag{Name: "med", Order: 0, Score: 0, Category: ""},      // medicine
		dbTag{Name: "met", Order: 0, Score: 0, Category: ""},      // meteorology
		dbTag{Name: "mil", Order: 0, Score: 0, Category: ""},      // military
		dbTag{Name: "mining", Order: 0, Score: 0, Category: ""},   // mining
		dbTag{Name: "music", Order: 0, Score: 0, Category: ""},    // music
		dbTag{Name: "noh", Order: 0, Score: 0, Category: ""},      // noh
		dbTag{Name: "ornith", Order: 0, Score: 0, Category: ""},   // ornithology
		dbTag{Name: "paleo", Order: 0, Score: 0, Category: ""},    // paleontology
		dbTag{Name: "pathol", Order: 0, Score: 0, Category: ""},   // pathology
		dbTag{Name: "pharm", Order: 0, Score: 0, Category: ""},    // pharmacy
		dbTag{Name: "phil", Order: 0, Score: 0, Category: ""},     // philosophy
		dbTag{Name: "photo", Order: 0, Score: 0, Category: ""},    // photography
		dbTag{Name: "physics", Order: 0, Score: 0, Category: ""},  // physics
		dbTag{Name: "physiol", Order: 0, Score: 0, Category: ""},  // physiology
		dbTag{Name: "politics", Order: 0, Score: 0, Category: ""}, // politics
		dbTag{Name: "print", Order: 0, Score: 0, Category: ""},    // printing
		dbTag{Name: "psy", Order: 0, Score: 0, Category: ""},      // psychiatry
		dbTag{Name: "psyanal", Order: 0, Score: 0, Category: ""},  // psychoanalysis
		dbTag{Name: "psych", Order: 0, Score: 0, Category: ""},    // psychology
		dbTag{Name: "rail", Order: 0, Score: 0, Category: ""},     // railway
		dbTag{Name: "rommyth", Order: 0, Score: 0, Category: ""},  // Roman mythology
		dbTag{Name: "Shinto", Order: 0, Score: 0, Category: ""},   // Shinto
		dbTag{Name: "shogi", Order: 0, Score: 0, Category: ""},    // shogi
		dbTag{Name: "ski", Order: 0, Score: 0, Category: ""},      // skiing
		dbTag{Name: "sports", Order: 0, Score: 0, Category: ""},   // sports
		dbTag{Name: "stat", Order: 0, Score: 0, Category: ""},     // statistics
		dbTag{Name: "stockm", Order: 0, Score: 0, Category: ""},   // stock market
		dbTag{Name: "sumo", Order: 0, Score: 0, Category: ""},     // sumo
		dbTag{Name: "telec", Order: 0, Score: 0, Category: ""},    // telecommunications
		dbTag{Name: "tradem", Order: 0, Score: 0, Category: ""},   // trademark
		dbTag{Name: "tv", Order: 0, Score: 0, Category: ""},       // television
		dbTag{Name: "vidg", Order: 0, Score: 0, Category: ""},     // video games
		dbTag{Name: "zool", Order: 0, Score: 0, Category: ""},     // zoology

		// <dial> dialect
		dbTag{Name: "bra", Order: 0, Score: 0, Category: ""},  // Brazilian
		dbTag{Name: "hob", Order: 0, Score: 0, Category: ""},  // Hokkaido-ben
		dbTag{Name: "ksb", Order: 0, Score: 0, Category: ""},  // Kansai-ben
		dbTag{Name: "ktb", Order: 0, Score: 0, Category: ""},  // Kantou-ben
		dbTag{Name: "kyb", Order: 0, Score: 0, Category: ""},  // Kyoto-ben
		dbTag{Name: "kyu", Order: 0, Score: 0, Category: ""},  // Kyuushuu-ben
		dbTag{Name: "nab", Order: 0, Score: 0, Category: ""},  // Nagano-ben
		dbTag{Name: "osb", Order: 0, Score: 0, Category: ""},  // Osaka-ben
		dbTag{Name: "rkb", Order: 0, Score: 0, Category: ""},  // Ryuukyuu-ben
		dbTag{Name: "thb", Order: 0, Score: 0, Category: ""},  // Touhoku-ben
		dbTag{Name: "tsb", Order: 0, Score: 0, Category: ""},  // Tosa-ben
		dbTag{Name: "tsug", Order: 0, Score: 0, Category: ""}, // Tsugaru-ben
	}
}
