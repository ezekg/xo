# `xo`

`xo` is a command line utility similar to `sed`, except it has one job: take an
input string from `stdin` and format the matches. That's it. I built this so
that I could read configuration files and use its contents to string
together a new command.

## Usage
Suppose we had a text file called `starwars.txt` containing some Star Wars quotes,
```
Vader: If only you knew the power of the Dark Side. Obi-Wan never told you what happened to your father.
Luke: He told me enough! He told me you killed him!
Vader: No, I am your father.
Luke: [shocked] No. No! That's not true! That's impossible!
```

and we wanted to do a little formatting, as if we're telling it as a story. Easy!
```bash
cat starwars.txt | xo 'mi/^(\w+):(\s*\[(.*?)\]\s*)?\s*([^\n]+)/$1 said, "$4" in a $3?:normal voice./'
# =>
#   Vader said "If only you knew the power of the Dark Side. Obi-Wan never told you what happened to your father." in a normal voice
#   Luke said "He told me enough! He told me you killed him!" in a normal voice
#   Vader said "No, I am your father." in a normal voice
#   Luke said "No. No! That's not true! That's impossible!" in a shocked voice
```

Okay, okay. Lets do something a little more useful. Now suppose we had a configuration
file called `servers.yml` containing some project information. Maybe it looks like
this,
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
cat servers.yml | xo 'mis/.*?(production):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/$4@$2 -p $3?:22/'
# =>
#  user-1@192.168.1.1 -p 1234

# Now let's actually use the output,
ssh $(cat servers.yml | xo 'mis/.*?(staging):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/$4@$2 -p $3?:22/')
# =>
#  ssh user-2@192.168.1.1 -p 22
```

Set that up as a nice `.bashrc` function, and then you're good to go:
```bash
function gosh() {
  ssh $(cat servers.yml | xo "mis/.*?($1):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/\$4@\$2 -p \$3?:22/")
}

# And then we can use it like,
gosh production
# =>
#  ssh user-1@192.168.1.1 -p 1234
```

## Usage
`xo` accepts the following syntax; all you have to do is feed it some `stdin` via
piped output (`echo 'hello' | xo ...`) or what have you. There's no flags, and no
additional arguments. Simple and easy to use.
```
xo 'modifiers/pattern/formatter/'
xo '/pattern/formatter/'
```

### Fallback values
You may specify fallback values for match indexes using `$N?:value`, where `N`
is the index that you want to assign the fallback value to. The fallback value
should be simple and cannot contain a space or newline.

### Delimiters
You may substitute `/` for any delimiter not found unescaped within your pattern
or formatter.

### Regular expression features
Please see [Go's regular expression documentation](https://golang.org/pkg/regexp/syntax/)
for additional usage options and features.
