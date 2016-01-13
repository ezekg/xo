# xo

`xo` is a command line utility similar to `sed`, except it has one job: take an
input string from `stdin` and format its matches. That's it. I built this so
that I could read configuration files and use its contents to string
together a new command.

Suppose we had a configuration file similar to this for every project,
```yml
# ...
stages:
  production:
    server: 192.168.1.1:1234
    user: user-1
  staging:
    server: 192.168.1.1
    user: user-2
# ...
```

and we wanted to SSH into it without having to remember the credentials. Easy!
```bash
ssh $(cat config.yml | xo "mis/.*?($1):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/\$4@\$2 -p \$3/")
```

Set that up as a nice `.bashrc` function, and then you're good to go:
```bash
function gosh() {
  ssh $(cat config.yml | xo "mis/.*?($1):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/\$4@\$2 -p \$3/")
}
```

`xo` accepts the following,
```
xo 'modifiers/pattern/formatter/'
xo '/pattern/formatter/'
```

You may substitute `/` for any delimiter not found unescaped within your pattern
or formatter. Please see [Go's regular expression documentation](https://golang.org/pkg/regexp/syntax/)
for additional usage options.
