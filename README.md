# Yomichan-Import #

Yomichan Import allows users of the [Yomichan](https://foosoft.net/projects/yomichan) extension to import custom dictionary files. It currently
supports the following formats:

*   [JMdict](http://www.edrdg.org/jmdict/edict_doc.html)
*   [JMnedict](http://www.edrdg.org/enamdict/enamdict_doc.html)
*   [KANJIDIC2](http://www.edrdg.org/kanjidic/kanjd2index.html)
*   [EPWING](https://ja.wikipedia.org/wiki/EPWING)
    *       [Daijirin](https://en.wikipedia.org/wiki/Daijirin) (三省堂　スーパー大辞林)
    *       [Daijisen](https://en.wikipedia.org/wiki/Daijisen) (大辞泉)
    *       [Kenkyusha](https://en.wikipedia.org/wiki/Kenky%C5%ABsha%27s_New_Japanese-English_Dictionary) (研究社　新和英大辞典　第５版)
    *       [Meikyou](https://ja.wikipedia.org/wiki/%E6%98%8E%E9%8F%A1%E5%9B%BD%E8%AA%9E%E8%BE%9E%E5%85%B8) (明鏡国語辞典)

Yomichan Import is being expanded to support other EPWING dictionaries based on user demand. This is a mostly
non-technical and (although laborious) process that requires writing regular expressions and creating font tables;
volunteer contributions are welcome.

## Installation ##

Yomichan Import is currently available for Linux, Mac OS X, and Windows:

*   [yomichan-import_linux.tar.gz](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_linux.tar.gz)
*   [yomichan-import_darwin.tar.gz](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_darwin.tar.gz)
*   [yomichan-import_windows.zip](https://foosoft.net/projects/yomichan-import/dl/yomichan-import_windows.zip)

## Usage ##

Yomichan Import is a simple command line application. If you are a Windows user and are not comfortable using the
terminal to input commands, you can use the provided `yomichan-import.bat` batch file instead of the
`yomichan-import.exe` executable to run the application in interactive mode.

When invoked without any arguments (or executed with `--help`), the application will output usage instructions:

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

In the vast majority of cases it is enough to simply provide the path to the dictionary resource you wish to process,
without explicitly specifying a format. Yomichan Import will attempt to automatically determine the format of the
dictionary based on the contents of the path:

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

Yomichan Import will now begin the conversion process, which can take a couple of minutes to complete:

```
2016/12/29 17:12:12 converting 'dict/Kokugo/Daijirin/' to '/tmp/yomichan_tmp_825860502' in 'epwing' format...
```

After dictionary processing is complete, the tool will start a local web server to enable the Yomichan extension to
retrieve dictionary data. Users of Windows will likely see a [firewall nag dialog](https://foosoft.net/projects/yomichan-import/img/firewall.png) at this point; you
must grant network access in order to make the converted dictionary data accessible to the extension.

```
2016/12/29 17:12:20 starting dictionary server on port 9876...
```

As a final step, open the Yomichan options dialog and choose the *Local dictionary* item in the dictionary importer
drop-down menu. When you see that `http://localhost:9876/index.json` displayed in the address text-box, you can press
the *Import* button to begin the import process. Once the imported dictionary is displayed on the options screen, it is
safe to terminate the Yomichan Import tool.

## License ##

MIT
