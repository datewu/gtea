package static

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
)

// FS implt http.FileSystem
type FS struct {
	NoDir   bool
	TryFile []string
	Root    string // http.Dir()
}

func (f FS) Open(name string) (http.File, error) {
	if f.TryFile != nil {
		indexs := []string{"index.html", "index.htm",
			"index.php", "main.jsx", "index", "readme.md"}
		f.TryFile = append(f.TryFile, indexs...)
	}
	dir := http.Dir(f.Root)
	file, err := dir.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) && f.TryFile != nil {
			for _, try := range f.TryFile {
				good, err := dir.Open(try)
				if err == nil {
					return good, nil
				}
			}
		}
		return nil, err
	}
	if f.NoDir {
		info, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			fmt.Println("name:", info.Name())
			if info.Name() == filepath.Base(f.Root) && f.TryFile != nil {
				return file, nil
			}
			return nil, fs.ErrPermission
		}
	}
	return file, nil
}
