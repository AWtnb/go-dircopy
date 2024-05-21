package dircopy_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWtnb/go-dircopy"
	"github.com/AWtnb/go-walk"
)

func makeTestDir(path string) error {
	err := os.MkdirAll(path, 0700)
	return err
}

func makeTestFile(path string) error {
	_, err := os.Create(path)
	return err
}

func makeTestTree(root string) error {
	if err := makeTestDir(root); err != nil {
		return err
	}
	ds := []string{"aa/bb", "aa/cc", "bb/ee"}
	for _, d := range ds {
		p := filepath.Join(root, d)
		if err := makeTestDir(p); err != nil {
			return err
		}
	}
	fs := []string{"aa/bb/cc.txt", "aa/ff.txt", "dd.txt"}
	for _, f := range fs {
		p := filepath.Join(root, f)
		if err := makeTestFile(p); err != nil {
			return err
		}
	}
	return nil
}

func TestMakeTestTree(t *testing.T) {
	p := `C:\Personal\gotemp\piyo`
	err := makeTestTree(p)
	if err != nil {
		t.Error(err)
	}
	showTreeContent(t, p)
}

func getChildItems(root string) []string {
	d := walk.Dir{All: true, Root: root}
	d.SetWalkDepth(-1)
	d.SetWalkException("")
	found, err := d.GetChildItem()
	if err != nil {
		fmt.Println(err)
	}
	return found
}

func showTreeContent(t *testing.T, root string) {
	d := filepath.Dir(root)
	for _, c := range getChildItems(root) {
		rel := strings.TrimPrefix(c, d)
		t.Log(rel)
	}
}

func TestCopy(t *testing.T) {
	from := `C:\Personal\gotemp\hoge`
	if err := makeTestTree(from); err != nil {
		t.Error(err)
	}
	to := `C:\Personal\gotemp\fuga`
	err := dircopy.Copy(from, to)
	if err != nil {
		t.Error(err)
	}
	t.Log("ORIGINAL TREE:")
	showTreeContent(t, from)
	t.Log("COPIED TREE:")
	showTreeContent(t, to)
}

func TestCopyFile(t *testing.T) {
	t.Log("this should raise error.")
	d := `C:\Personal\gotemp\hoge`
	makeTestDir(d)
	f := filepath.Join(d, "cc.txt")
	makeTestFile(f)
	err := dircopy.Copy(f, `C:\Personal\gotemp\cc.txt`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyFromChild(t *testing.T) {
	p := `C:\Personal\gotemp\hoge\aa\dd`
	makeTestDir(p)
	err := dircopy.Copy(p, `C:\Personal\gotemp\dd`)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyFromChildWarning(t *testing.T) {
	t.Log("This should raise error.")
	p := `C:\Personal\gotemp\hoge\aa\dd`
	makeTestDir(p)
	err := dircopy.Copy(p, `C:\Personal\gotemp`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyWarningSamePath(t *testing.T) {
	t.Log("This should raise error.")
	p := `C:\Personal\gotemp\hoge`
	makeTestDir(p)
	err := dircopy.Copy(p, p)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyWarningLoop(t *testing.T) {
	t.Log("This should raise error due to infinit loop.")
	d1 := `C:\Personal\gotemp\hoge`
	makeTestDir(d1)
	d2 := `C:\Personal\gotemp\hoge\piyo`
	makeTestDir(d2)
	err := dircopy.Copy(d1, d2)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}
