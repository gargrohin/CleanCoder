package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	_ "io/ioutil"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s ROOT MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	//mountpoint := flag.String("mp", "", "a directory to mount the fs on")
	flag.Parse()

	mountpoint := flag.Arg(0)

	if flag.NArg() != 1 {
		usage()
		fmt.Println("no mountpoint provided")
		os.Exit(2)
	}

	hw := &File{"hw", "hello world!\n", 1, 0777}
	//inode := &File{"inode", "inode works how?\n", 2, 0777}
	root := &Dir{"test", 2, 3, nil, hw, nil}
	myfs := &FS{"hellofs", root}

	con, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		fuse.Subtype("hwfs"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	err = fs.Serve(con, myfs)
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
	name string
	root *Dir
}

func (f FS) Root() (fs.Node, error) {
	return f.root, nil
}

// Root Directory Handler and Node interface implemantation

type Dir struct {
	name      string
	files     int //our helloworld.txt
	inode     uint64
	nextdir   *Dir
	nextfile  *File
	nextfile2 *File
}

func (d Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = d.inode
	a.Mode = os.ModeDir | 0777
	//d.files = 2
	return nil
}

//reading the directory, read the txt file name

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if d.nextfile != nil && name == d.nextfile.name {
		return *(d.nextfile), nil
	}
	if d.nextdir != nil && name == d.nextdir.name {
		return *(d.nextdir), nil
	}
	if d.nextfile2 != nil {
		return *(d.nextfile2), nil
	}
	return nil, fuse.ENOENT
}

//children of our root directory, basically reading a dir

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	//var rootchild []fuse.Dirent
	//if d.files == 1 {
	var rootchild = []fuse.Dirent{}
	if d.nextdir != nil {
		rootchild = append(rootchild, fuse.Dirent{Inode: d.nextdir.inode, Name: d.nextdir.name, Type: fuse.DT_Dir})
	}
	if d.nextfile != nil {
		rootchild = append(rootchild, fuse.Dirent{Inode: d.nextfile.inode, Name: d.nextfile.name, Type: fuse.DT_File})
	}
	if d.nextfile2 != nil {
		rootchild = append(rootchild, fuse.Dirent{Inode: d.nextfile2.inode, Name: d.nextfile2.name, Type: fuse.DT_File})
	}

	return rootchild, nil

}

func (d Dir) MkDir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	dir := &Dir{name: req.Name, files: 0, inode: 10, nextdir: nil, nextfile: nil, nextfile2: nil}
	d.nextdir = dir
	return dir, nil
}

func (d Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	f := &File{name: req.Name, text: "woah\n", inode: uint64(resp.LookupResponse.Node), osmode: req.Mode}
	//f, err := os.Create(req.Name)

	//b, _ := ioutil.ReadFile(f.Name())
	d.nextfile2 = f
	return d.nextfile2, d.nextfile2, nil
}

//File handlers

//const text = "Hello, World!\n"

type File struct {
	name   string
	text   string
	inode  uint64
	osmode os.FileMode
}

func (d Dir) Write(ctx context.Context, req fuse.WriteRequest, resp fuse.WriteResponse) error {
	resp.Size = (len(d.nextfile.text) + len(req.Data))
	d.nextfile.text = d.nextfile.text + string(req.Data)
	return nil
}

/*func (f File) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {

}*/

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = f.inode
	a.Mode = f.osmode
	a.Size = uint64(len(f.text))
	return nil
}

//reading the file
func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(f.text), nil
}

func (f File) Write(ctx context.Context, req fuse.WriteRequest, resp fuse.WriteResponse) error {
	resp.Size = (len(f.text) + len(req.Data))
	f.text = f.text + string(req.Data)
	return nil
}

/*func (f File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	return nil
}*/
