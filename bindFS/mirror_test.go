package mirror

import (
	"testing"

	"bazil.org/fuse"
	_ "bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

func Test(t *testing.T) {

	var ctx context.Context

	readlink_test := []struct {
		f      File
		req    *fuse.ReadlinkRequest
		target string
	}{
		{&File{"original/link"}, *fuse.ReadlinkRequest{}, "test1"},
	}

	for _, testcase := range readlink_test {
		rec, _ := (testcase.f).Readlink(ctx, testcase.req)
		if rec != testcase.target {
			t.Errorf("Bad Readlink: %v, need %v", rec, testcase.target)
		}
	}
}
