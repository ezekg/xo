/*
Command xo is a command line utility that takes an input string from stdin and
formats the regexp matches.
*/
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"unicode/utf8"
)

func main() {
	if len(os.Args) == 1 {
		help()
	}

	arg := os.Args[1]
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		throw("Nothing passed to stdin")
	}

	parts, err := split(arg)
	if err != nil {
		throw("Invalid argument string")
	}
	if len(parts) <= 1 {
		throw("No pattern or formatter specified")
	}
	if len(parts) > 3 {
		throw("Extra delimiter detected (maybe try one other than `/`)")
	}

	var (
		flags   string
		pattern string
		format  string
	)

	pattern = parts[0]
	format = parts[1]
	if len(parts) > 2 {
		flags = parts[2]
	}

	rx, err := regexp.Compile(fmt.Sprintf(`(?%s)%s`, flags, pattern))
	if err != nil {
		throw("Invalid regular expression")
	}

	in, _ := ioutil.ReadAll(os.Stdin)
	matches := rx.FindAllSubmatch(in, -1)
	if matches == nil {
		throw("No matches found")
	}

	fallbacks := make(map[int]string)

	for _, group := range matches {
		result := format

		for i, match := range group {
			value := string(match)

			rxFallback, err := regexp.Compile(fmt.Sprintf(`(\$%d)\?:([-_$A-za-z1-9]+)`, i))
			if err != nil {
				throw("Failed to parse default arguments")
			}

			fallback := rxFallback.FindStringSubmatch(result)
			if len(fallback) > 1 {
				// Store fallback values if key does not already exist
				if _, ok := fallbacks[i]; !ok {
					fallbacks[i] = string(fallback[2])
				}
				result = rxFallback.ReplaceAllString(result, "$1")
			}

			// Set default for empty values
			if value == "" {
				value = fallbacks[i]
			}

			// Replace values
			rxRepl, _ := regexp.Compile(fmt.Sprintf(`\$%d`, i))
			result = rxRepl.ReplaceAllString(result, value)
		}

		fmt.Printf("%s\n", result)
	}
}

// split slices str into all substrings separated by non-escaped values of the
// first rune and returns a slice of those substrings.
// It removes one backslash escape from any escaped delimiters.
func split(str string) ([]string, error) {
	if !utf8.ValidString(str) {
		return nil, errors.New("Invalid string")
	}

	// Grab the first rune.
	delim, size := utf8.DecodeRuneInString(str)
	str = str[size:]

	var (
		subs   []string
		buffer bytes.Buffer
	)
	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		str = str[size:]

		if r == '\\' {
			peek, peekSize := utf8.DecodeRuneInString(str)
			if peek == delim {
				buffer.WriteRune(peek)
				str = str[peekSize:]
				continue
			}
		}

		if r == delim {
			if buffer.Len() > 0 {
				subs = append(subs, buffer.String())
				buffer = *new(bytes.Buffer)
			}
			continue
		}

		buffer.WriteRune(r)
	}
	if buffer.Len() > 0 {
		subs = append(subs, buffer.String())
	}
	return subs, nil
}

// help prints how to use this whole `xo` thing then gracefully exits.
func help() {
	fmt.Printf("%s\n", "Usage: xo '/<pattern>/<formatter>/[flags]'")
	os.Exit(0)
}

// throw prints a bunch of strings and then exits with a non-zero exit code.
func throw(errs ...string) {
	for _, err := range errs {
		fmt.Printf("%s\n", err)
	}
	os.Exit(1)
}
