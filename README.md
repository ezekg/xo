# `xo`
[![Travis](https://img.shields.io/travis/ezekg/xo.svg?style=flat-square)](https://travis-ci.org/ezekg/xo)

`xo` is a command line utility that takes an input string from stdin and formats
the regexp matches. You might immediately think that this is a knockoff of `sed`,
but `xo` has only one job: to format matches. It does not handle search/replace,
because there are better tools for that, and again, it _only cares about matches_.
Unlike `sed`, it comes with the full power of [Go's regular expression syntax](https://golang.org/pkg/regexp/syntax/),
meaning it can handle multiline patterns and any other flag you can throw at it.

Enjoy.

## Installation
To install `xo`, please use `go get`. If you don't have Go installed, [get it here](https://golang.org/dl/).
If you would like to grab a precompiled binary, head over to the [releases](https://github.com/ezekg/xo/releases)
page. The precompiled `xo` binaries have no external dependencies.

```
go get github.com/ezekg/xo
```

## Usage
`xo` accepts the following syntax; all you have to do is feed it some `stdin` via
piped output (`echo 'hello' | xo ...`) or what have you. There's no flags, and no
additional arguments. Simple and easy to use.
```
xo '/<pattern>/<formatter>/[flags]'
```

## Examples
Let's start off a little simple, and then we'll ramp it up and get crazy. `xo`,
in its simplest form, does things like this,
```bash
echo 'Hi! My name is Bob.' | xo '/^(\w+)?! my name is (\w+)/$1, $2!/i'
# =>
#  Hi, Bob!
```

With that, what if the input string _forgot_ to specify a greeting, but we still want
to say "Hello"? Well, that sounds like a great job for a [fallback value](#fallback-values)!
Let's update the example a little bit,
```bash
echo 'Hi! My name is Bob.' | xo '/^(\w+)?! my name is (\w+)/$1?:Hello, $2!/i'
# =>
#  Hi, Bob!

echo 'My name is Sara.' | xo '/^((\w+)! )?my name is (\w+)/$2?:Hello, $3!/i'
# =>
#  Hello, Sara!
```

As you can see, we've taken the matches and created a new string out of them. We
also supplied a [fallback value](#fallback-values) for the second match (`$2`)
using the ternary `?:` operator that gets used if no match is found.

When you create a regular expression, wrapping a subexpression in parenthesis `(...)`
creates a new _capturing group_, numbered from left to right in order of opening
parenthesis. Submatch `$0` is the match of the entire expression, submatch `$1`
the match of the first parenthesized subexpression, and so on. These capturing
groups are what `xo` works with.

Now, suppose we had a text file called `starwars.txt` containing some Star Wars quotes,
```
Vader: If only you knew the power of the Dark Side. Obi-Wan never told you what happened to your father.
Luke: He told me enough! He told me you killed him!
Vader: No, I am your father.
Luke: [shocked] No. No! That's not true! That's impossible!
```

and we wanted to do a little formatting, as if we're telling it as a story. Easy!
```bash
cat starwars.txt | xo '/^(\w+):(\s*\[(.*?)\]\s*)?\s*([^\n]+)/$1 said, "$4" in a $3?:normal voice./mi'
# =>
#   Vader said, "If only you knew the power of the Dark Side. Obi-Wan never told you what happened to your father." in a normal voice.
#   Luke said, "He told me enough! He told me you killed him!" in a normal voice.
#   Vader said, "No, I am your father." in a normal voice.
#   Luke said, "No. No! That's not true! That's impossible!" in a shocked voice.
```

Okay, okay. Let's move on to something a little more useful. Suppose we had a
configuration file called `servers.yml` containing some project information.
Maybe it looks like this,
```yml
stages:
  production:
    server: 192.168.1.1:1234
    user: user-1
  staging:
    server: 192.168.1.1
    user: user-2
```

Now, let's say we have one of these configuration files for every project we've ever
worked on. Our day to day requires us to SSH into these projects a lot, and having
to read the config file for the IP address of the server, the SSH user, as well as
any potential port number gets pretty repetitive. Let's automate!
```bash
cat servers.yml | xo '/.*?(production):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/$4@$2 -p $3?:22/mis'
# =>
#  user-1@192.168.1.1 -p 1234

# Now let's actually use the output,
ssh $(cat servers.yml | xo '/.*?(staging):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/$4@$2 -p $3?:22/mis')
# =>
#  ssh user-2@192.168.1.1 -p 22
```

Set that up as a nice `~/.bashrc` function, and then you're good to go:
```bash
function shh() {
  ssh $(cat servers.yml | xo "/.*?($1):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/\$4@\$2 -p \$3?:22/mis")
}

# And then we can use it like,
shh production
# =>
#  ssh user-1@192.168.1.1 -p 1234
```

### Fallback values
You may specify fallback values for matches using `$n?:value`, where `n` is the
index that you want to assign the fallback value to. The fallback value should
be simple and can only contain letters, numbers, dashes and underscores; although,
it may contain other indices, in descending order e.g. `$2?:$1`, not `$1?:$2`.

### Delimiters
You may substitute `/` for any delimiter not found within your pattern or formatter.

### Regular expression features
Please see [Go's regular expression documentation](https://golang.org/pkg/regexp/syntax/)
for additional usage options and features.
