package xflag

import (
	"flag"
	"fmt"
	"strings"
)

// Positional is a FlagSet for parsing positional arguments.
type Positional struct {
	*FlagSet
	order        []string
	require, set map[string]bool
}

func (p *Positional) PrintDefaults() {
	p.VisitAll(func(f *flag.Flag) {
		n, u := flag.UnquoteUsage(f)
		fmt.Fprintf(p.Output(), "  %s %s\n    \t%s\n", f.Name, n, u)
	})
}

// Order sets the order of args
func (p *Positional) Order(usage string) {
	p.order = strings.Split(usage, " ")
}

func (p *Positional) Parse(args []string) error {
	err := p.FlagSet.FlagSet.Parse(p.interleave(args))
	if err != nil {
		return err
	}
	p.Visit(func(f *flag.Flag) { p.set[f.Name] = true })
	for r, _ := range p.require {
		if !p.set[r] {
			return p.failf("required but not set: %s", r)
		}
	}
	return nil
}

func (p *Positional) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(p.Output(), err)
	p.Usage()
	return err
}

func (p *Positional) interleave(values []string) []string {
	p.ensure()
	args := make([]string, 0, 2*len(values))
	i := 0
	for j, v := range values {
		if i < len(p.order) {
			args = append(args, "-"+p.clean(i))
		}
		args = append(args, v)
		switch {
		case j == len(values)-1 && i == len(p.order)-1:
			i++
		case !p.isRepeating(i):
			i++
		}
	}
	return args
}

func (p *Positional) cmd() string {
	p.ensure()
	return strings.Join(p.order, " ")
}

func (p *Positional) ensure() {
	if len(p.order) == 0 {
		p.order = []string{}
		p.VisitAll(func(f *flag.Flag) {
			p.order = append(p.order, f.Name)
		})
	}

	p.require = map[string]bool{}
	p.set = map[string]bool{}
	for i := range p.order {
		if !p.isOptional(i) {
			p.require[p.clean(i)] = true
		}
	}
}

func (p *Positional) isOptional(i int) bool {
	if i >= len(p.order) {
		return true
	}
	return strings.HasPrefix(p.order[i], "[") && strings.HasSuffix(p.order[i], "]")
}

func (p *Positional) isRepeating(i int) bool {
	if i >= len(p.order) {
		return true
	}
	s := strings.TrimSuffix(p.order[i], "]")
	return strings.HasSuffix(s, "...")
}

func (p *Positional) clean(i int) string {
	if i >= len(p.order) {
		return ""
	}
	s := strings.TrimPrefix(p.order[i], "[")
	s = strings.TrimSuffix(s, "]")
	return strings.TrimSuffix(s, "...")
}
