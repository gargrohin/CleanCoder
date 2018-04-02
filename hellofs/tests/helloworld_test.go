package foo

import (
	"bazil.org/fuse"
	"golang.org/x/net/context"
	"os"
	"testing"
)

var fhandle FS

func TestRoot(t *testing.T) {
	err := fhandle.Root()
	if err != nil {
		t.Error("root returning error")
	}
}

var dir Dir
var ctx context.Context
var f File

func TestDirAttr(t *testing.T) {

	d, a, err := dir.Attr(ctx)

	if err != nil || d.Files != 2 || a.Mode != os.ModeDir|0777 {
		t.Error("error in Directory's attributes")
	}

	a, err = f.Attr(ctx)

	if err != nil || a.Size != uint64(len(text)) || a.Mode != 0777 {
		t.Error("error in File's attributes")
	}

}

func TestReadDir(t *testing.T) {

	rootchild, err := dir.ReadDirAll(ctx)

	if rootchild[0].Type != fuse.DT_File || rootchild[0].Name != "" || err != nil {
		t.Error("error in links between root and children")
	}

	txt, err := f.ReadAll(ctx)

	if err != nil || text != string(txt[:]) {
		t.Error("error in File's contents")
	}

}
