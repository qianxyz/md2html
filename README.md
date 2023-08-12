# md2html

A simple script for converting Markdown files to HTML (powered by [GitHub
Markdown
API](https://docs.github.com/en/rest/markdown/markdown?apiVersion=2022-11-28))
with live preview.

## Quickstart

```shell
$ go install github.com/qianxyz/md2html@latest
$ md2html
Usage: md2html [-p port] file.md
```

The preview is at `http://0.0.0.0:8080` or on the port specified. After editing
and saving the Markdown file, refresh the page to see the changes go live.
