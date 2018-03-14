package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	_ "syscall"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s ROOT MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

var fileinfo os.FileInfo
var filelist []string
var err error
var origin string
var data []string

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		fmt.Println("2 arguments needed")
		os.Exit(2)
	}

	mountpoint := flag.Arg(0)
	origin = flag.Arg(1)

	file, err := os.Open(origin)
	if err != nil {
		log.Fatal(err)
	}

	fileinfo, err = os.Stat(origin)

	filelist, err = file.Readdirnames(-1)
	//fmt.Println(filelist, len(filelist))
	//a := []string{origin, filelist[0]}
	//path := strings.Join(a, "/")

	text, err := ioutil.ReadAll(file)
	data = append(data, string(text[:]))

	con, err := fuse.Mount(
		mountpoint,
		fuse.FSName("clone"),
		fuse.Subtype("cln"),
	)

	if err != nil {
		log.Fatal(err)
	}

	//defer con.Close()

	err = fs.Serve(con, FS{})
	if err != nil {
		log.Fatal(err)
	}

}

type FS struct{}

func (FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

type Dir struct {
	files int //our helloworld.txt
}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 0
	a.Mode = os.ModeDir | 0777
	//d.files = len(filelist)
	return nil
}

//reading the directory, read the txt file name

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == filelist[0] {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

//child of our root directory, basically reading a dir

func (Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	a := []string{origin, filelist[0]}
	fileinfo, _ = os.Stat(strings.Join(a, "/"))
	if fileinfo.IsDir() == true {
		tp := fuse.DT_Dir
		var rootchild = []fuse.Dirent{{Inode: 1, Name: filelist[0], Type: tp}}
		return rootchild, nil
	} else {
		tp := fuse.DT_File
		var rootchild = []fuse.Dirent{{Inode: 1, Name: filelist[0], Type: tp}}
		return rootchild, nil
	}

}

type File struct{}

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = 0777
	a.Size = uint64(fileinfo.Size())
	return nil
}

//reading the file
func (File) ReadAll(ctx context.Context) ([]byte, error) {
	file, err := os.Open(origin)
	if err != nil {
		log.Fatal(err)
	}

	text, err := ioutil.ReadAll(file)
	return text, nil
}
