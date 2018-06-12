package xflag

import (
	"flag"
	"reflect"
	"testing"
)

func TestBuffer(t *testing.T) {
	var ss = []struct {
		i []string
		x string
	}{
		{[]string{}, ""},
		{[]string{"-b", ""}, "\n"},
		{[]string{"-b", "foo", "-b", "baz"}, "foo\nbaz\n"},
	}
	for _, s := range ss {
		f := &FlagSet{FlagSet: &flag.FlagSet{}}
		b := f.Buffer("b", "")
		f.Parse(s.i)
		r := b.String()
		if r != s.x {
			t.Errorf("expected %q, got %q", s.x, r)
		}
	}
}

func TestStrings(t *testing.T) {
	var ss = []struct {
		i []string
		x []string
	}{
		{[]string{}, []string{}},
		{[]string{"-s", ""}, []string{""}},
		{[]string{"-s", "foo", "-s", "baz"}, []string{"foo", "baz"}},
	}
	for _, s := range ss {
		f := &FlagSet{FlagSet: &flag.FlagSet{}}
		r := f.Strings("s", "")
		f.Parse(s.i)
		if !reflect.DeepEqual(*r, s.x) {
			t.Errorf("expected %q, got %q", s.x, *r)
		}
	}
}
