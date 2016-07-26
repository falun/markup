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

There isn't much in the way of config options:

    $ markup -h
		Usage of markup:
		  -index.exclude string
		        these directories (comma separated) will not be included in the
            generated /index. They are resolved relative to index.root
		  -index.root string
		        the root serving directory (default current directory)
		  -serve.home string
		        the file returned if '/' is requested; resolved relative to server.root (default "README.md")
		  -serve.ip string
		        the interface we should be listening on (default "localhost")
		  -serve.port int
		        the port markup will listen on (default 8080)

All requests are resolved relative to `-index.root`. Relative links in your
docs that need to descend below root will not resolve correctly.

Typical usage for me looks lomething like `markup -serve.index README.md` which
will end up dropped into `handbook.sh`, `engdocs.sh`, and the like.

## Caching!

Every time you hit a page it reads it from disk  and renders it anew.

## Other Features!

`/index` will serve a list of all markdown files found by descending from
`serve.root`. This file list is built once at startup, follows symlinks,
doesn't enforce any kind of max depth, and blocks serving until the walk
finishes (so there are some improvements that can be made).
