# `xo`
[![Travis](https://img.shields.io/travis/ezekg/xo.svg?style=flat-square)](https://travis-ci.org/ezekg/xo)

`xo` is a command line utility that composes regular expression match groups.
What differentiates `xo` from tools like `sed` and `awk` is that `xo` is designed
to do a single job well, and that is compose together match groups into a new
string. `xo` is not meant to handle things that `sed` and `awk` are already
good at; namely, search and replace. It simply performs one of the few things
those tools cannot accomplish intuitively. For example,

```bash
echo 'Hello! My name is C3PO, human cyborg relations.' | xo '/^(\w+)! my name is (\w+)/$1, $2!/i'
# =>
#  Hello, C3PO!
```

You may find yourself using `xo` to format logs into something a bit more human-readable,
compose together command output into a new command, or even normalize some data using
[fallback values](#fallback-values). `xo` also comes with the full power of
[Go's regular expression syntax](https://golang.org/pkg/regexp/syntax/); meaning
it can handle multi-line patterns (something `sed` doesn't do very well!), as well
as any other flag you want to throw at it.

## Installation
To install `xo`, please use `go get`. If you don't have Go installed, [get it here](https://golang.org/dl/).
If you would like to grab a precompiled binary, head over to the [releases](https://github.com/ezekg/xo/releases)
page. The precompiled `xo` binaries have no external dependencies.

```
go get github.com/ezekg/xo
```

## Usage
`xo` accepts the following syntax; all you have to do is feed it some `stdin` via
piped output (`echo 'hello' | xo ...`) or what have you. There's no command line
flags, and no additional arguments. Simple and easy to use.
```
xo '/<pattern>/<formatter>/[flags]'
```

## Examples
Let's start off a little simple, and then we'll ramp it up and get crazy. `xo`,
in its simplest form, does things like this,
```bash
echo 'Hello! My name is C3PO, human cyborg relations.' | xo '/^(\w+)?! my name is (\w+)/$1, $2!/i'
# =>
#  Hello, C3PO!
```

Here's a quick breakdown of what each piece of the puzzle is,
```bash
echo 'Hello! My name is C3PO.' | xo '/^(\w+)?! my name is (\w+)/$1, $2!/i'
^                              ^     ^^                       ^ ^     ^ ^
|______________________________|     ||_______________________| |_____| |
                |                    + Delimiter |                 |    + Flag
                + Piped output                   + Pattern         + Formatter
```

When you create a regular expression, wrapping a subexpression in parenthesis `(...)`
creates a new _capturing group_, numbered from left to right in order of opening
parenthesis. Submatch `$0` is the match of the entire expression, submatch `$1`
the match of the first parenthesized subexpression, and so on. These capturing
groups are what `xo` works with.

What about the question mark? The question mark makes the preceding token in the
regular expression optional. `colou?r` matches both `colour` and `color`. You can
make several tokens optional by _grouping_ them together using parentheses, and
placing the question mark after the closing parenthesis, e.g. `Nov(ember)?`
matches `Nov` and `November`.

With that, what if the input string _forgot_ to specify a greeting, but we, desiring
to be polite, still wanted to say "Hello"? Well, that sounds like a great job for
a [fallback value](#fallback-values)! Let's update the example a little bit,
```bash
echo 'Hello! My name is C3PO.' | xo '/^(?:(\w+)! )?my name is (\w+)/$1?:Greetings, $2!/i'
# =>
#  Hello, C3PO!

echo 'My name is Chewbacca, uuuuuur ahhhhhrrr uhrrr ahhhrrr aaargh.' | xo '/^(?:(\w+)! )?my name is (\w+)/$1?:Greetings, $2!/i'
# =>
#  Greetings, Chewbacca!
```

As you can see, we've taken the matches and created a new string out of them. We
also supplied a [fallback value](#fallback-values) for the second match (`$2`)
that gets used if no match is found, using the ternary `?:` operator.

Now that we have the basics of `xo` out of the way, let's pick up the pace a little
bit. Suppose we had a text file called `starwars.txt` containing some Star Wars quotes,
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

Okay, okay. Let's move away from Star Wars references and on to something a little
more useful. Suppose we had a configuration file called `servers.yml` containing
some project information. Maybe it looks like this,
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

Lastly, what about reading sensitive credentials from an ignored configuration
file to pass to a process, say, `rails s`? Let's use Stripe keys as an example
of something we might not want to log to our terminal history,
```bash
cat secrets/stripe.yml | xo '/test_secret_key:\s([\w]+).*?test_publishable_key:\s([\w]+)/PUBLISHABLE_KEY=$1 SECRET_KEY=$2 rails s/mis' | sh
```

Pretty cool, huh?

## Fallback values
You may specify fallback values for matches using the ternary operator, `$i?:value`,
where `i` is the index that you want to assign the fallback value to. The fallback
value may contain any sequence of characters, though anything other than letters,
numbers, dashes and underscores must be escaped; it may also contain other match
group indices if they are in descending order e.g. `$2?:$1`, not `$1?:$2`.

## Delimiters
You may substitute `/` for any delimiter. If the delimiter is found within your pattern
or formatter, it must be escaped. If it would normally be escaped in your pattern
or formatter, it must be escaped again. For example,

```bash
# Using the delimiter `|`,
echo 'Hello! My name is C3PO, human cyborg relations.' | xo '|^(\w+)?! my name is (\w+)|$1, $2!|i'

# Using the delimiter `w`,
echo 'Hello! My name is C3PO, human cyborg relations.' | xo 'w^(\\w+)?! my name is (\\w+)w$1, $2!wi'
```

## Regular expression features
Please see [Go's regular expression documentation](https://golang.org/pkg/regexp/syntax/)
for additional usage options and features.
