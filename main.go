package main

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/dave/console"
	"github.com/dave/dropper"
	"github.com/dave/saver"
	"honnef.co/go/js/dom"
)

func main() {
	go run()
}

func run() {
	w := &console.Writer{}
	w.Message("Drag files here to zip")

	// initialise the drag+drop function with the github.com/dave/dropper package
	events := dropper.Initialise(dom.GetWindow().Document())

	// the dropper package creates a channel of events
	for ev := range events {
		switch ev := ev.(type) {
		case dropper.DropEvent:
			// accept a dropped file
			w.Message("Processing")

			// choose a filename
			var name string
			if len(ev) == 1 {
				name = strings.TrimSuffix(ev[0].Name(), filepath.Ext(ev[0].Name())) + ".zip"
			} else {
				name = "files.zip"
			}

			// zip the files using the standard library zip package
			buf := &bytes.Buffer{}
			zw := zip.NewWriter(buf)
			for _, f := range ev {
				w, err := zw.Create(f.Path())
				if err != nil {
					panic(err)
				}
				if _, err := io.Copy(w, f.Reader()); err != nil {
					panic(err)
				}
			}
			zw.Close()

			// save the file as a browser download using the github.com/dave/saver package
			saver.Save(name, "application/zip", buf.Bytes())
			w.Message("Done")

		case dropper.EnterEvent:
			// drag event enters the page
			w.Message("Drop here")

		case dropper.LeaveEvent:
			// drag event exits the page
			w.Message("Drag files here to zip")

		}
	}
}
