package adm

import (
	"io"
	"os"
)

func CopyFile(dst, src string) (written int64, err error) {
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer out.Close()

	in, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		return
	}

	return io.Copy(out, in)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
