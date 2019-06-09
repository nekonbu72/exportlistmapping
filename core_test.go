package exportlistmapping

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/nekonbu72/xemlsx"
	"github.com/tealeg/xlsx"
)

const (
	dir = "test"
)

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
				FileName: "File#" + strconv.Itoa(i),
				File:     f,
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
	ch2 := ToJSON(done, ch)

	var js []string
	for j := range ch2 {
		js = append(js, j)
	}

	if len(js) != len(testPaths()) {
		t.Errorf("len: %v\n", len(js))
	}

	for _, j := range js {
		log.Println(j)
	}
}

func TestToData2(t *testing.T) {
	s, _ := NewSetting()

	done := make(chan interface{})
	defer close(done)

	xlsxStream := testXLSXStream(done)
	xlsxDataStream := toXLSXData(done, s, xlsxStream)
	dataStream := toData2(done, xlsxDataStream)

	var data []*Data
	for d := range dataStream {
		data = append(data, d)
	}

	for _, d := range data {
		log.Println(d)
	}
}
