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
      -default-file string
            the file returned if '/' is requested; resolved relative to root (default "README.md")
      -ip string
            the interface we should be listening on (default "localhost")
      -port int
            the port markup will listen on (default 8080)
      -root string
            the root serving directory (default is the current directory)

In the most recent push sever flags were changed. The following flags are
_deprecated_ but still function:

old flag      | equivalent new flag
--------------|--------------------
index.root    | root
serve.default | default-file
serve.ip      | ip
serve.port    | port

Additionally the flag `index.exclude` no longer functions due to changing the
way the `/index` endpoint works and will be removed in the future. It is retained
for now for backwards compatibility with scripts that have been written to
start `markup`.

All requests are resolved relative to `root`. Relative links in your
docs that need to descend below root will not resolve correctly.

Typical usage for me looks lomething like `markup -default SOMEFILE.md` which
will end up dropped into `handbook.sh`, `engdocs.sh`, and the like.

## Caching!

Every time you hit a page it reads it from disk  and renders it anew.

## Other Features!

`/index` serves a directory listing starting at `-root`. Each markdown file
will be rendered with a link to its parent directory.
