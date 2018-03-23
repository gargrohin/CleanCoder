package foo

import (
	_ "flag"
	_ "fmt"
	_ "log"
	"os"

	"bazil.org/fuse"
	_ "bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

func mount(mountpoint string) error {
	//mountpoint := flag.String("mp", "", "a directory to mount the fs on")

	_, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		fuse.Subtype("hwfs"),
	)

	return err

	//defer con.Close()

	/*err = fs.Serve(con, FS{})
	if err != nil {
		log.Fatal(err)
	}

	<-con.Ready
	if err := con.MountError; err != nil {
		log.Fatal(err)
	}
	return err*/
}

// Hello World FUSE startup

type FS struct {
	root *Dir
}

func (FS) Root() error {
	return nil
}

// Root Directory Handler and Node interface implemantation

type Dir struct {
	Files int //our helloworld.txt
}

type Atr struct {
	Inode int
	Mode  os.FileMode
	Size  uint64
}

var a Atr

func (d *Dir) Attr(ctx context.Context) (*Dir, Atr, error) {
	//var a *fuse.Attr
	a.Inode = 0
	a.Mode = os.ModeDir | 0777
	d.Files = 2
	return d, a, nil
}

//reading the directory, read the txt file name
func (d *Dir) Lookup(ctx context.Context, name string) error {
	if name == "" {
		return nil
	}
	/*if name == "inode" {
		return &File2{}, nil
	}*/
	return fuse.ENOENT
}

//child of our root directory, basically reading a dir

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	//var rootchild []fuse.Dirent
	//if d.files == 1 {
	var rootchild = []fuse.Dirent{{Inode: 1, Name: "", Type: fuse.DT_File}, {Inode: 2, Name: "inode", Type: fuse.DT_File}} // {Inode: 3, Type: fuse.DT_File, Name: "three"}}
	//}

	return rootchild, nil

}

//File handlers

const text = "Hello, World!\n"

type File struct{}

func (f *File) Attr(ctx context.Context) (Atr, error) {
	a.Inode = 1
	a.Mode = 0777
	a.Size = uint64(len(text))
	return a, nil
}

//reading the file
func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(text), nil
}

// File 2 to see how inodes work

/*const txt2 = "File 2. Check Inode working"

type File2 struct{}

func (f *File2) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 5
	a.Mode = 0777
	a.Size = uint64(len(txt2))
	return nil
}

func (f *File2) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(txt2), nil
}*/
