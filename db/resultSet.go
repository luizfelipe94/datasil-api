package db

import (
	"math"
)

type ResultSet[T any] struct {
	Data      []T `json:"data"`
	Page      int `json:"page"`
	Take      int `json:"take"`
	ItemCount int `json:"itemCount"`
	PageCount int `json:"pageCount"`
}

func NewResultSet[T any](data []T, page int, take int, count int) ResultSet[T] {
	rs := ResultSet[T]{
		Data:      data,
		Page:      page,
		Take:      take,
		ItemCount: count,
	}
	rs.PageCount = int(math.Ceil(float64(rs.ItemCount) / float64(rs.Take)))
	return rs
}
