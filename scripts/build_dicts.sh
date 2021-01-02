#!/bin/bash

go get github.com/FooSoft/yomichan-import/yomichan

mkdir -p src
mkdir -p dst

if [ ! -f src/JMdict ]; then
    wget http://ftp.monash.edu/pub/nihongo/JMdict.gz
    gunzip -c JMdict.gz > src/JMdict
fi

yomichan -language="dutch"     -title="JMdict (Dutch)"     src/JMdict dst/jmdict_dutch.zip
yomichan -language="english"   -title="JMdict (English)"   src/JMdict dst/jmdict_english.zip
yomichan -language="french"    -title="JMdict (French)"    src/JMdict dst/jmdict_french.zip
yomichan -language="german"    -title="JMdict (German)"    src/JMdict dst/jmdict_german.zip
yomichan -language="hungarian" -title="JMdict (Hungarian)" src/JMdict dst/jmdict_hungarian.zip
yomichan -language="russian"   -title="JMdict (Russian)"   src/JMdict dst/jmdict_russian.zip
yomichan -language="slovenian" -title="JMdict (Slovenian)" src/JMdict dst/jmdict_slovenian.zip
yomichan -language="spanish"   -title="JMdict (Spanish)"   src/JMdict dst/jmdict_spanish.zip
yomichan -language="swedish"   -title="JMdict (Swedish)"   src/JMdict dst/jmdict_swedish.zip

if [ ! -f src/JMnedict.xml ]; then
    wget http://ftp.monash.edu/pub/nihongo/JMnedict.xml.gz
    gunzip -c JMnedict.xml.gz > src/JMnedict.xml
fi

yomichan src/JMnedict.xml dst/jmnedict.zip

if [ ! -f src/kanjidic2.xml ]; then
    wget http://www.edrdg.org/kanjidic/kanjidic2.xml.gz
    gunzip -c kanjidic2.xml.gz > src/kanjidic2.xml
fi

yomichan -language="english"    -title="KANJIDIC (English)"    src/kanjidic2.xml dst/kanjidic_english.zip
yomichan -language="french"     -title="KANJIDIC (French)"     src/kanjidic2.xml dst/kanjidic_french.zip
yomichan -language="portuguese" -title="KANJIDIC (Portuguese)" src/kanjidic2.xml dst/kanjidic_portuguese.zip
yomichan -language="spanish"    -title="KANJIDIC (Spanish)"    src/kanjidic2.xml dst/kanjidic_spanish.zip
