# xo

`xo` is a command line utility similar to `sed`, except it has one job: take an
input string from `stdin` and format its matches. That's it. I built this so
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
cat starwars.txt | xo 'mis/^(\w+):(.*?\[(.*?)\].*?)?([^\n]+)/$1 said "$4" in a $3?:normal voice/'
```

## Syntax
`xo` accepts the following syntax,
```
xo 'modifiers/pattern/formatter/'
xo '/pattern/formatter/'
```

## Fallback values
You may specify fallback values for match indexes using `$N?:value`, where `N`
is the index that you want to assign the fallback value to. The fallback value
should be simple and cannot contain a space or newline.

## Delimiters
You may substitute `/` for any delimiter not found unescaped within your pattern
or formatter. Please see [Go's regular expression documentation](https://golang.org/pkg/regexp/syntax/)
for additional usage options.
