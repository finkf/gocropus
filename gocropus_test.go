package gocropus

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOpenTxtFile(t *testing.T) {
	tests := []struct {
		test, want string
		wantErr    bool
	}{
		{"testdata/00001.gt.txt", "Fritſch, ein unverheyratheter Mann von hoͤchſt ein—", false},
		{"testdata/00001.txt", "Fritſch, ein unverheyratheter Mann von hochſt ein⸗", false},
		{"testdata/00001.xyz", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got, err := OpenTxtFile(tc.test)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("got error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestOpenLLocsFile(t *testing.T) {
	tests := []struct {
		test    string
		want    LLocs
		wantErr bool
	}{
		{"testdata/00003.llocs", LLocs{
			{'e', 60.7, 0},
			{'r', 77.8, 0},
			{'l', 91.8, 0},
			{'i', 107.3, 0},
			{'n', 121.3, 0},
			{'.', 143.1, 0},
		}, false},
		{"testdata/00004.llocs", LLocs{
			{'e', 60.7, 0.2},
			{'r', 77.8, 0.9},
			{'l', 91.8, 4e-06},
			{'i', 107.3, 0.1},
			{'n', 121.3, 0.1},
			{'.', 143.1, 0.1},
		}, false},
		{"invalid/00004.llocs", nil, true},
		{"testdata/00001.llocs", nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got, err := OpenLLocsFile(tc.test)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

func TestOpenLLocsFileString(t *testing.T) {
	tests := []struct {
		test, want string
	}{
		{"testdata/00003.llocs", "erlin."},
		{"testdata/00002.llocs", "n r in e  ein r i ch."},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			llocs, err := OpenLLocsFile(tc.test)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if got := llocs.String(); got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestOpenLLocsFileCuts(t *testing.T) {
	tests := []struct {
		test string
		want []int
	}{
		{"testdata/00003.llocs", []int{60, 77, 91, 107, 121, 143}},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			llocs, err := OpenLLocsFile(tc.test)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if got := llocs.Cuts(); !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

func TestOpenLLocsFileConfs(t *testing.T) {
	tests := []struct {
		test string
		want []float32
	}{
		{"testdata/00003.llocs", make([]float32, 6)},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			llocs, err := OpenLLocsFile(tc.test)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if got := llocs.Confs(); !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

func TestOpenImgFile(t *testing.T) {
	tests := []struct {
		test    string
		wantErr bool
	}{
		{"testdata/00001.nrm.png", false},
		{"testdata/00002.bin.png", false},
		{"testdata/00003.dew.png", false},
		{"testdata/00004.png", false},
		{"testdata/00007.png", true},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			_, err := OpenImgFile(tc.test)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("got error: %v", err)
			}
		})
	}
}

func TestStrip(t *testing.T) {
	tests := []struct {
		test, want string
	}{
		{"/foo/bar/abc.a", "/foo/bar/abc"},
		{"../foo.bar.baz", "../foo"},
		{".././foo/bar/.abc.a", ".././foo/bar/"},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			if got := Strip(tc.test); tc.want != got {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestImageFromFile(t *testing.T) {
	tests := []struct {
		test, want string
		ok         bool
	}{
		{"testdata/00001", "testdata/00001.nrm.png", true},
		{"testdata/00002", "testdata/00002.bin.png", true},
		{"testdata/00003", "testdata/00003.dew.png", true},
		{"testdata/00004", "testdata/00004.png", true},
		{"testdata/00005", "", false},
		{"testdata/00006", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got, ok := ImageFromFile(tc.test)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("expected (%q,%t); got (%q,%t)",
					tc.want, tc.ok, got, ok)
			}
		})
	}
}

func TestGTFromFile(t *testing.T) {
	tests := []struct {
		test, want string
		stat, ok   bool
	}{
		{"testdata/00001", "testdata/00001.gt.txt", true, true},
		{"testdata/00002", "testdata/00002.gt.txt", true, true},
		{"testdata/00003", "testdata/00003.gt.txt", true, true},
		{"testdata/00004", "testdata/00004.gt.txt", true, true},
		{"testdata/00001", "testdata/00001.gt.txt", false, false},
		{"testdata/00002", "testdata/00002.gt.txt", false, false},
		{"testdata/00003", "testdata/00003.gt.txt", false, false},
		{"testdata/00004", "testdata/00004.gt.txt", false, false},
		{"testdata/00005", "testdata/00005.gt.txt", true, false},
		{"testdata/00006", "testdata/00006.gt.txt", true, false},
		{"testdata/00005", "testdata/00005.gt.txt", false, false},
		{"testdata/00006", "testdata/00006.gt.txt", false, false},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got, ok := GTFromFile(tc.test, tc.stat)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("expected (%q,%t); got (%q,%t)",
					tc.want, tc.ok, got, ok)
			}
		})
	}
}

func TestTxtLocsFromFile(t *testing.T) {
	tests := []struct {
		test, want string
		stat, ok   bool
	}{
		{"testdata/00001", "testdata/00001.txt", true, true},
		{"testdata/00002", "testdata/00002.txt", true, true},
		{"testdata/00003", "testdata/00003.txt", true, true},
		{"testdata/00004", "testdata/00004.txt", true, true},
		{"testdata/00001", "testdata/00001.txt", false, false},
		{"testdata/00002", "testdata/00002.txt", false, false},
		{"testdata/00003", "testdata/00003.txt", false, false},
		{"testdata/00004", "testdata/00004.txt", false, false},
		{"testdata/00005", "testdata/00005.txt", true, false},
		{"testdata/00006", "testdata/00006.txt", true, false},
		{"testdata/00005", "testdata/00005.txt", false, false},
		{"testdata/00006", "testdata/00006.txt", false, false},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got, ok := TxtFromFile(tc.test, tc.stat)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("expected (%q,%t); got (%q,%t)",
					tc.want, tc.ok, got, ok)
			}
		})
	}
}

func TestLLocsFromFile(t *testing.T) {
	tests := []struct {
		test, want string
		stat, ok   bool
	}{
		{"testdata/00001", "testdata/00001.llocs", true, true},
		{"testdata/00002", "testdata/00002.llocs", true, true},
		{"testdata/00003", "testdata/00003.llocs", true, true},
		{"testdata/00004", "testdata/00004.llocs", true, true},
		{"testdata/00001", "testdata/00001.llocs", false, false},
		{"testdata/00002", "testdata/00002.llocs", false, false},
		{"testdata/00003", "testdata/00003.llocs", false, false},
		{"testdata/00004", "testdata/00004.llocs", false, false},
		{"testdata/00005", "testdata/00005.llocs", true, false},
		{"testdata/00006", "testdata/00006.llocs", true, false},
		{"testdata/00005", "testdata/00005.llocs", false, false},
		{"testdata/00006", "testdata/00006.llocs", false, false},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got, ok := LLocsFromFile(tc.test, tc.stat)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("expected (%q,%t); got (%q,%t)",
					tc.want, tc.ok, got, ok)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	tests := []struct {
		dir, ext, test string
		rec, want      bool
	}{
		{"testdata", TxtExt, "testdata/00001.gt.txt", false, true},
		{"testdata", BinPngExt, "testdata/00002.gt.txt", false, true},
		{"testdata", BinPngExt, "testdata/00001.gt.txt", false, false},
		{"testdata", PngExt, "testdata/00005.gt.txt", false, false},
		{"testdata", TxtExt, "testdata/00007.txt", false, true},
		{"testdata", TxtExt, "testdata/00007.gt.txt", false, false},
		// recursive
		{"./", TxtExt, "testdata/00001.gt.txt", true, true},
		{"./", BinPngExt, "testdata/00002.gt.txt", true, true},
		{"./", BinPngExt, "testdata/00001.gt.txt", true, false},
		{"./", PngExt, "testdata/00005.gt.txt", true, false},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			found := make(map[string]bool)
			err := Walk(tc.dir, tc.ext, tc.rec, func(gt, img, txt, llocs string) error {
				found[gt] = true
				found[img] = true
				found[txt] = true
				found[llocs] = true
				return nil
			})
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			found[""] = false // make sure that the empty string returns false
			if got := found[tc.test]; got != tc.want {
				t.Fatalf("expected found=%t; got found=%t", tc.want, got)
			}
		})
	}
}

func TestWalkError(t *testing.T) {
	tests := []struct {
		dir string
		f   WalkFunc
	}{
		{"testdata", func(string, string, string, string) error { return fmt.Errorf("error") }},
		{"not-exists", func(string, string, string, string) error { return nil }},
	}
	for _, tc := range tests {
		t.Run(tc.dir, func(t *testing.T) {
			err := Walk(tc.dir, TxtExt, false, tc.f)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	}
}
