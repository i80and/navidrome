package cache

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/djherbis/fscache"
	"github.com/karrick/godirwalk"
	"gopkg.in/djherbis/atime.v1"
	"gopkg.in/djherbis/stream.v1"
)

type spreadFS struct {
	root string
	mode os.FileMode
	init func() error
}

// NewSpreadFS returns a FileSystem rooted at directory dir. It
// Dir is created with perms if it doesn't exist.
func NewSpreadFS(dir string, mode os.FileMode) (fscache.FileSystem, error) {
	fs := &spreadFS{root: dir, mode: mode, init: func() error {
		return os.MkdirAll(dir, mode)
	}}
	return fs, fs.init()
}

func (fs *spreadFS) Reload(f func(key string, name string)) error {
	return godirwalk.Walk(fs.root, &godirwalk.Options{
		Callback: func(absoluteFilePath string, de *godirwalk.Dirent) error {
			path, err := filepath.Rel(fs.root, absoluteFilePath)
			if err != nil {
				return nil
			}

			// Skip if name is not in the format XX/XX/XXXXXXXXXXXX
			parts := strings.Split(path, string(os.PathSeparator))
			if len(parts) != 3 || len(parts[0]) != 2 || len(parts[1]) != 2 {
				return nil
			}

			f(absoluteFilePath, absoluteFilePath)
			return nil
		},
		Unsorted: true,
	})
}

func (fs *spreadFS) Create(name string) (stream.File, error) {
	path := filepath.Dir(name)
	err := os.MkdirAll(path, fs.mode)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
}

func (fs *spreadFS) Open(name string) (stream.File, error) {
	return os.Open(name)
}

func (fs *spreadFS) Remove(name string) error {
	return os.Remove(name)
}

func (fs *spreadFS) Stat(name string) (fscache.FileInfo, error) {
	stat, err := os.Stat(name)
	if err != nil {
		return fscache.FileInfo{}, err
	}
	return fscache.FileInfo{FileInfo: stat, Atime: atime.Get(stat)}, nil
}

func (fs *spreadFS) RemoveAll() error {
	if err := os.RemoveAll(fs.root); err != nil {
		return err
	}
	return fs.init()
}

func (fs *spreadFS) KeyMapper(key string) string {
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(key)))
	return filepath.Join(fs.root, hash[0:2], hash[2:4], hash)
}
