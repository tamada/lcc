package wildcat

import (
	"os"
	"strings"
	"testing"
)

func match(list []string, wonts []string) bool {
	for _, wont := range wonts {
		found := false
		for _, item := range list {
			if strings.Contains(item, wont) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestStdin(t *testing.T) {
	testdata := []struct {
		stdinFile     string
		opts          *ReadOptions
		listSize      int
		wontFileNames []string
		wontErrorSize int
	}{
		{"testdata/wc/london_bridge_is_broken_down.txt", &ReadOptions{false, true, true}, 1, []string{"<stdin>"}, 0},
		{"testdata/filelist.txt", &ReadOptions{true, false, false}, 3, []string{"humpty_dumpty.txt", "sakura_sakura.txt", "london_bridge_is_broken_down.txt"}, 0},
	}
	for _, td := range testdata {
		file, _ := os.Open(td.stdinFile)
		origStdin := os.Stdin
		os.Stdin = file
		defer func() {
			os.Stdin = origStdin
			file.Close()
		}()
		args := NewArguments()
		args.Options = td.opts
		ec := NewErrorCenter()
		rs := args.CountAll(func() Counter { return NewCounter(All) }, ec)

		if len(rs.list) != td.listSize {
			t.Errorf("ResultSet size did not match, wont %d, got %d (%v)", td.listSize, len(rs.list), rs.list)
		}
		if !match(rs.list, td.wontFileNames) {
			t.Errorf("ResultSet files did not match, wont %v, got %v", td.wontFileNames, rs.list)
		}
		if len(ec.errs) != td.wontErrorSize {
			t.Errorf("ErrorSize did not match, wont %d, got %d (%v)", td.wontErrorSize, len(ec.errs), ec.errs)
		}
	}
}

func TestCountAll(t *testing.T) {
	testdata := []struct {
		args          []string
		opts          *ReadOptions // FileList, NoIgnore, NoExtract
		listSize      int
		wontFileNames []string
		wontErrorSize int
	}{
		{[]string{"testdata/wc"}, &ReadOptions{false, false, false}, 3, []string{"humpty_dumpty.txt", "sakura_sakura.txt", "london_bridge_is_broken_down.txt"}, 0},
		{[]string{"https://www.apache.org/licenses/LICENSE-2.0.txt"}, &ReadOptions{false, false, false}, 1, []string{"https://www.apache.org/licenses/LICENSE-2.0.txt"}, 0},
		{[]string{"testdata/ignores"}, &ReadOptions{false, false, false}, 2, []string{"notIgnore.txt", "notIgnore_sub.txt"}, 0},
		{[]string{"testdata/ignores"}, &ReadOptions{false, true, false}, 7, []string{"ignore.test", "ignore.test2", "notIgnore.txt", "notIgnore_sub.txt", "ignore_sub.test"}, 0},
		{[]string{"testdata/filelist.txt"}, &ReadOptions{true, false, false}, 3, []string{"humpty_dumpty.txt", "sakura_sakura.txt", "london_bridge_is_broken_down.txt"}, 0},
		{[]string{"testdata/not_found.txt"}, &ReadOptions{true, false, false}, 0, []string{}, 1},
		{[]string{"https://example.com/not_found"}, &ReadOptions{false, false, false}, 0, []string{}, 1},
	}
	for _, td := range testdata {
		args := NewArguments()
		args.Args = td.args
		args.Options = td.opts
		ec := NewErrorCenter()
		rs := args.CountAll(func() Counter { return NewCounter(All) }, ec)

		if len(rs.list) != td.listSize {
			t.Errorf("ResultSet size did not match, wont %d, got %d (%v)", td.listSize, len(rs.list), rs.list)
		}
		if !match(rs.list, td.wontFileNames) {
			t.Errorf("ResultSet files did not match, wont %v, got %v", td.wontFileNames, rs.list)
		}
		if len(ec.errs) != td.wontErrorSize {
			t.Errorf("ErrorSize did not match, wont %d, got %d (%v)", td.wontErrorSize, len(ec.errs), ec.errs)
		}
	}
}