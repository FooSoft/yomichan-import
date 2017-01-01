#!/bin/sh
gox -os="linux windows" -arch="386"

mkdir -p yomichan-import_windows/bin/windows
mv yomichan-import_windows_386.exe yomichan-import_windows/yomichan-import.exe
cp bin/windows/* yomichan-import_windows/bin/windows/
7z a yomichan-import_windows.zip yomichan-import_windows
rm -rf yomichan-import_windows

mkdir -p yomichan-import_linux/bin/linux
mv yomichan-import_linux_386 yomichan-import_linux/yomichan-import
cp bin/linux/* yomichan-import_linux/bin/linux/
tar czvf yomichan-import_linux.tar.gz yomichan-import_linux
rm -rf yomichan-import_linux
