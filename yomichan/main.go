package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	yomichan "foosoft.net/projects/yomichan-import"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] input-path output-path\n", path.Base(os.Args[0]))
	fmt.Fprint(os.Stderr, "https://foosoft.net/projects/yomichan-import/\n\n")
	fmt.Fprint(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	var (
		format   = flag.String("format", yomichan.DefaultFormat, "dictionary format [edict|enamdict|epwing|kanjidic|rikai]")
		language = flag.String("language", yomichan.DefaultLanguage, "dictionary language (if supported)")
		title    = flag.String("title", yomichan.DefaultTitle, "dictionary title")
		stride   = flag.Int("stride", yomichan.DefaultStride, "dictionary bank stride")
		pretty   = flag.Bool("pretty", yomichan.DefaultPretty, "output prettified dictionary JSON")
	)

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	if err := yomichan.ExportDb(flag.Arg(0), flag.Arg(1), *format, *language, *title, *stride, *pretty); err != nil {
		log.Fatal(err)
	}
}
