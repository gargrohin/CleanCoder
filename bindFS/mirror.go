package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

func main() {
	flag.Parse()
	mountpoint := flag.Arg(0)

	root := &Dir{flag.Arg(1)}
	myFS := &FS{flag.Arg(1), root}
	con, err := fuse.Mount(
		mountpoint, fuse.FSName("FLCL-mirror"),
		fuse.Subtype("bindfs"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	err = fs.Serve(con, myFS)
	if err != nil {
		log.Fatal(err)
	}

	<-con.Ready
	if err := con.MountError; err != nil {
		log.Fatal(err)
	}

}

type FS struct {
	path string
	root *Dir
}

func (f FS) Root() (fs.Node, error) {
	return f.root, nil
}

// Root Directory Handler and Node interface implemantation

type Dir struct {
	name string
}

func (d Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	//dir,_:=os.Open(d.name)
	var datr syscall.Stat_t
	syscall.Stat(d.name, &datr)
	a.Inode = datr.Ino
	m := os.FileMode(datr.Mode)
	a.Mode = os.ModeDir | m
	return nil
}

//reading the directory, read the txt file name

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	dir, _ := os.Open(d.name)
	fileinfo, _ := dir.Readdir(-1)
	for _, child := range fileinfo {
		if name == child.Name() {
			if child.IsDir() {
				return Dir{d.name + "/" + child.Name()}, nil
			} else {
				return File{d.name + "/" + name}, nil
			}
		}
	}

	return nil, fuse.ENOENT
}

//children of our root directory, basically reading a dir

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	file, _ := os.Open(d.name)

	var rootchild = []fuse.Dirent{}

	fileinfo, _ := file.Readdir(-1)
	for _, child := range fileinfo {
		datr := syscall.Stat_t{}
		syscall.Stat(d.name+"/"+child.Name(), &datr)

		if child.IsDir() {
			rootchild = append(rootchild, fuse.Dirent{Inode: datr.Ino, Name: child.Name(), Type: fuse.DT_Dir})
		} else {
			rootchild = append(rootchild, fuse.Dirent{Inode: datr.Ino, Name: child.Name(), Type: fuse.DT_File})
		}
	}

	return rootchild, nil

}

type File struct {
	name string
}

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {

	var fatr syscall.Stat_t
	syscall.Stat(f.name, &fatr)
	a.Inode = fatr.Ino
	a.Mode = os.FileMode(fatr.Mode)
	a.Size = uint64(fatr.Size)
	return nil
}

//reading the file
func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	//file,_=os.Open(f.name)
	txt, err := ioutil.ReadFile(f.name)
	if err != nil {
		log.Fatal(err)
	}
	return txt, nil
}
