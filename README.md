# xflag

a sidecar for go's flag package

```
go get -d github.com/dedelala/xflag
```

## example

```
s := flag.String("s", "idk", "some string for something")
ss := xflag.Strings("ss", "many strings for something else")
flag.Parse()
```

## funky flags

Append multiple args into a newline separated buffer.
```
bb := xflag.Buffer("bb", "something supplied many times")
```

Open files for reading or "-" for stdin.
```
i := xflag.InFiles("i", "input files")
```

Open files for writing or "-" for stdout.
```
o := xflag.OutFiles("o", "output files")
```

## positionals

```
s := xflag.Pos().String("arg", "foo", "the main arg for the program")
// the xflag parse method is required to parse positionals, otherwise flag.Parse will be fine.
xflag.Parse()
```
