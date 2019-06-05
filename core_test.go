package exportlistjson_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/nekonbu72/exportlistjson"
	"github.com/nekonbu72/mailg"
	"github.com/nekonbu72/xemlsx"
	"github.com/tealeg/xlsx"
)

const dir string = "test"

func testPaths() []string {
	return dirwalk(dir)
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func TestPaths(t *testing.T) {
	for _, s := range testPaths() {
		fmt.Println(s)
	}
}

func testFiles() []*xlsx.File {
	var fs []*xlsx.File
	for _, p := range testPaths() {
		f, err := xlsx.OpenFile(p)
		if err != nil {
			continue
		}
		fs = append(fs, f)
	}
	return fs
}

func TestFiles(t *testing.T) {
	fs := testFiles()
	if len(fs) != 3 {
		t.Errorf("len: %v\n", len(fs))
	}
}

func testXLSXStream(done chan interface{}) <-chan *xemlsx.XLSX {
	xlsxStream := make(chan *xemlsx.XLSX)
	go func() {
		defer close(xlsxStream)
		for i, f := range testFiles() {
			select {
			case <-done:
				return
			case xlsxStream <- &xemlsx.XLSX{
				Attachment: &mailg.Attachment{
					Filename: strconv.Itoa(i),
					Reader:   nil,
				},
				File: f,
			}:
			}
		}
	}()
	return xlsxStream
}

func TestToJSON(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	ch := testXLSXStream(done)
	ch2 := exportlistjson.ToJSON(done, ch)

	var js []string
	for j := range ch2 {
		js = append(js, j)
	}

	if len(js) != len(testPaths()) {
		t.Errorf("len: %v\n", len(js))
	}
}