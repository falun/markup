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

The most common case is to run markup from the root directory of a project that
you want to read docs in. The generated config will be stored in
`./markup/conf.json` and then used for subsequent executions.

To customize behavior there is a small but moving set of configuration options.
The core set of flags are:

| Flag                 | Meaning | Default |
| -------------------- | ------- | ------- |
| config-dir           | directory housing the config file | `<pwd>/.markup` |
| config-file          | name of the file housing config; may be overwritten by passing a file as argument 1 | `<config-dir>/conf.json` |
| default-file         | the file that will be rendered initially when accessing `/` | README.md |
| port                 | what port should the server listen on for requests | 8080 |
| root                 | the directory that will be used as a starting point to locate all requested files | `<pwd>` |
| dc                   | prints out the configuration that will be used if markup is run with the provided flags | |
| force-install-assets | will overwrite the contents of `<asset-dir>` with the default set of web pages | false |

A full set of config flags cmay be found by `markup -h`.

As suggested in the table markup's config is found by checking in `<config-dir>/<config-file>`
but may be explicitly specified (overriding both `config-X` options) by passing the
file on the command line: `markup ./custom-config.json`.

## Caching?

None!

Every time you hit a page it reads it from disk and renders it anew.

Every time you view a page it looks for a template on disk and parses it if found.

## Other Features!

`/index` serves a directory listing starting at `-root`. Each markdown file
will be rendered with a link to its parent directory.

`/asset` will load a file from the asset dir without applying any transformations.
This may be used for supplemental resources (see `assets/browse.tmpl`'s use of
`/asset/style.css`).

`/i/search` can be used to find files based on file name matching. It has the
following query parameters:

| Parameter        | Meaning |
| ---------------- | ------- |
| `term`           | The file being searched for. Spaces are replaced with a wildcard, regexp may be used. |
| `ft`             | A comma-separated list of file extenstions. Leading `.` is not required and case is ignored. |
| `case-sensitive` | If passed `term` is matched in a case sensitive manner. |
| `include-dir`    | If passed the full path (vs file name only) will be searched for a match to `term`. |
| `page`           | Indicates which page of results to return; 0-indexed. |

Example:

```shell
$ curl -s 'localhost:8080/i/search?term=(browse|readme)&ft=md,go'
{
  "search_term": "(browse|readme)",
  "file_types": [
    ".md",
    ".go"
  ],
  "last_scan_ms": 1487234369950,
  "match": [
    "app\\browse.go",
    "README.md",
    "vendor\\github.com\\russross\\blackfriday\\README.md",
    "vendor\\github.com\\shurcooL\\sanitized_anchor_name\\README.md",
    "web\\browse.go"
  ],
  "count": 5
}
```

Directories may be excluded from this search by specifying `exclude-dirs` at
startup and the frequency with which the index will be rebuilt is controled by
`scan-freq`. For more information see `markup -h`.
