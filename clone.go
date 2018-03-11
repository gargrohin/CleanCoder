func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hw" {
		return &File{}, nil
	}
	/*if name == "inode" {
		return &File2{}, nil
	}*/
	return nil, fuse.ENOENT
}

//child of our root directory, basically reading a dir

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	//var rootchild []fuse.Dirent
	//if d.files == 1 {
	var rootchild = []fuse.Dirent{{Inode: 1, Name: "hw", Type: fuse.DT_File}, {Inode: 2, Name: "inode", Type: fuse.DT_File}} // {Inode: 3, Type: fuse.DT_File, Name: "three"}}
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
