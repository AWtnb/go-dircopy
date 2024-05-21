package dircopy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Copy(src string, newPath string) error {
	if isFile(src) {
		return fmt.Errorf("non dir-path: %s", src)
	}
	if isFile(newPath) {
		return fmt.Errorf("non dir-path: %s", newPath)
	}
	if src == newPath {
		return fmt.Errorf("two args are the same path")
	}
	if strings.HasPrefix(newPath, src) {
		return fmt.Errorf("coping '%s' as its own subdirectory '%s' will cause infinit-loop", src, newPath)
	}
	if strings.HasPrefix(src, newPath) {
		return fmt.Errorf("creating '%s' may remove current directory tree '%s'", newPath, src)
	}
	if _, err := os.Stat(newPath); err == nil {
		err := os.RemoveAll(newPath)
		if err != nil {
			return fmt.Errorf("failed to remove pre-existing dest path")
		}
	}
	return copyItem(src, newPath)
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func isLink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0 || fi.Mode()&os.ModeDevice != 0
}

func copyItem(src string, newPath string) error {
	if isLink(src) {
		return fmt.Errorf("'%s' is a link to atnother location", src)
	}
	fs, err := os.Stat(src)
	if err != nil {
		return err
	}

	if fs.IsDir() {
		return addDir(src, newPath)
	}

	return addFile(src, newPath)
}

func addDir(src string, newPath string) error {
	fi, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	fs, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Mkdir(newPath, fs.Mode()&os.ModePerm); err != nil {
		return err
	}

	for _, f := range fi {
		sp := filepath.Join(src, f.Name())
		np := filepath.Join(newPath, f.Name())
		err := copyItem(sp, np)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFile(src string, newPath string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	nf, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer nf.Close()
	if _, err = io.Copy(nf, sf); err != nil {
		return err
	}
	return nil
}
