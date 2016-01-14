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
	shouldEqual(t, `echo 'Hello there!' | xo '/hello(.*)/Hi$1/i'`,
		`Hi there!
`)
	shouldEqual(t, `echo 'Hello! - Bob' | xo '/(hello).*?-.*?(\w+)/Why $1, $2!/i'`,
		`Why Hello, Bob!
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
	shouldEqual(t, `echo 'abc' | xo '/(\w)(\w)(\w)(\w)?/$1$2$3$4?:$1/'`,
		`abca
`)
	shouldEqual(t, `echo ',2,3' | xo '/^(\d)?,(\d),(\d)/$3,$2,$1?:$3/'`,
		`3,2,
`)
	shouldExit(t, `echo '1' | xo '/^(\s)/$1/'`, 1)
	shouldExit(t, `echo '1' | xo '/1/'`, 1)
	shouldExit(t, `echo '1' | xo ///`, 1)
	shouldExit(t, `xo ///`, 1)
}

func execShellCommand(t *testing.T, cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		t.Fatalf(err.Error())
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
