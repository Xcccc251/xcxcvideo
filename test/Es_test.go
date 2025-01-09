package test

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/es"
	"XcxcVideo/common/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestEs(t *testing.T) {
	esClient := es.Init()
	indexName := "xcxc_video"
	res, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode == 404 {
		// 如果索引不存在，则创建它
		res, err = esClient.Indices.Create(indexName)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		defer res.Body.Close()
		fmt.Printf("Index %s created\n", indexName)
	} else {
		fmt.Printf("Index %s already exists\n", indexName)
	}

	var videoVoList []models.VideoVo
	models.Db.Model(new(models.VideoVo)).Find(&videoVoList)
	for _, videoVo := range videoVoList {
		videoVoJson, err := json.Marshal(videoVo)
		if err != nil {
			t.Error(err)
		}
		_, err = esClient.Index(indexName, bytes.NewReader(videoVoJson), esClient.Index.WithDocumentID(strconv.Itoa(videoVo.Vid)))
		if err != nil {
			t.Error(err)
		}

	}
	// 查询索引中的所有文档
	searchResult, err := esClient.Search(
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(`{"query": {"match_all": {}}}`)),
	)
	if err != nil {
		t.Error(err)
	}
	defer searchResult.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, searchResult.Body)
	if err != nil {
		t.Error(err)
	}

	var searchResponse map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &searchResponse)
	if err != nil {
		t.Error(err)
	}

	// 打印查询结果
	fmt.Println(searchResponse)

}

