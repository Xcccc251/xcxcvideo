package service

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func askCurrentChunkByHash(hash string) {

}
func getFilesNumberWithPrefix(dir string, prefix string) (number int) {
	//sw := sync.WaitGroup{}
	//c := make(chan struct{}, 10)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(d.Name(), prefix) {
			number++
		}
		return nil

	})
	if err != nil {
		panic(err)
	}
	return
}
