package models

type EsSearchWord struct {
	Content string `json:"content"`
}
type HotSearch struct {
	Content string  `json:"content"`
	Score   float64 `json:"score"`
	Type    int     `json:"type"`
}
type HotSearchWords struct {
	Keyword string  `json:"keyword"`
	Score   float64 `json:"score"`
}
