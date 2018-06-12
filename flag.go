// Package xflag is a collection of extra flags for the flag package.
package xflag

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// CommandLine is a FlagSet for flag.CommandLine.
var CommandLine = &FlagSet{FlagSet: flag.CommandLine}

// FlagSet is a wrapper for flag.FlagSet.
type FlagSet struct {
	*flag.FlagSet
	pos *Positional
}

func Pos() *Positional {
	return CommandLine.Pos()
}

func (f *FlagSet) Pos() *Positional {
	if f.pos == nil {
		f.pos = &Positional{FlagSet: &FlagSet{FlagSet: &flag.FlagSet{}}}
	}
	return f.pos
}

type bufferValue struct {
	*bytes.Buffer
}

// Set writes arg to the buffer with a newline.
func (b *bufferValue) Set(arg string) error {
	_, err := b.WriteString(arg + "\n")
	return err
}

// Buffer defines a bytes.Buffer flag with specified name and usage string.
// The return value is the address of a bytes.Buffer that stores the value of the flag.
// A newline character is written after each value that is stored.
func Buffer(name, usage string) *bytes.Buffer {
	return CommandLine.Buffer(name, usage)
}

// Buffer defines a bytes.Buffer flag with specified name and usage string.
// The return value is the address of a bytes.Buffer that stores the value of the flag.
// A newline character is written after each value that is stored.
func (f *FlagSet) Buffer(name, usage string) *bytes.Buffer {
	b := bytes.NewBuffer([]byte{})
	f.BufferVar(b, name, usage)
	return b
}

// BufferVar defines a bytes.buffer flag with specified name and usage string.
// The argument p points to a bytes.Buffer variable in which to store the value of the flag.
// A newline character is written after each value that is stored.
func BufferVar(b *bytes.Buffer, name, usage string) {
	CommandLine.Var(&bufferValue{b}, name, usage)
}

// BufferVar defines a bytes.buffer flag with specified name and usage string.
// The argument p points to a bytes.Buffer variable in which to store the value of the flag.
// A newline character is written after each value that is stored.
func (f *FlagSet) BufferVar(b *bytes.Buffer, name, usage string) {
	f.Var(&bufferValue{b}, name, usage)
}

type stringsValue struct {
	p *[]string
}

// Set appends arg to the slice of strings.
func (s *stringsValue) Set(arg string) error {
	*s.p = append(*s.p, arg)
	return nil
}

func (s *stringsValue) String() string {
	if s.p == nil {
		return ""
	}
	return fmt.Sprint(*s.p)
}

// Strings defines a slice of strings flag with specified name and usage string.
// The return value is the address of a string slice that stores the value of the flag.
func Strings(name, usage string) *[]string {
	return CommandLine.Strings(name, usage)
}

// Strings defines a slice of strings flag with specified name and usage string.
// The return value is the address of a string slice that stores the value of the flag.
func (f *FlagSet) Strings(name, usage string) *[]string {
	ss := []string{}
	f.StringsVar(&ss, name, usage)
	return &ss
}

// StringsVar defines a slice of strings flag with specified name and usage string.
// The argument p points to a string slice in which to store the value of the flag.
func StringsVar(p *[]string, name, usage string) {
	CommandLine.StringsVar(p, name, usage)
}

// StringsVar defines a slice of strings flag with specified name and usage string.
// The argument p points to a string slice in which to store the value of the flag.
func (f *FlagSet) StringsVar(p *[]string, name, usage string) {
	f.Var(&stringsValue{p: p}, name, usage)
}

type writerValue struct {
	io.Writer
}

// Set appends arg to the slice of strings.
func (w *writerValue) Set(arg string) error {
	_, err := w.Write([]byte(arg))
	return err
}

func (w *writerValue) String() string {
	return "writer"
}

// WriterVar defines a writer flag with specified name and usage string.
// The argument w is a writer in which to write the value of the flag.
func WriterVar(w io.Writer, name, usage string) {
	CommandLine.WriterVar(w, name, usage)
}

// WriterVar defines a writer flag with specified name and usage string.
// The argument w is a writer in which to write the value of the flag.
func (f *FlagSet) WriterVar(w io.Writer, name, usage string) {
	f.Var(&writerValue{w}, name, usage)
}

type InFileValue struct {
	io.ReadCloser
}

func inOrStdin(arg string) (io.ReadCloser, error) {
	if arg == "-" {
		return ioutil.NopCloser(os.Stdin), nil
	}
	return os.Open(arg)
}

func (i *InFileValue) Set(arg string) (err error) {
	i.ReadCloser, err = inOrStdin(arg)
	return err
}

func (i *InFileValue) String() string {
	return "input file"
}

func (i *InFileValue) Read(b []byte) (int, error) {
	if i.ReadCloser == nil {
		return 0, io.EOF
	}
	return i.ReadCloser.Read(b)
}

func (i *InFileValue) Close() error {
	if i.ReadCloser == nil {
		return nil
	}
	return i.ReadCloser.Close()
}

// InFile defines a slice of io.ReadCloser flag with specified name and usage string.
// The return value is the address of a slice that stores the value of the flag.
func InFile(name, usage string) io.ReadCloser {
	return CommandLine.InFile(name, usage)
}

