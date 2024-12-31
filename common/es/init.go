package es

import (
	"MyBlogv2/blog-common/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"strconv"
)

var esClient = Init()

func Init() *elasticsearch.Client {
	// 初始化 Elasticsearch 客户端
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9201", // 替换为你的 Elasticsearch 地址
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 检查 Elasticsearch 集群健康状态
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	return esClient
}

func SaveArticle(article models.ArticleEs, id int) {
	fmt.Println(article.Id)
	// 将 Article 转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(article); err != nil {
		log.Fatalf("Error encoding article: %s", err)
	}

	// 将文档存储到 Elasticsearch 的 "articles" 索引中
	res, err := esClient.Index(
		"article_es", // 索引名称
		&buf,         // 文档内容
		esClient.Index.WithDocumentID(fmt.Sprintf("%d", id)), // 使用文章 ID 作为文档 ID
		esClient.Index.WithContext(context.Background()),     // 上下文
	)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	// 打印响应
	log.Println(res)
	return
}

func UpdateArticle(article models.ArticleEs, id int) {
	// 将 Article 转换为 JSON 用于更新的内容
	updateBody := map[string]interface{}{
		"doc": article, // 使用 "doc" 字段来指定部分更新内容
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updateBody); err != nil {
		log.Fatalf("Error encoding update body: %s", err)
	}
	fmt.Println(buf)

	// 使用 Elasticsearch 的 Update API 更新文档
	res, err := esClient.Update(
		"article_es",          // 索引名称
		fmt.Sprintf("%d", id), // 文档 ID
		&buf,                  // 更新的部分内容
		esClient.Update.WithContext(context.Background()), // 上下文
	)
	if err != nil {
		log.Fatalf("Error updating document: %s", err)
	}
	defer res.Body.Close()

	// 打印响应
	log.Println(res)
	return
}

func SearchArticleByContent(content string) ([]models.ArticleEs, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  content,
				"fields": []string{"articleContent"},
				"type":   "best_fields",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("article_es"),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)

	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Error: %s", res.String())
	}

	var r struct {
		Hits struct {
			Total struct {
				Value float64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source models.ArticleEs `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	var articleEsList []models.ArticleEs
	for _, hit := range r.Hits.Hits {
		articleEsList = append(articleEsList, hit.Source)
	}

	fmt.Printf("Total Hits: %d\n", int(r.Hits.Total.Value))
	return articleEsList, nil
}

func DeleteArticle(ids []int) {
	for _, id := range ids {
		// 从 Elasticsearch 的 "article_es" 索引中删除文档
		res, err := esClient.Delete(
			"article_es",     // 索引名称
			strconv.Itoa(id), // 文档 ID
			esClient.Delete.WithContext(context.Background()), // 上下文
		)
		if err != nil {
			log.Fatalf("Error deleting document: %s", err)
		}
		defer res.Body.Close()

		// 打印响应
		log.Println(res)
	}
}
