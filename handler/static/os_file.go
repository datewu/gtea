package static

import (
	"errors"
	"io/fs"
	"net/http"
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
					return good, err
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
			return nil, fs.ErrPermission
		}
	}
	return file, nil
}