// InFile defines a slice of io.ReadCloser flag with specified name and usage string.
// The return value is the address of a slice that stores the value of the flag.
func (f *FlagSet) InFile(name, usage string) io.ReadCloser {
	r := &InFileValue{}
	f.InFileVar(r, name, usage)
	return r
}

// InFileVar defines a slice of io.ReadCloser flag with specified name and usage string.
// The argument p is a pointer to a slice in which to store the value of the flag.
func InFileVar(p *InFileValue, name, usage string) {
	CommandLine.InFileVar(p, name, usage)
}

// InFileVar defines a slice of io.ReadCloser flag with specified name and usage string.
// The argument p is a pointer to a slice in which to store the value of the flag.
func (f *FlagSet) InFileVar(p *InFileValue, name, usage string) {
	f.Var(&InFileValue{p}, name, usage)
}

type inFilesValue struct {
	p *[]io.ReadCloser
}

func (i *inFilesValue) Set(arg string) error {
	r, err := inOrStdin(arg)
	*i.p = append(*i.p, r)
	return err
}

func (i *inFilesValue) String() string {
	return "input files"
}

// InFiles defines a slice of io.ReadCloser flag with specified name and usage string.
// The return value is the address of a slice that stores the value of the flag.
func InFiles(name, usage string) *[]io.ReadCloser {
	return CommandLine.InFiles(name, usage)
}

// InFiles defines a slice of io.ReadCloser flag with specified name and usage string.
// The return value is the address of a slice that stores the value of the flag.
func (f *FlagSet) InFiles(name, usage string) *[]io.ReadCloser {
	rs := []io.ReadCloser{}
	f.InFilesVar(&rs, name, usage)
	return &rs
}

// InFilesVar defines a slice of io.ReadCloser flag with specified name and usage string.
// The argument p is a pointer to a slice in which to store the value of the flag.
func InFilesVar(p *[]io.ReadCloser, name, usage string) {
	CommandLine.InFilesVar(p, name, usage)
}

// InFilesVar defines a slice of io.ReadCloser flag with specified name and usage string.
// The argument p is a pointer to a slice in which to store the value of the flag.
func (f *FlagSet) InFilesVar(p *[]io.ReadCloser, name, usage string) {
	f.Var(&inFilesValue{p}, name, usage)
}

type outFilesValue struct {
	p *[]io.WriteCloser
}

func (i *outFilesValue) Set(arg string) error {
	if arg == "-" {
		*i.p = append(*i.p, &nopWriteCloser{os.Stdout})
		return nil
	}
	f, err := os.Create(arg)
	if err != nil {
		return err
	}

	*i.p = append(*i.p, f)
	return nil
}

func (i *outFilesValue) String() string {
	return "output files"
}

// OutFiles defines a slice of io.WriteCloser flag with specified name and usage string.
// The return value is the address of a slice that stores the value of the flag.
func OutFiles(name, usage string) *[]io.WriteCloser {
	return CommandLine.OutFiles(name, usage)
}

// OutFiles defines a slice of io.WriteCloser flag with specified name and usage string.
// The return value is the address of a slice that stores the value of the flag.
func (f *FlagSet) OutFiles(name, usage string) *[]io.WriteCloser {
	ws := []io.WriteCloser{}
	f.OutFilesVar(&ws, name, usage)
	return &ws
}

// OutFilesVar defines a slice of io.WriteCloser flag with specified name and usage string.
// The argument p is a pointer to a slice in which to store the value of the flag.
func OutFilesVar(p *[]io.WriteCloser, name, usage string) {
	CommandLine.OutFilesVar(p, name, usage)
}

// OutFilesVar defines a slice of io.WriteCloser flag with specified name and usage string.
// The argument p is a pointer to a slice in which to store the value of the flag.
func (f *FlagSet) OutFilesVar(p *[]io.WriteCloser, name, usage string) {
	f.Var(&outFilesValue{p}, name, usage)
}

func Parse() error {
	return CommandLine.Parse(os.Args[1:])
}

func (f *FlagSet) Parse(args []string) error {
	f.Usage = f.defaultUsage
	if err := f.FlagSet.Parse(args); err != nil {
		return err
	}

	if f.pos != nil {
		f.pos.FlagSet.Usage = f.defaultUsage
		if err := f.pos.Parse(f.Args()); err != nil {
			return f.handle(err)
		}
	}
	return nil
}

func (f *FlagSet) handle(err error) error {
	if err == nil {
		return err
	}
	switch f.ErrorHandling() {
	case flag.ContinueOnError:
		return err
	case flag.ExitOnError:
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
	return err
}

func (f *FlagSet) defaultUsage() {
	u := fmt.Sprintf("Usage: %s", f.Name())
	n := 0
	flag.VisitAll(func(f *flag.Flag) { n++ })

	if n > 0 {
		u += " [options]"
	}
	if f.pos != nil {
		u += " " + f.pos.cmd()
	}
	fmt.Fprintln(f.Output(), u)
	if f.pos != nil {
		f.pos.PrintDefaults()
	}
	if n > 0 {
		fmt.Fprintln(f.Output(), "Options:")
	}
	f.PrintDefaults()
}

type nopWriteCloser struct {
	io.Writer
}

func (w *nopWriteCloser) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}

func (w *nopWriteCloser) Close() error {
	return nil
}