func TestSearch(t *testing.T) {
	esClient := es.Init()
	indexName := "xcxc_video"
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy": map[string]interface{}{
				"title": map[string]interface{}{
					"value":     "晨光", // 模糊查询的关键词
					"fuzziness": "1",  // 自动模糊级别，支持 0、1、2 或 "AUTO"
				},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
		},
	}

	// 将查询转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索请求
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 解析搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	videoVoList := []models.VideoVo{}
	// 打印匹配的结果和高亮部分
	fmt.Printf("Found %d hits:\n", int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		highlight := hitMap["highlight"].(map[string]interface{})
		var videoVo models.VideoVo
		//封装结构体
		videoVoList = append(videoVoList, videoVo)
		// 打印文档信息
		fmt.Printf("Document: %v\n", source)

		// 打印高亮部分
		for field, highlights := range highlight {
			fmt.Printf("Highlighted %s: %v\n", field, highlights)
		}
	}
	fmt.Println(videoVoList)
}
func TestSearch4(t *testing.T) {
	esClient := es.Init()
	indexName := "xcxc_video"
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy": map[string]interface{}{
				"title": map[string]interface{}{
					"value":     "晨光", // 模糊查询的关键词
					"fuzziness": "1",  // 自动模糊级别，支持 0、1、2 或 "AUTO"
				},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
		},
	}

	// 将查询转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索请求
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 解析搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source    models.VideoVo      `json:"_source"`
				Highlight map[string][]string `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	videoVoList := make([]models.VideoVo, 0)
	fmt.Printf("Found %d hits:\n", r.Hits.Total.Value)

	for _, hit := range r.Hits.Hits {
		video := hit.Source

		// 处理高亮字段，假设需要将高亮的内容保存到结构体中
		if highlights, ok := hit.Highlight["title"]; ok && len(highlights) > 0 {
			video.Title = highlights[0] // 将第一个高亮片段存入结构体
		}

		videoVoList = append(videoVoList, video)

		// 打印文档信息和高亮字段
		fmt.Printf("Document: %+v\n", video)
		if len(hit.Highlight) > 0 {
			fmt.Printf("Highlight: %+v\n", hit.Highlight)
		}
	}

	fmt.Println("Final Result:", videoVoList)
}

func TestEs2(t *testing.T) {
	esClient := es.Init()
	indexName := define.ES_SEARCH_WORD
	res, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode == 404 {
		// 如果索引不存在，则创建它
		res, err = esClient.Indices.Create(indexName)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		defer res.Body.Close()
		fmt.Printf("Index %s created\n", indexName)
	} else {
		fmt.Printf("Index %s already exists\n", indexName)
	}
}

func TestEsAdd(t *testing.T) {
	esClient := es.Init()
	searchWord := models.EsSearchWord{}
	searchWord.Content = "测试"
	searchWordJson, err := json.Marshal(searchWord)

	_, err = esClient.Index(define.ES_SEARCH_WORD, bytes.NewReader(searchWordJson))
	if err != nil {
		t.Error(err)
	}

	// 查询索引中的所有文档
	searchResult, err := esClient.Search(
		esClient.Search.WithIndex(define.ES_SEARCH_WORD),
		esClient.Search.WithBody(strings.NewReader(`{"query": {"match_all": {}}}`)),
	)
	if err != nil {
		t.Error(err)
	}
	defer searchResult.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, searchResult.Body)
	if err != nil {
		t.Error(err)
	}

	var searchResponse map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &searchResponse)
	if err != nil {
		t.Error(err)
	}

	// 打印查询结果
	fmt.Println(searchResponse)

}

func TestEsDel(t *testing.T) {
	esClient := es.Init()
	ids := []string{"2", "4", "10"}
	for _, id := range ids {
		_, err := esClient.Delete(define.ES_VIDEO, id)
		if err != nil {
			t.Error(err)
		}
	}

}

func TestEsAdd3(t *testing.T) {
	esClient := es.Init()
	indexName := "xcxc_video"

	videoVo := models.VideoVo{}
	videoVo.Vid = 10
	videoVo.Title = "测试"
	videoVoJson, err := json.Marshal(videoVo)
	if err != nil {
		t.Error(err)
	}
	_, err = esClient.Index(indexName, bytes.NewReader(videoVoJson), esClient.Index.WithDocumentID(strconv.Itoa(videoVo.Vid)))
	if err != nil {
		t.Error(err)
	}

	// 查询索引中的所有文档
	searchResult, err := esClient.Search(
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(`{"query": {"match_all": {}}}`)),
	)
	if err != nil {
		t.Error(err)
	}
	defer searchResult.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, searchResult.Body)
	if err != nil {
		t.Error(err)
	}

	var searchResponse map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &searchResponse)
	if err != nil {
		t.Error(err)
	}

	// 打印查询结果
	fmt.Println(searchResponse)
}

func TestEsUpdate(t *testing.T) {
	esClient := es.Init()
	indexName := "xcxc_video"
	updateBody := map[string]interface{}{
		"doc": map[string]interface{}{
			"status": 0,
		},
	}
	searchWordJson, err := json.Marshal(updateBody)
	if err != nil {
		t.Error(err)
	}
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: "4",
		Body:       bytes.NewReader(searchWordJson),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		t.Errorf("Error updating document: %s", res.String())
	} else {
		var response map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			t.Error(err)
		}
		t.Logf("Update response: %v", response)
	}

}

func TestEsSearch2(t *testing.T) {
	esClient := es.Init()
	indexName := "xcxc_search_word"
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"content": "晨", // 支持中文分词的全文搜索
						},
					},
					{
						"match_phrase_prefix": map[string]interface{}{
							"content": "晨", // 匹配以“测试”开头的分词
						},
					},
				},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
		},
	}

	// 将查询转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索请求
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 解析搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	searchWordList := []string{}
	// 打印匹配的结果和高亮部分
	fmt.Printf("Found %d hits:\n", int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		highlight := hitMap["highlight"].(map[string]interface{})
		searchWordList = append(searchWordList, source["content"].(string))

		// 打印文档信息
		fmt.Printf("Document: %v\n", source)

		// 打印高亮部分
		for field, highlights := range highlight {
			fmt.Printf("Highlighted %s: %v\n", field, highlights)
		}
	}
}
