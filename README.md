# markup

`markdown` is a collection of tools that convert specially formatted text into
HTML. It's designed to be easy to read when it's viewed in its plantext form
and not only when rendered.

At my company markdown is lingua franca for anything that isn't code. Our eng
docs, company handbook, interview guides, etc are all markdown. This is great
for editing but a bit of a pain to actually read.

`markup` is just a simple server that renders markdown.

## Installing markup

```
go get github.com/falun/markup
go install github.com/falun/markup
```

## Running markup

There isn't much in the way of config options:

    $ markup -h
    Usage of markup:
      -serve.index string
          the file returned if '/' is requested; resolved relative to server.root (default README.md)
      -serve.ip string
          the interface we should be listening on (default "localhost")
      -serve.port int
          the port markup will listen on (default 8080)
      -serve.root string
          the root serving directory (default current directory)

All requests are resolved relative to `-serve.root`. Relative links in your
docs that need to descend below root will not resolve correctly.

Typical usage for me looks lomething like `markup -serve.index README.md` which
will end up dropped into `handbook.sh`, `engdocs.sh`, and the like.

## Caching!

Every time you hit a page it reads it from disk  and renders it anew.

## Other Features!

Nope.
