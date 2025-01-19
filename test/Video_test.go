package test

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	mergeSlices("chunk-" + "d8f4477e6e2eadf67026ceab0806ba09")
}
func TestAdd(t *testing.T) {
	var videoList []models.Video
	var videoStatsList []models.VideoStats
	models.Db.Model(new(models.Video)).Find(&videoList)
	models.Db.Model(new(models.VideoStats)).Find(&videoStatsList)
	for _, v := range videoList {
		vJson, _ := json.Marshal(v)
		models.RDb.Set(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(v.Vid), vJson, 0)
	}
	for _, v := range videoStatsList {
		vJson, _ := json.Marshal(v)
		models.RDb.Set(context.Background(), define.VIDEOSTATS_PREFIX+strconv.Itoa(v.Vid), vJson, 0)
	}
}
func mergeSlices(prefix string) error {
	// 搜索以 hash 开头的文件
	files, err := findFilesByHash(prefix, define.CHUNK_PATH)
	if err != nil {
		return fmt.Errorf("搜索文件失败: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("未找到以哈希值 %s 开头的文件", prefix)
	}

	// 合并文件到目标路径
	outputFile := filepath.Join(define.VIDEO_PATH, prefix+"-merged.mp4")
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer out.Close()

	for _, file := range files {
		if err := appendFileToOutput(file, out); err != nil {
			return err
		}
	}

	return nil
}

func findFilesByHash(prefix string, directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 按文件名排序以确保正确顺序
	sort.Strings(files)

	return files, nil
}

func appendFileToOutput(sliceFile string, out *os.File) error {
	in, err := os.Open(sliceFile)
	if err != nil {
		return fmt.Errorf("打开切片文件 %s 失败: %w", sliceFile, err)
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}
	return nil
}

func isContain(a []int, x int) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func TestContain(t *testing.T) {
	fmt.Println(isContain([]int{1}, 1))
}
