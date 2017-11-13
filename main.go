package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/trace"

	"github.com/ifreddyrondon/gocapture/capture"
)

func getJsonPath(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

func main() {
	trace.Start(os.Stdout)
	defer trace.Stop()

	p := new(capture.Path)
	filePath := filepath.Join("capture", "testdata", "merida_path_no_date_fixture.json")
	p.UnmarshalJSON(getJsonPath(filePath))
}
