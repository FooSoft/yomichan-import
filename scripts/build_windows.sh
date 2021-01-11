#!/bin/bash
mkdir -p dst
mkdir -p yomichan-import

export CXX=x86_64-w64-mingw32-g++.exe
export CC=x86_64-w64-mingw32-gcc.exe
go build github.com/FooSoft/yomichan-import/yomichan
go build -ldflags="-H windowsgui" github.com/FooSoft/yomichan-import/yomichan-gtk

mv yomichan.exe yomichan-import
mv yomichan-gtk.exe yomichan-import

7za a yomichan-import_windows.zip yomichan-import

rm -rf yomichan-import
