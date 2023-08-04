package parameter

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/errors"
	"gorm.io/gorm"
)

type Parameter[T any] struct {
	Page, Limit                               int
	ctx                                       *app.HttpContext
	searchTerm, sortOrder                     string
	selectColumns, searchColumns, sortColumns []string
	processEachItem                           bool
	processEachItemFn                         func(*T)
}

func New[T any](ctx *app.HttpContext) *Parameter[T] {
	p := new(Parameter[T])
	p.ctx = ctx
	p.processEachItem = false

	p.loadParams()

	return p
}

func (p *Parameter[_]) loadParams() {
	pageParam, ok := p.ctx.GetParam("page")
	if !ok {
		pageParam = "1"
	}
	page, _ := strconv.Atoi(pageParam)
	p.Page = page

	limitParam, ok := p.ctx.GetParam("limit")
	if !ok {
		limitParam = "25"
	}
	limit, _ := strconv.Atoi(limitParam)
	p.Limit = limit

	if searchTerm, ok := p.ctx.GetParam("searchTerm"); ok {
		p.searchTerm = strings.TrimSpace(searchTerm)
	}
}

func (p *Parameter[_]) DBLikeSearch() (string, bool) {
	if strings.TrimSpace(p.searchTerm) != "" {
		return "%" + p.searchTerm + "%", true
	}
	return "", false
}

func (p *Parameter[T]) Paginate(db *gorm.DB) *gorm.DB {
	return db.Offset(p.Limit * (p.Page - 1)).Limit(p.Limit)
}

func (p *Parameter[T]) SelectColumns(columns ...string) *Parameter[T] {
	p.selectColumns = columns
	return p
}

func (p *Parameter[T]) WithSortAscending(columns ...string) *Parameter[T] {
	p.sortColumns = columns
	p.sortOrder = "ASC"
	return p
}

func (p *Parameter[T]) WithSortDescending(columns ...string) *Parameter[T] {
	p.sortColumns = columns
	p.sortOrder = "DESC"
	return p
}

func (p *Parameter[T]) SearchColumns(columns ...string) *Parameter[T] {
	p.searchColumns = columns

	return p
}

func (p *Parameter[T]) WithEachItemProcess(fn func(*T)) *Parameter[T] {
	p.processEachItemFn = fn
	p.processEachItem = true

	return p
}

func (p *Parameter[T]) Execute(db *gorm.DB) (*ListResponse[T], error) {
	var (
		totalRecords int64
		result       []T
		wg           sync.WaitGroup
	)

	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		if err := db.WithContext(context.Background()).Count(&totalRecords).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		qry := db.WithContext(context.Background())

		if searchTerm, ok := p.DBLikeSearch(); ok {
			for i, column := range p.searchColumns {
				if i == 0 {
					qry = qry.Debug().Where(column+" LIKE ?", searchTerm)
				} else {
					qry = qry.Or(column+" LIKE ?", searchTerm)
				}
			}
		}

		if len(p.selectColumns) > 0 {
			fields := strings.Join(p.selectColumns, ",")
			qry = qry.Select(fields)
		}

		if len(p.sortColumns) > 0 {
			orderClause := strings.Join(p.sortColumns, ",") + " " + p.sortOrder
			qry = qry.Order(orderClause)
		}

		qry = p.Paginate(qry)
		if err := qry.Find(&result).Error; err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return nil, errors.NewBadRequestError(consts.BadRequest, err)
	default:
	}

	// process each item
	if p.processEachItem {
		for i, item := range result {
			p.processEachItemFn(&item)
			result[i] = item
		}
	}

	return NewListResponse[T](p.Page, p.Limit, totalRecords, result), nil
}

type ListResponse[T any] struct {
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	Limit    int64 `json:"limit"`
	LastPage int64 `json:"last_page"`
	From     int64 `json:"from"`
	To       int64 `json:"to"`
	Data     []T   `json:"data"`
}

func NewListResponse[T any](page, limit int, total int64, data []T) *ListResponse[T] {
	response := new(ListResponse[T])
	response.Page = int64(page)
	response.Limit = int64(limit)
	response.From = ((response.Page - 1) * response.Limit) + 1
	response.To = response.From + response.Limit - 1
	response.Total = total
	response.Data = data

	// calculate last page
	lp := float64(total) / float64(limit)
	lastPage := int64(lp)
	if lp > float64(lastPage) {
		lastPage++
	}
	response.LastPage = lastPage

	return response

}
