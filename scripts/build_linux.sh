#!/bin/bash
go build github.com/FooSoft/yomichan-import/yomichan
go build github.com/FooSoft/yomichan-import/yomichan-gtk

mkdir yomichan-import

mv yomichan yomichan-import
mv yomichan-gtk yomichan-import

7za a yomichan-import_linux.7z yomichan-import

rm -rf yomichan-import
