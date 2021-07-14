package elio

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/maruel/panicparse/stack"
)

// PanicParse panic parse
func PanicParse() (string, error) {
	stack := string(debug.Stack())
	out, err := StackParse(stack)
	if nil != err {
		return stack, err
	}
	return out, err
}

// StackParse stack parse
func StackParse(in string) (out string, err error) {
	i := bytes.NewBufferString(in)
	o := &bytes.Buffer{}
	var c *stack.Context
	if c, err = stack.ParseDump(i, o, true); nil != err {
		fmt.Printf("error:\n%s", err.Error())
		return "", err
	}

	// Find out similar goroutine traces and group them into buckets.
	buckets := stack.Aggregate(c.Goroutines, stack.AnyValue)

	// Calculate alignment.
	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.SrcLine()); l > srcLen {
				srcLen = l
			}
			if l := len(line.Func.PkgName()); l > pkgLen {
				pkgLen = l
			}
		}
	}

	var builder strings.Builder
	for _, bucket := range buckets {
		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}
		if c := bucket.CreatedByString(false); c != "" {
			extra += " [Created by " + c + "]"
		}
		builder.WriteString(fmt.Sprintf("%d: %s%s\n", len(bucket.IDs), bucket.State, extra))

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			builder.WriteString(fmt.Sprintf("    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.PkgName(), srcLen, line.SrcLine(),
				line.Func.Name(), &line.Args))
		}
		if bucket.Stack.Elided {
			builder.WriteString("    (...)\n")
		}
	}

	return builder.String(), nil
}
