package wildcat

import (
	"path/filepath"
	"testing"
)

func TestReadFilesInDir(t *testing.T) {
	ec := NewErrorCenter()
	targets := readFilesInDir("testdata", ec)
	if len(targets) != 3 {
		t.Errorf("readFilesInDir(\"testdata\") size did not match, wont %d, got %d", 3, len(targets))
	}
	testdata := []string{"ja/sakura_sakura.txt", "humpty_dumpty.txt", "london_bridge_is_broken_down.txt"}
	for _, td := range testdata {
		if !Contains(targets, filepath.Join("testdata/", td)) {
			t.Errorf("readFilesInDir did not contains %s", td)
		}
	}
}