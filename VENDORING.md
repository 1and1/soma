# VENDORING

Before using `godep` to update vendored sources, remove
the build ignore flag from `script/render_markdown.go`.

Otherwise the markdown processor used to generate the
somaadm help text will not be tracked.
