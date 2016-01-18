package main

import (
	"bytes"
	"os/exec"
	"syscall"
	"testing"
)

func TestMain(t *testing.T) {
	shouldEqual(t, `xo`,
		`Usage: xo '/<pattern>/<formatter>/[flags]'
`)
	shouldEqual(t, `echo 'Hello there!' | xo '~hello(.*)~Hi$1~i'`,
		`Hi there!
`)
	shouldEqual(t, `echo 'Hello! - Luke' | xo '/(hello).*?-.*?(\w+)/Why $1, $2!/i'`,
		`Why Hello, Luke!
`)
	shouldEqual(t, `cat fixtures/servers.yml | xo '/.*?(production):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/$4@$2 -p $3?:22/mis'`,
		`user-1@192.168.1.1 -p 1234
`)
	shouldEqual(t, `cat fixtures/starwars.txt | xo '/^(\w+):(\s*\[(.*?)\]\s*)?\s*([^\n]+)/$1 said, "$4" in a $3?:normal voice./mi'`,
		`Vader said, "If only you knew the power of the Dark Side. Obi-Wan never told you what happened to your father." in a normal voice.
Luke said, "He told me enough! He told me you killed him!" in a normal voice.
Vader said, "No, I am your father." in a normal voice.
Luke said, "No. No! That's not true! That's impossible!" in a shocked voice.
`)
	shouldEqual(t, `echo '123' | xo '/(\d)(\d)(\d)(\d)?(\d)?/$1, $2, $3, 4?:FOUR $5?:FIVE/'`,
		`1, 2, 3, 4?:FOUR FIVE
`)
	shouldEqual(t, `echo 'abc' | xo '%(\w)(\w)(\w)(\w)?%$1$2$3$4?:$1%'`,
		`abca
`)
	shouldEqual(t, `echo 'Hello! My name is C3PO, human cyborg relations.' | xo '/^((\w+)! )?my name is (\w+)/$2?:Hello, $3!/i'`,
		`Hello, C3PO!
`)
	shouldEqual(t, `echo 'My name is Chewbacca, uuuuuur ahhhhhrrr uhrrr ahhhrrr aaargh!' | xo '|^((\w+)! )?my name is (\w+)|$2?:Greetings, $3!|i'`,
		`Greetings, Chewbacca!
`)
	shouldEqual(t, `cat fixtures/romans.txt | xo '/\d\s(\w+).*?to all that are in (\w+),.*?24 \[the (grace)? of ([\w\s]{21})/Romans is a letter written by $1 addressed to the people of $2 about the $3?:gospel of $4./mis'`,
		`Romans is a letter written by Paul addressed to the people of Rome about the grace of our Lord Jesus Christ.
`)
	shouldEqual(t, `echo 'hi' | xo '/(hi)/te\/st/mi'`,
		`te/st
`)
	shouldExit(t, `echo '1' | xo '/^(\s)/$1/'`, 1)
	shouldExit(t, `echo '1' | xo '/1/'`, 1)
	shouldExit(t, `echo '1' | xo ///`, 1)
	shouldExit(t, `xo ///`, 1)
}

func TestSplit(t *testing.T) {
	tests := map[string][]string{
		`%bc%b\%%`:    []string{"bc", "b%"},
		`⌘abc⌘bca⌘`:   []string{"abc", "bca"},
		`⌘abc⌘bca⌘\⌘`: []string{"abc", "bca", "⌘"},
		`\bc\bc\`:     []string{"bc", "bc"},
		`\b\\c\bc\`:   []string{`b\c`, `bc`},
		`[\[xy[xy[`:   []string{"[xy", "xy"},
		`[\\[xy[xy[`:  []string{`\[xy`, "xy"},
		`[\\[xy[xy[i`: []string{`\[xy`, "xy", "i"},
		`///`:         []string{},
		`///a`:        []string{"a"},
		``:            []string{},
	}
outer:
	for test, expected := range tests {
		actual, err := split(test)
		if err != nil {
			t.Fatalf("failed to split `%q`\n", test)
			continue
		}

		if len(actual) != len(expected) {
			t.Fatalf("`%v` should be `%v` for `%q`\n", actual, expected, test)
			continue
		}

		for i := range actual {
			if actual[i] != expected[i] {
				t.Fatalf("`%v` should be `%v` for `%q`\n", actual, expected, test)
				continue outer
			}
		}
	}
}

func execShellCommand(t *testing.T, cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	return string(out)
}

func shouldEqual(t *testing.T, cmd string, expected string) {
	result := execShellCommand(t, cmd)
	if result != expected {
		t.Fatalf("`%s` should be `%s`", result, expected)
	}
}

func shouldExit(t *testing.T, cmd string, expected int) {
	c := exec.Command("bash", "-c", cmd)
	stderr := &bytes.Buffer{}
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		code := 0
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code = status.ExitStatus()
			}
		}
		if code != expected {
			t.Fatalf("exit status `%d` should be `%d`", code, expected)
		}
	} else {
		t.Fatalf("command `%s` failed to spawn", cmd)
	}
}
