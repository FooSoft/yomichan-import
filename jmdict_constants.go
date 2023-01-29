package yomichan

type LangCode struct {
	language string
	code     string
}

const edrdgAttribution = "This publication has included material from the JMdict (EDICT, etc.) dictionary files in accordance with the licence provisions of the Electronic Dictionaries Research Group. See http://www.edrdg.org/"

const prioritySymbol = "★"
const rareKanjiSymbol = "🅁"
const irregularSymbol = "⚠"
const outdatedSymbol = "⛬"
const defaultSymbol = "㊒"

const priorityTagName = "⭐"
const rareKanjiTagName = "R"
const irregularTagName = "⚠️"
const outdatedTagName = "⛬"
const atejiTagName = "ateji"
const gikunTagName = "gikun"

const langMarker = "'🌐 '"
const noteMarker = "'📝 '"
const infoMarker = "'ℹ️ '"
const refMarker = "'➡️ '"
const antonymMarker = "'🔄 '"

var ISOtoFlag = map[string]string{
	"":    "'🇬🇧 '",
	"eng": "'🇬🇧 '",
	"dut": "'🇳🇱 '",
	"fre": "'🇫🇷 '",
	"ger": "'🇩🇪 '",
	"hun": "'🇭🇺 '",
	"ita": "'🇮🇹 '",
	"jpn": "'🇯🇵 '",
	"rus": "'🇷🇺 '",
	"slv": "'🇸🇮 '",
	"spa": "'🇪🇸 '",
	"swe": "'🇸🇪 '",
}

var langNameToCode = map[string]string{
	"":          "eng",
	"english":   "eng",
	"dutch":     "dut",
	"french":    "fre",
	"german":    "ger",
	"hungarian": "hun",
	"italian":   "ita",
	"russian":   "rus",
	"slovenian": "slv",
	"spanish":   "spa",
	"swedish":   "swe",
}

var glossTypeCodeToName = map[LangCode]string{
	LangCode{"eng", "lit"}:  "literally",
	LangCode{"eng", "fig"}:  "figuratively",
	LangCode{"eng", "expl"}: "", // don't need to tell the user that an explanation is an explanation
	LangCode{"eng", "tm"}:   "trademark",
}

var refNoteHint = map[LangCode]string{
	LangCode{"eng", "xref"}: "see",
	LangCode{"eng", "ant"}:  "antonym",
}

var sourceLangTypeCodeToType = map[LangCode]string{
	LangCode{"eng", "part"}: "partial",
	LangCode{"eng", ""}:     "", // implied "full"
}

var langCodeToName = map[LangCode]string{
	LangCode{"eng", "afr"}: "Afrikaans",
	LangCode{"eng", "ain"}: "Ainu",
	LangCode{"eng", "alg"}: "Algonquian",
	LangCode{"eng", "amh"}: "Amharic",
	LangCode{"eng", "ara"}: "Arabic",
	LangCode{"eng", "arn"}: "Mapudungun",
	LangCode{"eng", "bnt"}: "Bantu",
	LangCode{"eng", "bre"}: "Breton",
	LangCode{"eng", "bul"}: "Bulgarian",
	LangCode{"eng", "bur"}: "Burmese",
	LangCode{"eng", "chi"}: "Chinese",
	LangCode{"eng", "chn"}: "Chinook Jargon",
	LangCode{"eng", "cze"}: "Czech",
	LangCode{"eng", "dan"}: "Danish",
	LangCode{"eng", "dut"}: "Dutch",
	LangCode{"eng", "eng"}: "English",
	LangCode{"eng", "epo"}: "Esperanto",
	LangCode{"eng", "est"}: "Estonian",
	LangCode{"eng", "fil"}: "Filipino",
	LangCode{"eng", "fin"}: "Finnish",
	LangCode{"eng", "fre"}: "French",
	LangCode{"eng", "geo"}: "Georgian",
	LangCode{"eng", "ger"}: "German",
	LangCode{"eng", "glg"}: "Galician",
	LangCode{"eng", "grc"}: "Ancient Greek",
	LangCode{"eng", "gre"}: "Modern Greek",
	LangCode{"eng", "haw"}: "Hawaiian",
	LangCode{"eng", "heb"}: "Hebrew",
	LangCode{"eng", "hin"}: "Hindi",
	LangCode{"eng", "hun"}: "Hungarian",
	LangCode{"eng", "ice"}: "Icelandic",
	LangCode{"eng", "ind"}: "Indonesian",
	LangCode{"eng", "ita"}: "Italian",
	LangCode{"eng", "khm"}: "Khmer",
	LangCode{"eng", "kor"}: "Korean",
	LangCode{"eng", "kur"}: "Kurdish",
	LangCode{"eng", "lat"}: "Latin",
	LangCode{"eng", "mal"}: "Malayalam",
	LangCode{"eng", "mao"}: "Maori",
	LangCode{"eng", "may"}: "Malay",
	LangCode{"eng", "mnc"}: "Manchu",
	LangCode{"eng", "mol"}: "Moldavian", // ISO 639 deprecated (https://iso639-3.sil.org/code/mol)
	LangCode{"eng", "mon"}: "Mongolian",
	LangCode{"eng", "nor"}: "Norwegian",
	LangCode{"eng", "per"}: "Persian",
	LangCode{"eng", "pol"}: "Polish",
	LangCode{"eng", "por"}: "Portuguese",
	LangCode{"eng", "rum"}: "Romanian",
	LangCode{"eng", "rus"}: "Russian",
	LangCode{"eng", "san"}: "Sanskrit",
	LangCode{"eng", "scr"}: "Croatian", // Code doesn't seem to exist in ISO 639. Should be "hrv" instead? (https://iso639-3.sil.org/code/hrv)
	LangCode{"eng", "slo"}: "Slovak",
	LangCode{"eng", "slv"}: "Slovenian",
	LangCode{"eng", "som"}: "Somali",
	LangCode{"eng", "spa"}: "Spanish",
	LangCode{"eng", "swa"}: "Swahili",
	LangCode{"eng", "swe"}: "Swedish",
	LangCode{"eng", "tah"}: "Tahitian",
	LangCode{"eng", "tam"}: "Tamil",
	LangCode{"eng", "tgl"}: "Tagalog",
	LangCode{"eng", "tha"}: "Thai",
	LangCode{"eng", "tib"}: "Tibetan",
	LangCode{"eng", "tur"}: "Turkish",
	LangCode{"eng", "ukr"}: "Ukrainian",
	LangCode{"eng", "urd"}: "Urdu",
	LangCode{"eng", "vie"}: "Vietnamese",
	LangCode{"eng", "yid"}: "Yiddish",
}

