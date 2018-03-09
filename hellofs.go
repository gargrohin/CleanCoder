package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
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

	//defer con.Close()

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
	a.Inode = 0
	a.Mode = os.ModeDir | 0777
	d.files = 2
	return nil
}

//reading the directory, read the txt file name

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hw" {
		return &File{}, nil
	}
	if name == "inode" {
		return &File2{}, nil
	}
	return nil, fuse.ENOENT
}

//child of our root directory, basically reading a dir

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	//var rootchild []fuse.Dirent
	//if d.files == 1 {
	var rootchild = []fuse.Dirent{{Inode: 1, Name: "hw", Type: fuse.DT_File}, {Inode: 2, Name: "inode", Type: fuse.DT_File}, {Inode: 3, Type: fuse.DT_File, Name: "three"}}
	//}

	return rootchild, nil

}

//File handlers

const text = "Hello, World!\n"

type File struct{}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = 0777
	a.Size = uint64(len(text))
	return nil
}

//reading the file
func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(text), nil
}

// File 2 to see how inodes work

const txt2 = "File 2. Check Inode working"

type File2 struct{}

func (f *File2) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 5
	a.Mode = 0777
	a.Size = uint64(len(txt2))
	return nil
}

func (f *File2) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(txt2), nil
}
