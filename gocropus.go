package gocropus // import "github.com/finkf/gocropus"
import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LLoc represents character information for one recognized character.
type LLoc struct {
	Char rune
	Cut  float32
	Conf float32
}

// LLocs represents character information for one recognized line.
type LLocs []LLoc

func (l LLocs) String() string {
	wstr := make([]rune, len(l))
	for i := range l {
		wstr[i] = l[i].Char
	}
	return string(wstr)
}

// Cuts returns the right cuts of the llocs as int array.
func (l LLocs) Cuts() []int {
	cuts := make([]int, len(l))
	for i := range l {
		cuts[i] = int(l[i].Cut)
	}
	return cuts
}

// Confs returns the confidences as an array.
func (l LLocs) Confs() []float32 {
	confs := make([]float32, len(l))
	for i := range l {
		confs[i] = l[i].Conf
	}
	return confs
}

// OpenLLocsFile opens a llocs file and returns its contents.
func OpenLLocsFile(path string) (LLocs, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %v", path, err)
	}
	defer in.Close()
	return ReadLLocsFile(in)
}

// ReadLLocsFile read the contents from a llocs file.
func ReadLLocsFile(in io.Reader) (LLocs, error) {
	var res LLocs
	s := bufio.NewScanner(in)
	for s.Scan() {
		line := s.Text()
		if line == "" || line[0] == '\t' { // skip bad lines
			continue
		}
		var lloc LLoc
		if _, err := fmt.Sscanf(line, "%c\t%f\t%f", &lloc.Char, &lloc.Cut, &lloc.Conf); err == nil {
			res = append(res, lloc)
			continue
		}
		if _, err := fmt.Sscanf(line, "%c\t%f", &lloc.Char, &lloc.Cut); err == nil {
			res = append(res, lloc)
			continue
		}
		return nil, fmt.Errorf("cannot parse line: %q", line)
	}
	return res, s.Err()
}

// OpenTxtFile opens a txt or gt file and reads it content line.
func OpenTxtFile(path string) (string, error) {
	in, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %v", path, err)
	}
	defer in.Close()
	return ReadTxtFile(in)
}

// ReadTxtFile read the content line from a txt or gt file.
func ReadTxtFile(in io.Reader) (string, error) {
	var res string
	s := bufio.NewScanner(in)
	if s.Scan() {
		res = s.Text()
	}
	return res, s.Err()
}

// OpenImgFile reads the image's data from a png encoded file.
func OpenImgFile(path string) (image.Image, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %v", path, err)
	}
	defer in.Close()
	return ReadImgFile(in)
}

// ReadImgFile reads the image's data from a png encoded file.
func ReadImgFile(in io.Reader) (image.Image, error) {
	return png.Decode(in)
}

// Strip returns the bare file path for a given path with all
// extensions stripped.  If the path's file name starts with a leading
// dot, the whole file name will be removed.
func Strip(p string) string {
	for ext := filepath.Ext(p); ext != ""; ext = filepath.Ext(p) {
		p = p[0 : len(p)-len(ext)]
	}
	return p
}

// ImageExtensions defines the different possible extensions for line
// image files.  The order of the extensions defines which files are
// used for image files.  Change this if you need other image file
// priorities.
var ImageExtensions = []string{
	BinPngExt,
	DewPngExt,
	PngExt,
	NrmPngExt,
}

// File extensions for gt, img, txt and llocs files.
const (
	GTExt     = ".gt.txt"
	TxtExt    = ".txt"
	LLocsExt  = ".llocs"
	BinPngExt = ".bin.png"
	DewPngExt = ".dew.png"
	PngExt    = ".png"
	NrmPngExt = ".nrm.png" /* GT4HistOCR */
)

// ImageFromStripped returns the according line image file for the
// given stripped or unstripped path and whether it exists.  In order
// to identify the right extension, the file path is checked with
// stat.  If no existing image file path can be found this function
// returns "", false.
func ImageFromStripped(stripped string) (string, bool) {
	for _, ext := range ImageExtensions {
		p := stripped + ext
		if isFile(p) {
			return p, true
		}
	}
	return "", false
}

// GTFromStripped returns the according gt file for the given stripped
// or unstripped path and whether it exists.  If stat is set to false,
// just the according gt path and false are returend; it is not
// checked in this case if the resulting file path exists.  In any
// case the according gt file path is returned.
func GTFromStripped(p string, stat bool) (string, bool) {
	return checkStrippedWithExt(p, GTExt, stat)
}

// TxtFromStripped returns the according txt file for the given
// stripped or unstripped path and whether it exists.  If stat is set
// to false, just the according gt path and false are returend; it is
// not checked in this case if the resulting file path exists.  In any
// case the according txt file path is returned.
func TxtFromStripped(p string, stat bool) (string, bool) {
	return checkStrippedWithExt(p, TxtExt, stat)
}

// LLocsFromStripped returns the according llocs file for the given
// stripped path and whether it exists.  If stat is set to false, just
// the according gt path and false are returend; it is not checked in
// this case if the resulting file path exists.  In any case the
// according llocs file path is returned.
func LLocsFromStripped(p string, stat bool) (string, bool) {
	return checkStrippedWithExt(p, LLocsExt, stat)
}

// WalkFunc defines the type for the callback function used in Walk.
// It is called with the paths of the existing Ocropy file image
// set. The first path is the gt, the second path is the img, the
// third path is the txt and the fourth path is the llocs file path.
// If any file path does not exist the according value is set to the
// empty string "".
type WalkFunc func(string, string, string, string) error

// Walk iterates over all files in the given directory and calls the
// given callback function for each set of Ocropy files.  The set of
// Ocropy files is calculated based on the the given file extension.
// If recursive is false, sub directories are ignored.
func Walk(dir, ext string, recursive bool, f WalkFunc) error {
	err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if p == dir {
			return nil
		}
		if fi.IsDir() && recursive || p == dir { // do not skip sub dirs
			return nil
		}
		if fi.IsDir() { // do skip dir
			return filepath.SkipDir
		}
		if !strings.HasSuffix(p, ext) { // keep on looking
			return nil
		}
		// we have a valid file; check for the other files of the set
		s := Strip(p)
		gt := pathString(GTFromStripped(s, true))
		img := pathString(ImageFromStripped(s))
		txt := pathString(TxtFromStripped(s, true))
		llocs := pathString(LLocsFromStripped(s, true))
		return f(gt, img, txt, llocs)
	})
	return err
}

// pathString returns the empty string if ok is false; else the given
// string is returned.
func pathString(p string, ok bool) string {
	if ok {
		return p
	}
	return ""
}

// checkStrippedWithExt checks a stripped file with the given
// extension and returns if the file exists or not, depending on the
// stat argument.
func checkStrippedWithExt(p, ext string, stat bool) (string, bool) {
	p = Strip(p) + ext
	return p, stat && isFile(p)
}

// isFile returns true if the given path exists and is not a
// directory.
func isFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}
