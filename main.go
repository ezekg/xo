/*
Command xo is a command line utility similar to `sed`, except it has one job: take an
input string from `stdin` and format its matches. That's it. I built this so that I could
read configuration files and use its contents to string together a new command.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		throwError("No arguments were passed")
	}

	arg := os.Args[1]
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		throwError("Nothing passed to stdin")
	}

	delimiter := string(arg[len(arg)-1])
	pieces := delEmpty(strings.Split(arg, delimiter))
	if len(pieces) <= 1 {
		throwError("No pattern or formatter specified")
	}

	var (
		mods    string
		pattern string
		format  string
	)

	if len(pieces) > 2 {
		mods = pieces[0]
		pattern = pieces[1]
		format = pieces[2]
	} else {
		pattern = pieces[0]
		format = pieces[1]
	}

	rx, err := regexp.Compile(fmt.Sprintf(`(?%s)%s`, mods, pattern))
	if err != nil {
		throwError("Invalid regular expression", err.Error())
	}

	stdin, _ := ioutil.ReadAll(os.Stdin)
	matches := rx.FindAllSubmatch(stdin, -1)
	if matches == nil {
		throwError("No matches found")
	}

	var (
		fallbacks = make(map[int]string)
	)

	for _, match := range matches {
		var result string

		for i, m := range match {
			str := string(m)

			// Store fallback values if key does not already exist
			if _, ok := fallbacks[i]; !ok {
				rxFallback, err := regexp.Compile(fmt.Sprintf(`\$%d\?:([^\s\n]+)`, i))
				if err != nil {
					throwError("Failed to parse default arguments", err.Error())
				}

				cond := rxFallback.FindStringSubmatch(format)
				if len(cond) > 1 {
					fallbacks[i] = string(cond[1])
				}
			}

			// Set default for empty values
			if str == "" {
				str = fallbacks[i]
			}

			// Replace values
			rxRepl, _ := regexp.Compile(fmt.Sprintf(`\$%d`, i))
			if result == "" {
				result = str
			}
			result = rxRepl.ReplaceAllString(result, format)
		}

		fmt.Printf("%s\n", result)
		fmt.Println(fallbacks)
	}

	// rx, err := regexp.Compile("(?s)(.*?)\?:(.*?)")
	// if err != nil {
	// 	throwError("Invalid regular expression")
	// }
	//
	// fmt.Print(result)
}

func delEmpty(strs []string) []string {
	var result []string

	for _, str := range strs {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}

func throwError(errs ...string) {
	fmt.Println("%s\n", errs)
	os.Exit(1)
}
