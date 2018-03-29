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

//creates a dir

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	err := os.Mkdir(d.name+"/"+req.Name, req.Mode)
	dchild := &Dir{name: d.name + "/" + req.Name}
	return dchild, err
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	_, err := os.Create(d.name + "/" + req.Name)
	f := &File{name: d.name + "/" + req.Name}
	return f, f, err
}

type File struct {
	name string
}

type Handle struct {
	name   string
	handle *os.File
}

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {

	var fatr syscall.Stat_t
	syscall.Lstat(f.name, &fatr)
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

func (fh Handle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	txt, err := ioutil.ReadFile(fh.name)
	resp.Data = txt
	return err
}

func (fh Handle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	n, err := fh.handle.Write(req.Data)
	resp.Size = int(n)
	return err

}

func (fh Handle) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	fh.handle.Sync()
	return nil
}

func (f File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	fileinfo, err := os.Lstat(f.name)

	fh, err := os.OpenFile(f.name, int(req.Flags), fileinfo.Mode())

	resp.Handle = fuse.HandleID(req.Header.Node)
	resp.Flags = fuse.OpenResponseFlags(req.Flags)

	return &Handle{name: f.name, handle: fh}, err
}

func (fh Handle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	err := fh.handle.Close()
	return err
}
