// +build ignore

// called by go generate in subdirectory cmd/somaadm

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gholt/blackfridaytext"
	"github.com/gholt/brimtext"
)

func main() {
	align := brimtext.NewSimpleAlignOptions()
	opt := &blackfridaytext.Options{
		Color:             true,
		TableAlignOptions: align,
	}
	fs, _ := ioutil.ReadDir(os.Args[1])
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), `.md`) {
			inFile := filepath.Join(os.Args[1], f.Name())
			outFile := filepath.Join(os.Args[2],
				strings.TrimSuffix(f.Name(), `.md`)+`.fmt`)
			markdown, _ := ioutil.ReadFile(inFile)
			_, rndrOut := blackfridaytext.MarkdownToText(markdown, opt)
			ioutil.WriteFile(outFile, rndrOut, os.FileMode(0640))
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
