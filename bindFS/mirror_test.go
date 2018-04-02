package mirror

import (
	"testing"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

var ctx context.Context

func TestReadLink(t *testing.T) {

	readlink_test := []struct {
		f      File
		req    *fuse.ReadlinkRequest
		target string
	}{
		{File{"original/file2"}, &fuse.ReadlinkRequest{}, "file1"},
	}

	for _, testcase := range readlink_test {
		rec, _ := (testcase.f).Readlink(ctx, testcase.req)
		if rec != testcase.target {
			t.Errorf("Bad Readlink: %v, need %v", rec, testcase.target)
		}
	}

}

func TestReadDir(t *testing.T) {

	readdirall_test := []struct {
		d     Dir
		child []fuse.Dirent
	}{
		{Dir{"original"}, []fuse.Dirent{{Name: "dir1", Type: fuse.DT_Dir}, {Name: "file.png", Type: fuse.DT_File},
			{Name: "file1", Type: fuse.DT_File}, {Name: "file2", Type: fuse.DT_Link}}},
		{Dir{"original/dir1"}, []fuse.Dirent{{Name: ".rc", Type: fuse.DT_File}, {Name: "a.py", Type: fuse.DT_File},
			{Name: "dir2", Type: fuse.DT_Dir}}},
		{Dir{"original/dir1/dir2"}, []fuse.Dirent{}},
	}

	for _, testcase := range readdirall_test {
		rec, _ := (testcase.d).ReadDirAll(ctx)
		if len(rec) != len(testcase.child) {
			t.Errorf("Bad ReadDirAll: %v, need length %v", len(rec), len(testcase.child))
		}
		for _, dirents := range rec {
			flag := false
			for _, ch := range testcase.child {
				if dirents.Name == ch.Name && dirents.Type == ch.Type {
					flag = true
				}
			}
			if !flag {
				t.Errorf("Bad ReadDirAll: %v, need %v", rec, testcase.child)
			}
		}
	}
}

func TestLookup(t *testing.T) {

	lookup_test := []struct {
		d    Dir
		name string
		node fs.Node
	}{
		{Dir{"original"}, "file2", File{"original/file2"}},
		{Dir{"original"}, "file.png", File{"original/file.png"}},
		{Dir{"original"}, "dir1", Dir{"original/dir1"}},
		{Dir{"original"}, "random", nil},
		{Dir{"original/dir1"}, "dir2", Dir{"original/dir1/dir2"}},
		{Dir{"original/dir1"}, ".rc", File{"original/dir1/.rc"}},
		{Dir{"original/dir1"}, "random", nil},
		{Dir{"original/dir1/dir2"}, "random", nil},
	}

	for _, testcase := range lookup_test {
		if rec, _ := (testcase.d).Lookup(ctx, testcase.name); testcase.node != rec {
			t.Errorf("Bad Lookup: %v, need %v", rec, testcase.node)
		}
	}
}
