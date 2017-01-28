# markup

`markdown` is a style of formatting plain text and a collection of tools that
convert this format into HTML. It's designed to be easy to write and readable
in its plantext form.

At my company markdown is lingua franca for anything that isn't code. Our eng
docs, company handbook, interview guides, etc are all markdown. This is great
for editing but a bit of a pain to actually read in a rendered form.

`markup` is just a simple server that renders markdown on the fly.

## Installing markup

```
go get github.com/falun/markup
go install github.com/falun/markup
```

## Running markup

There is a small but moving set of configuration opts options for `markup`. The
core set necessary to meaningfully use it are:

| Flag         | Meaning | Default |
| ------------ | ------- | ------- |
| config-dir   | directory housing the config file | `<pwd>/.markup` |
| config-file  | name of the file housing config; may be overwritten by passing a file as argument 1 | `<config-dir>/conf.json` |
| default-file | the file that will be rendered initially when accessing `/` | README.md |
| port         | what port should the server listen on for requests | 8080 |
| root         | the directory that will be used as a starting point to locate all requested files | `<pwd>` |
| dc           | prints out the configuration that will be used if markup is run with the provided flags | |

The full set of config flags cmay be found by `markup -h`.

As indicated above markup config is found by checking in `config-dir` + `config-file`
but may be explicitly specified (overriding both `config-X` options) by passing the
file on the command line: 'markup ./custom-config.json`.

## Caching!

Every time you hit a page it reads it from disk  and renders it anew.

## Other Features!

`/index` serves a directory listing starting at `-root`. Each markdown file
will be rendered with a link to its parent directory.
