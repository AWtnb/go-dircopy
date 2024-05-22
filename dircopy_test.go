package dircopy_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWtnb/go-dircopy"
	"github.com/AWtnb/go-testtree"
)

func makeTestTree(root string) error {
	ds := []string{"aa/bb", "aa/cc", "bb/ee"}
	fs := []string{"aa/bb/cc.txt", "aa/ff.txt", "dd.txt"}
	return testtree.MakeTestTree(root, ds, fs)
}

func TestMakeTestTree(t *testing.T) {
	p := `C:\Personal\gotemp\piyo`
	err := makeTestTree(p)
	if err != nil {
		t.Error(err)
	}
	showTreeContent(t, p)
}

func showTreeContent(t *testing.T, root string) {
	d := filepath.Dir(root)
	for _, c := range testtree.GetChildItems(root) {
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
	testtree.MakeTestDir(d)
	f := filepath.Join(d, "cc.txt")
	testtree.MakeTestFile(f)
	err := dircopy.Copy(f, `C:\Personal\gotemp\cc.txt`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyFromChild(t *testing.T) {
	p := `C:\Personal\gotemp\hoge\aa\dd`
	testtree.MakeTestDir(p)
	err := dircopy.Copy(p, `C:\Personal\gotemp\dd`)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyFromChildWarning(t *testing.T) {
	t.Log("This should raise error.")
	p := `C:\Personal\gotemp\hoge\aa\dd`
	testtree.MakeTestDir(p)
	err := dircopy.Copy(p, `C:\Personal\gotemp`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyWarningSamePath(t *testing.T) {
	t.Log("This should raise error.")
	p := `C:\Personal\gotemp\hoge`
	testtree.MakeTestDir(p)
	err := dircopy.Copy(p, p)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyWarningLoop(t *testing.T) {
	t.Log("This should raise error due to infinit loop.")
	d1 := `C:\Personal\gotemp\hoge`
	testtree.MakeTestDir(d1)
	d2 := `C:\Personal\gotemp\hoge\piyo`
	testtree.MakeTestDir(d2)
	err := dircopy.Copy(d1, d2)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}
