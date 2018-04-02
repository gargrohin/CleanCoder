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
	"time"
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
	defer file.Close()

	var rootchild = []fuse.Dirent{}

	fileinfo, _ := file.Readdir(-1)
	for _, child := range fileinfo {
		datr := syscall.Stat_t{}
		syscall.Stat(d.name+"/"+child.Name(), &datr)

		file_type := []struct {
			check     os.FileMode
			fuse_type fuse.DirentType
		}{
			{os.ModeDir, fuse.DT_Dir},
			{os.ModeSymlink, fuse.DT_Link},
			{os.ModeSocket, fuse.DT_Socket},
			{os.ModeCharDevice, fuse.DT_Char},
			{os.ModeNamedPipe, fuse.DT_FIFO},
			{os.FileMode(0xffffffff), fuse.DT_File},
		}

		for _, ftype := range file_type {
			if child.Mode()&ftype.check != 0 {
				rootchild = append(rootchild, fuse.Dirent{Inode: datr.Ino, Name: child.Name(), Type: ftype.fuse_type})
				break
			}
		}
	}

	return rootchild, nil

}

//creates a dir, i.e implementation of mkdir command

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	err := os.Mkdir(d.name+"/"+req.Name, req.Mode)
	dchild := &Dir{name: d.name + "/" + req.Name}
	return dchild, err
}

//creates a file, i.e. implementation of touch command

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	file_opened, err := os.Create(d.name + "/" + req.Name)
	f := &File{name: d.name + "/" + req.Name}
	fh := &Handle{name: d.name + "/" + req.Name, handle: file_opened}
	return f, fh, err
}

func (d Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	err := os.Remove(d.name + "/" + req.Name)
	return err
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
	fileinfo, _ := os.Lstat(f.name)
	a.Inode = fatr.Ino
	a.Mode = os.FileMode(fatr.Mode)
	a.Size = uint64(fatr.Size)
	a.Mtime = time.Time(fileinfo.ModTime())
	a.Blocks = uint64(fatr.Blocks)
	a.BlockSize = uint32(fatr.Blksize)
	a.Nlink = uint32(fatr.Nlink)
	a.Uid = uint32(fatr.Uid)
	a.Gid = uint32(fatr.Gid)
	a.Rdev = uint32(fatr.Rdev)
	return nil
}

func (f File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	if req.Valid.Size() {
		os.Truncate(f.name, int64(req.Size))
	}
	if req.Valid.Atime() || req.Valid.Mtime() {
		os.Chtimes(f.name, req.Atime, req.Mtime)
	}
	if req.Valid.Gid() || req.Valid.Uid() {
		os.Chown(f.name, int(req.Uid), int(req.Gid))
	}
	if req.Valid.Mode() {
		os.Chmod(f.name, req.Mode)
	}

	return nil
}

//reading the file

func (fh Handle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	txt, err := ioutil.ReadFile(fh.name)
	resp.Data = txt
	return err
}

func (fh Handle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	n, err := fh.handle.Write(req.Data)
	//log.Println(req.Offset)
	resp.Size = int(n)
	//newsize := req.Offset + int64(n)
	//err = os.Truncate(fh.name, newsize)
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

func (f File) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
	target, err := os.Readlink(f.name)
	//log.Println(target)
	return target, err
}
