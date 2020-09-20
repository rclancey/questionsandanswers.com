package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultSortOrder = SortAskDateDesc
	DefaultPageSize = 10
	DefaultPageNum = 0
	MaxPageSize = 100
	MaxPayloadSize = 0x10000 // 64kb
)

type ListMeta struct {
	PageNum int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	SortOrder string `json:"sortOrder"`
	ResultCount int `json:"resultCount"`
	TotalResults int `json:"totalResults"`
	TotalPages int `json:"totalPages"`
	PrevPage string `json:"prevPage,omitempty"`
	ThisPage string `json:"thisPage"`
	NextPage string `json:"nextPage,omitempty"`
}

type ListResponse struct {
	Meta ListMeta `json:"meta"`
	Questions []*Question `json:"results"`
}

func ListQuestions(w http.ResponseWriter, r *http.Request, pathinfo string) (interface{}, error) {
	q := r.URL.Query()
	sortOrder := q.Get("sort")
	if sortOrder == "" {
		sortOrder = DefaultSortOrder
	}
	pageSize, err := strconv.Atoi(q.Get("count"))
	if err != nil || pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	pageNum, err := strconv.Atoi(q.Get("page"))
	if err != nil || pageNum < 0 {
		pageNum = DefaultPageNum
	}
	count, err := db.CountQuestions()
	if err != nil {
		return nil, err
	}
	questions, err := db.ListQuestions(sortOrder, pageSize, pageNum)
	if err != nil {
		return nil, err
	}
	totalPages := count / pageSize
	if totalPages % pageSize > 0 {
		totalPages += 1
	}
	query := url.Values{}
	query.Set("sort", sortOrder)
	query.Set("count", strconv.Itoa(pageSize))
	query.Set("page", strconv.Itoa(pageNum))
	thisPage := &url.URL{Path: r.URL.Path, RawQuery: query.Encode()}
	meta := ListMeta{
		PageNum: pageNum,
		PageSize: pageSize,
		SortOrder: sortOrder,
		ResultCount: len(questions),
		TotalResults: count,
		TotalPages: totalPages,
		ThisPage: thisPage.String(),
	}
	if pageNum > 0 {
		query.Set("page", strconv.Itoa(pageNum - 1))
		prevPage := &url.URL{Path: r.URL.Path, RawQuery: query.Encode()}
		meta.PrevPage = prevPage.String()
	}
	if pageNum < totalPages - 1 {
		query.Set("page", strconv.Itoa(pageNum + 1))
		nextPage := &url.URL{Path: r.URL.Path, RawQuery: query.Encode()}
		meta.NextPage = nextPage.String()
	}
	return &ListResponse{meta, questions}, nil
}

func CreateQuestion(w http.ResponseWriter, r *http.Request, pathinfo string) (interface{}, error) {
	data, err := readAtMost(r, MaxPayloadSize)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, BadRequest(nil)
	}
	q, err := db.CreateQuestion(string(data))
	if err != nil {
		return nil, err
	}
	return q, nil
}

func GetQuestion(w http.ResponseWriter, r *http.Request, pathinfo string) (interface{}, error) {
	questionId := strings.Trim(pathinfo, "/")
	if questionId == "" {
		return nil, NotFound(nil)
	}
	q, err := db.ReadQuestion(questionId)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, NotFound(nil)
	}
	return q, nil
}

func AnswerQuestion(w http.ResponseWriter, r *http.Request, pathinfo string) (interface{}, error) {
	questionId := strings.Trim(pathinfo, "/")
	if questionId == "" {
		return nil, NotFound(nil)
	}
	q, err := db.ReadQuestion(questionId)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, NotFound(nil)
	}
	data, err := readAtMost(r, MaxPayloadSize)
	if err != nil {
		return nil, err
	}
	err = q.SetAnswer(string(data))
	if err != nil {
		return nil, err
	}
	err = db.UpdateQuestion(q)
	if err == ErrNoSuchQuestion {
		return nil, NotFound(nil)
	}
	if err != nil {
		return nil, err
	}
	return q, nil
}