// https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry
var ISOtoHTML = map[string]string{
	"afr": "af",  // Afrikaans
	"ain": "ain", // Ainu
	"alg": "alg", // Algonquian
	"amh": "am",  // Amharic
	"ara": "ar",  // Arabic
	"arn": "arn", // Mapudungun
	"bnt": "bnt", // Bantu
	"bre": "br",  // Breton
	"bul": "bg",  // Bulgarian
	"bur": "my",  // Burmese
	"chi": "zh",  // Chinese
	"chn": "chn", // Chinook Jargon
	"cze": "cs",  // Czech
	"dan": "da",  // Danish
	"dut": "nl",  // Dutch
	"eng": "en",  // English
	"epo": "eo",  // Esperanto
	"est": "et",  // Estonian
	"fil": "fil", // Filipino
	"fin": "fi",  // Finnish
	"fre": "fr",  // French
	"geo": "ka",  // Georgian
	"ger": "de",  // German
	"glg": "gl",  // Galician
	"grc": "grc", // Ancient Greek
	"gre": "el",  // Modern Greek
	"haw": "haw", // Hawaiian
	"heb": "he",  // Hebrew
	"hin": "hi",  // Hindi
	"hun": "hu",  // Hungarian
	"ice": "is",  // Icelandic
	"ind": "id",  // Indonesian
	"ita": "it",  // Italian
	"jpn": "ja",  // Japanese
	"khm": "km",  // Khmer
	"kor": "ko",  // Korean
	"kur": "ku",  // Kurdish
	"lat": "la",  // Latin
	"mal": "ml",  // Malayalam
	"mao": "mi",  // Maori
	"may": "ms",  // Malay
	"mnc": "mnc", // Manchu
	"mol": "ro",  // Moldavian
	"mon": "mn",  // Mongolian
	"nor": "no",  // Norwegian
	"per": "fa",  // Persian
	"pol": "pl",  // Polish
	"por": "pt",  // Portuguese
	"rum": "ro",  // Romanian
	"rus": "ru",  // Russian
	"san": "sa",  // Sanskrit
	"scr": "hr",  // Croatian
	"slo": "sk",  // Slovak
	"slv": "sl",  // Slovenian
	"som": "so",  // Somali
	"spa": "es",  // Spanish
	"swa": "sw",  // Swahili
	"swe": "sv",  // Swedish
	"tah": "ty",  // Tahitian
	"tam": "ta",  // Tamil
	"tgl": "tl",  // Tagalog
	"tha": "th",  // Thai
	"tib": "bo",  // Tibetan
	"tur": "tr",  // Turkish
	"ukr": "uk",  // Ukrainian
	"urd": "ur",  // Urdu
	"vie": "vi",  // Vietnamese
	"yid": "yi",  // Yiddish
}
