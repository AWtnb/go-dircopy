package dircopy_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWtnb/go-dircopy"
	"github.com/AWtnb/go-walk"
)

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
	err := dircopy.Copy(`C:\Personal\gotemp\hoge\cc.txt`, `C:\Personal\gotemp\cc.txt`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyFromChild(t *testing.T) {
	err := dircopy.Copy(`C:\Personal\gotemp\hoge\aa\dd`, `C:\Personal\gotemp\dd`)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyWarningSamePath(t *testing.T) {
	t.Log("This should raise error.")
	err := dircopy.Copy(`C:\Personal\gotemp\hoge`, `C:\Personal\gotemp\hoge`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}

func TestCopyWarningLoop(t *testing.T) {
	t.Log("This should raise error due to infinit loop.")
	err := dircopy.Copy(`C:\Personal\gotemp\hoge`, `C:\Personal\gotemp\hoge\piyo`)
	if err != nil {
		t.Logf("error was raised: %s", err)
	}
}
