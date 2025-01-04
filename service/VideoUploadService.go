package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/minIO"
	"XcxcVideo/common/models"
	"XcxcVideo/common/oss"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func AskCurrentChunkByHash(c *gin.Context) {
	hash := c.Query("hash")
	number := getFilesNumberWithPrefix(define.CHUNK_PATH, "chunk-"+hash+"-")
	if number == 0 {
		response.ResponseOKWithData(c, "获取成功", 0)
	} else {
		response.ResponseOKWithData(c, "获取成功", number)
	}
	return

}

func UploadVideoChunk(c *gin.Context) {
	chunk, _ := c.FormFile("chunk")
	hash := c.PostForm("hash")
	index, _ := strconv.Atoi(c.PostForm("index"))
	chunkFileName := "chunk-" + hash + "-" + strconv.Itoa(index)
	chunkPath := define.CHUNK_PATH + "/" + chunkFileName
	err := c.SaveUploadedFile(chunk, chunkPath)
	if err != nil {
		response.ResponseFailWithData(c, http.StatusInternalServerError, "保存文件失败", nil)
		return
	}
	response.ResponseOK(c)
	return
}

func UploadVideo(c *gin.Context) {
	cover, _ := c.FormFile("cover")
	hash := c.PostForm("hash")
	title := c.PostForm("title")
	typeOfVideo := c.PostForm("type")
	auth := c.PostForm("auth")
	duration := c.PostForm("duration")
	mcid := c.PostForm("mcid")
	scid := c.PostForm("scid")
	tags := c.PostForm("tags")
	descr := c.PostForm("descr")
	userId, _ := c.Get("userId")
	title = strings.Trim(title, " ")
	if title == "" {
		response.ResponseFailWithData(c, 500, "标题不能为空", nil)
		return
	}
	if len(title) > 80 {
		response.ResponseFailWithData(c, 500, "标题过长", nil)
		return
	}
	if len(descr) > 2000 {
		response.ResponseFailWithData(c, 500, "简介过长", nil)
		return
	}
	var videoUploadInfoDto models.VideoUploadInfoDto
	videoUploadInfoDto.Auth, _ = strconv.Atoi(auth)
	videoUploadInfoDto.Duration, _ = strconv.ParseFloat(duration, 64)
	videoUploadInfoDto.McId = mcid
	videoUploadInfoDto.ScId = scid
	videoUploadInfoDto.Type, _ = strconv.Atoi(typeOfVideo)
	videoUploadInfoDto.Tags = tags
	videoUploadInfoDto.Title = title
	videoUploadInfoDto.Hash = hash
	coverUrl, _ := oss.UploadImage(cover)
	videoUploadInfoDto.CoverUrl = coverUrl
	videoUploadInfoDto.Uid = userId.(int)
	go func() {
		videoUrl, _ := mergeSlices("chunk-" + hash)
		var video models.Video
		copier.Copy(&video, &videoUploadInfoDto)
		video.VideoUrl = videoUrl
		models.Db.Model(new(models.Video)).Create(&video)
		var videoStats models.VideoStats
		videoStats.Vid = video.Vid
		models.Db.Model(new(models.VideoStats)).Create(&videoStats)
		var videoVo models.VideoVo
		models.Db.Model(new(models.VideoVo)).Where("id = ?", video.Vid).Find(&videoVo)
		//todo es
		go func() {
			videoJson, _ := json.Marshal(videoVo)
			videoStatsJson, _ := json.Marshal(videoStats)
			models.RDb.Set(context.Background(), define.VIDEOSTATS_PREFIX+strconv.Itoa(video.Vid), videoStatsJson, 0)
			models.RDb.Set(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(video.Vid), videoJson, 0)
			models.RDb.SAdd(context.Background(), define.VIDEO_STATUS_0, video.Vid)
		}()
	}()
	response.ResponseOKWithData(c, "上传成功", nil)
	return

}

func mergeSlices(prefix string) (videoPath string, err error) {
	// 搜索以 hash 开头的文件
	files, err := findFilesByHash(prefix, define.CHUNK_PATH)
	if err != nil {
		return "", fmt.Errorf("搜索文件失败: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("未找到以哈希值 %s 开头的文件", prefix)
	}

	// 合并文件到目标路径
	outputFile := filepath.Join(define.VIDEO_PATH, prefix+"-merged.mp4")
	out, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer out.Close()

	for _, file := range files {
		if err := appendFileToOutput(file, out); err != nil {
			return "", err
		}
	}

	url, err := minIO.UploadMP4(prefix+"-merged.mp4", out)
	if err != nil {
		return "", err
	}
	go func() {
		//todo kafka
		filepath.WalkDir(define.CHUNK_PATH, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasPrefix(d.Name(), prefix) {
				os.Remove(path)
			}
			return nil
		})
	}()

	return url, nil
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

func CancelUpload(c *gin.Context) {
	hash := c.Query("hash")
	err := filepath.WalkDir(define.CHUNK_PATH, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), "chunk-"+hash+"-") {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		response.ResponseFailWithData(c, http.StatusInternalServerError, "取消上传失败", nil)
		return
	}
	response.ResponseOKWithData(c, "取消上传成功", nil)
	return

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
