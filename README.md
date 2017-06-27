# Yomichan Import #

Yomichan Import allows users of the [Yomichan](https://foosoft.net/projects/yomichan) extension to import custom dictionary files. It currently
supports the following formats:

*   [JMdict](http://www.edrdg.org/jmdict/edict_doc.html)
*   [JMnedict](http://www.edrdg.org/enamdict/enamdict_doc.html)
*   [KANJIDIC2](http://www.edrdg.org/kanjidic/kanjd2index.html)
*   [EPWING](https://ja.wikipedia.org/wiki/EPWING)
    *   [Daijirin](https://en.wikipedia.org/wiki/Daijirin) (三省堂　スーパー大辞林)
    *   [Daijisen](https://en.wikipedia.org/wiki/Daijisen) (大辞泉)
    *   [Kenkyusha](https://en.wikipedia.org/wiki/Kenky%C5%ABsha%27s_New_Japanese-English_Dictionary) (研究社　新和英大辞典　第５版)
    *   [Kotowaza](http://www.web-nihongo.com/wn/dictionary/dic_21/d-index.html) (故事ことわざの辞典)
    *   [Meikyou](https://ja.wikipedia.org/wiki/%E6%98%8E%E9%8F%A1%E5%9B%BD%E8%AA%9E%E8%BE%9E%E5%85%B8) (明鏡国語辞典)

Yomichan Import is being expanded to support other EPWING dictionaries based on user demand. This is a mostly
non-technical (although laborious) process that requires writing regular expressions and creating font tables; volunteer
contributions are welcome.

## Installation ##

Builds of Yomichan Import are currently available for Linux, Mac OS X, and Windows. The necessary version of
[Zero-EPWING](https://foosoft.net/projects/zero-epwing) is included for processing EPWING dictionaries.

*   [yomichan-import_linux.tar.gz](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_linux.tar.gz): (GTK+ 3 required for GUI)
*   [yomichan-import_darwin.tar.gz](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_darwin.tar.gz)
*   [yomichan-import_windows.zip](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_windows.zip) (64 bit Vista or above, no console output)

## Basic Usage ##

Please follow the steps outlined below to import your custom dictionary into Yomichan:

1.  Launch the `yomichan-import` executable.
2.  Specify the source path of the dictionary you wish to convert.
3.  Specify the target path of the dictionary ZIP archive that you wish to create.
4.  Press the button labeled *Import dictionary...* and wait for processing to complete.
5.  On the Yomichan options page, browse to the dictionary ZIP archive file you created.
6.  Wait for the import progress to complete before closing the options page.

[![Importer](https://foosoft.net/projects/yomichan-import/img/import-thumb.png)](https://foosoft.net/projects/yomichan-import/img/import.png)

## License ##

MIT
