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

Builds of Yomichan Import are currently available for Linux, Mac OS X, and Windows. The required version of
[Zero-EPWING](https://foosoft.net/projects/zero-epwing) is included for processing EPWING dictionaries.

*   [yomichan-import_linux.tar.gz](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_linux.tar.gz): (GTK+ 3 required for GUI)
*   [yomichan-import_darwin.tar.gz](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_darwin.tar.gz)
*   [yomichan-import_windows.zip](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_windows.zip) (64 bit Vista or above, no console output)

## Using the Graphical Interface ##

In most cases, it is sufficient to run the application without command line arguments and use the graphical interface.
Follow the steps below to import your dictionary into Yomichan:

1.  Launch the `yomichan-import` executable.
2.  Specify the path to the dictionary you wish to convert (path to `CATALOGS` file for EPWING dictionaries).
3.  Specify a network port to use (the default port `9876` should be fine for most configurations).
4.  Specify the dictionary format from the provided options.
5.  Press the button labeled *Import dictionary...* and wait for processing to complete.
6.  Once you the message `starting dictionary server on port 9876...`, the dictionary data is ready to be imported.
7.  In Yomichan, open the options page and select the *Local dictionary* item in the dictionary importer drop-down menu.
8.  When `http://localhost:9876/index.json` is displayed in the address text-box, press the *Import* button to begin import.
9.  Wait for the import progress to complete (a progress bar is displayed during dictionary processing).
9.  Close Yomichan Import once the import process has finished.

[![Import window](https://foosoft.net/projects/yomichan-import/img/import-thumb.png)](https://foosoft.net/projects/yomichan-import/img/import.png)

## Using the Command Line ##

Yomichan Import can be used as a command line application. When executed with the `--help` argument, usage instructions
will be displayed (except on Windows).

```
Usage: yomichan-import [options] input-path [output-dir]
https://foosoft.net/projects/yomichan-import/

Parameters:
  -format string
    	dictionary format [edict|enamdict|kanjidic|epwing]
  -port int
    	port to serve dictionary JSON on (default 9876)
  -pretty
    	output prettified dictionary JSON
  -serve
    	serve dictionary JSON for extension
  -stride int
    	dictionary bank stride (default 10000)
  -title string
    	dictionary title
```

In most cases it is sufficient to simply provide the path to the dictionary resource you wish to process, without
explicitly specifying a format. Yomichan Import will attempt to automatically determine the format of the dictionary
based on the contents of the path:

| Format       | Resource                             |
| ------------ | ------------------------------------ |
| **edict**    | file named `JMDict_e.xml`            |
| **enamdict** | file named `JMNedict.xml`            |
| **kanjidic** | file named `kanjidic2.xml`           |
| **epwing**   | directory with file named `CATALOGS` |

For example, if you wanted to process an EPWING dictionary titled Daijirin, you could do so with the following command
(shown on Linux):

```
$ ./yomichan-import dict/Kokugo/Daijirin/
```

Yomichan Import will now begin the conversion process, which can take a couple of minutes to complete. Once you see the
message `starting dictionary server on port 9876...` output to your console, you can use Yomichan to import the
processed dictionary data using the same steps as described in the *Using the Graphical Interface* section.

## License ##

MIT
