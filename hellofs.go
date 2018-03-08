package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
)

func main() {
	//mountpoint := flag.String("mp", "", "a directory to mount the fs on")
	flag.Parse()

	mountpoint := flag.Arg(0)

	if flag.NArg() != 1 {
		fmt.Println("no mountpoint provided")
		os.Exit(2)
	}

	con, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		fuse.Subtype("hwfs"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	err = fs.Serve(con, FS{})
	if err != nil {
		log.Fatal(err)
	}

	<-con.Ready
	if err := con.MountError; err != nil {
		log.Fatal(err)
	}

}

// Hello World FUSE startup

type FS struct {
	root *Dir
}

func (FS) Root() (fs.Node, error) {
	return &Dir{}, nil
}

// Root Directory Handler and Node interface implemantation

type Dir struct {
	files int //our helloworld.txt
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	d.files = 1
	return nil
}

//reading the directory, read the txt file name

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hello" {
		return &File{}, nil
	}
	return nil, fuse.ENOENT
}

//child of our root directory, basically reading a dir

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	//var rootchild []fuse.Dirent
	//if d.files == 1 {
	var rootchild = []fuse.Dirent{{Inode: 2, Name: "hello", Type: fuse.DT_File}}
	//}

	return rootchild, nil

}

//File handler

const text = "Hello, World!\n"

type File struct{}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(text))
	return nil
}

//reading the file
func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(text), nil
}
