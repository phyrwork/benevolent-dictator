package graph

import (
	"fmt"
	"github.com/phyrwork/benevolent-dictator/pkg/api/graph/model"
	"gorm.io/gorm"
)

func PointersOf[T any](s []T) []*T {
	t := make([]*T, len(s))
	for i := range s {
		t[i] = &s[i]
	}
	return t
}

func MapOf[T, U any](s []T, f func(T) U) []U {
	t := make([]U, len(s))
	for i, v := range s {
		t[i] = f(v)
	}
	return t
}

func MapOfError[T, U any](s []T, f func(T) (U, error)) ([]U, error) {
	t := make([]U, len(s))
	for i, v := range s {
		var err error
		t[i], err = f(v)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func MapPointersOf[T, U any](s []T, f func(T) U) []*U {
	t := MapOf(s, f)
	return PointersOf(t)
}

func MapPointersOfError[T, U any](s []T, f func(T) (U, error)) ([]*U, error) {
	t, err := MapOfError(s, f)
	if err != nil {
		return nil, err
	}
	return PointersOf(t), nil
}

type PageItem interface {
	IDRef() *int
	IDAfter() func(*gorm.DB) *gorm.DB
	IDBeforeOrEqual() func(*gorm.DB) *gorm.DB
}

type PageReader[T PageItem] struct {
	Query     *gorm.DB
	After     T
	Limit     int
	Rows      []T
	PrevCount int64
	NextCount int64
}

func (p *PageReader[T]) StartRow() *T {
	if len(p.Rows) == 0 {
		return nil
	}
	return &p.Rows[0]
}

func (p *PageReader[T]) EndRow() *T {
	if len(p.Rows) == 0 {
		return nil
	}
	return &p.Rows[len(p.Rows)-1]
}

func (p *PageReader[T]) Read() error {
	return p.Query.Transaction(func(tx *gorm.DB) error {
		// Page select.
		qry := tx.Scopes(p.After.IDAfter())
		if p.Limit != 0 {
			qry = qry.Limit(p.Limit)
		}
		if err := qry.Find(&p.Rows).Error; err != nil {
			return fmt.Errorf("page select error: %w", err)
		}
		// Prev./next page identify.
		if err := tx.Scopes(p.After.IDBeforeOrEqual()).Count(&p.PrevCount).Error; err != nil {
			return fmt.Errorf("prev. count error: %w", err)
		}
		if len(p.Rows) == 0 {
			p.NextCount = 0
		} else if err := tx.Scopes((*p.EndRow()).IDAfter()).Count(&p.NextCount).Error; err != nil {
			return fmt.Errorf("next count error: %w", err)
		}
		return nil
	})
}

func (p *PageReader[T]) Info() *model.PageInfo {
	pageInfo := model.PageInfo{
		HasPreviousPage: p.PrevCount > 0,
		HasNextPage:     p.NextCount > 0,
	}
	if len(p.Rows) != 0 {
		pageInfo.StartCursor = (*p.StartRow()).IDRef()
		pageInfo.EndCursor = (*p.EndRow()).IDRef()
	}
	return &pageInfo
}
