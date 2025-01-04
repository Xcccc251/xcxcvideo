package test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestWalkDir(t *testing.T) {
	number := getFilesNumberWithPrefix("C:\\Users\\86150\\GolandProjects\\XcXcVideo\\fileDir\\chunk", "chunk-")
	fmt.Println(number)

}
func getFilesNumberWithPrefix(dir string, prefix string) (number int) {
	//sw := sync.WaitGroup{}
	//c := make(chan struct{}, 10)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), prefix) {
			number++
			fmt.Println(d.Name())
		}
		return nil

	})
	if err != nil {
		panic(err)
	}
	return
}
