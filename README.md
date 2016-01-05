Babysit
=======

## What Do?

This is a very simple wrapper which runs a target program with a set of input files and logs non-zero exit codes. This approach _is_ sufficient to catch exits that won't cause a debugger break-in like /GS failures and FailFast. In combination with a simple target that processes one input and then exits, you can use this as a fairly fast first cut across a number of files, and then triage further with something like [BugId](https://github.com/SkyLined/BugId). Under the hood it just uses `CreateProcess`, `WaitForSingleObject` and `GetExitCodeProcess`.

## Documentation

```
Usage: bin\babysit -i c:\inputs\*.txt C:\path\to\target -in @@ -other -targetopts
( @@ will be substituted with each input )

  -i string
        Glob for input files, eg: c:\files\*.pdf
  -t int
        Process wait timeout (ms), -1 for INFINITE (default 10000)
```

Using powershell, you can filter the output with something like `babysit -i [ARGS TO BABYSIT] 2>&1 | select-string -match "0xC0000409" | out-file "fun.txt"`. (Note the redirection of stderr into stdout there)

## Installation

1. You should follow the [instructions](https://golang.org/doc/install) to
install Go, if you haven't already done so.

2. You'll probably want a way for `go get` to work. There are a few ways to achieve this on Windows. Personally, I use the [github client](https://desktop.github.com/) which also ships with a "Git Shell" powershell shortcut and some handy binaries.

3. To build this particular binary, cgo needs to find a compiler it likes. The simplest way to get that is via something like [mingw-w64](http://mingw-w64.org/doku.php) which will give you gcc.

4. To make stuff easier you might want to set your system PATH so it can find `git` and `gcc` among other things. Mine ends with `C:\Go\bin;C:\Program Files\mingw-w64\x86_64-5.2.0-win32-seh-rt_v4-rev0\mingw64\bin;C:\Users\ben\AppData\Local\GitHub\PortableGit_c2ba306e536fdf878271f7fe636a147ff37326ad\bin`

5. Once you have done all that, `go get github.com/bnagy/babysit` should download, build and install `babysit.exe` into your `%GOPATH%\bin`

*IF ALL ELSE FAILS* and you are a maniac, there is a pre-built x64 binary in the `/bin` directory. You could download [master](https://github.com/bnagy/babysit/archive/master.zip) and just run that, but I will respect you less.

## Screenshots

```
2016/01/04 16:25:33 [OK] babysit starting up...
2016/01/04 16:25:33 [OK] Found 208 input files.
2016/01/04 16:25:37 [!!] "C:\\jxrs\\Maui-12bppYCC420.jxr" 0x88982f62
2016/01/04 16:25:37 [!!] "C:\\jxrs\\Maui-20bppYCC422.jxr" 0x88982f62
2016/01/04 16:25:37 [!!] "C:\\jxrs\\Maui-24bppYCC444.jxr" 0x88982f62
2016/01/04 16:25:37 [!!] "C:\\jxrs\\Maui-32bppCMYKDIRECT.jxr" 0x88982f62
2016/01/04 16:25:37 [!!] "C:\\jxrs\\Maui-32bppYCC422.jxr" 0x88982f62
2016/01/04 16:25:43 [OK] All done. 208 files in 9.4809055s (21.94/s)

## TODO

Nothing planned. Open issues for feature requests.

## Contributing

Fork and send a pull request to contribute code.

Report issues.

## License

BSD style, see LICENSE file for details.
