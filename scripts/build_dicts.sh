#!/bin/bash

go get foosoft.net/projects/yomichan-import/yomichan

mkdir -p src
mkdir -p dst

function refresh_source () {
    NOW=$(date '+%s')
    YESTERDAY=$((NOW - 86400)) # 86,400 seconds in 24 hours
    if [ ! -f "src/$1" ]; then
        wget "ftp.edrdg.org/pub/Nihongo/$1.gz"
        gunzip -c "$1.gz" > "src/$1"
    elif [[ $YESTERDAY -gt $(date -r "src/$1" '+%s') ]]; then
        rsync "ftp.edrdg.org::nihongo/$1" "src/$1"
    fi
}

refresh_source "JMdict_e_examp"
yomichan -language="english" -title="JMdict" src/JMdict_e_examp dst/jmdict_english_with_examples.zip

refresh_source "JMdict"
yomichan -language="english"   -title="JMdict"             src/JMdict dst/jmdict_english.zip
yomichan -language="dutch"     -title="JMdict (Dutch)"     src/JMdict dst/jmdict_dutch.zip
yomichan -language="french"    -title="JMdict (French)"    src/JMdict dst/jmdict_french.zip
yomichan -language="german"    -title="JMdict (German)"    src/JMdict dst/jmdict_german.zip
yomichan -language="hungarian" -title="JMdict (Hungarian)" src/JMdict dst/jmdict_hungarian.zip
yomichan -language="russian"   -title="JMdict (Russian)"   src/JMdict dst/jmdict_russian.zip
yomichan -language="slovenian" -title="JMdict (Slovenian)" src/JMdict dst/jmdict_slovenian.zip
yomichan -language="spanish"   -title="JMdict (Spanish)"   src/JMdict dst/jmdict_spanish.zip
yomichan -language="swedish"   -title="JMdict (Swedish)"   src/JMdict dst/jmdict_swedish.zip

yomichan -format="forms"       -title="JMdict Forms"       src/JMdict dst/jmdict_forms.zip

refresh_source "JMnedict.xml"
yomichan src/JMnedict.xml dst/jmnedict.zip

refresh_source "kanjidic2.xml"
yomichan -language="english"    -title="KANJIDIC"              src/kanjidic2.xml dst/kanjidic_english.zip
yomichan -language="french"     -title="KANJIDIC (French)"     src/kanjidic2.xml dst/kanjidic_french.zip
yomichan -language="portuguese" -title="KANJIDIC (Portuguese)" src/kanjidic2.xml dst/kanjidic_portuguese.zip
yomichan -language="spanish"    -title="KANJIDIC (Spanish)"    src/kanjidic2.xml dst/kanjidic_spanish.zip
