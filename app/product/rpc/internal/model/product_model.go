package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	ProductModel interface {
		Insert(ctx context.Context, data *Product) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Product, error)
		FindPageList(ctx context.Context, page, size int32) ([]*Product, error)
		FindPageListCount(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *Product) error
		Delete(ctx context.Context, id int64) error
		DeductStock(ctx context.Context, id, amount int64) (int64, error)
	}

	customProductModel struct {
		defaultProductModel
	}

	Product struct {
		Id          int64  `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		Price       int64  `db:"price"`
		Stock       int64  `db:"stock"`
		Image       string `db:"image"`
		CreatedAt   int64  `db:"created_at"`
		UpdatedAt   int64  `db:"updated_at"`
	}
)

func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf) ProductModel {
	return &customProductModel{
		defaultProductModel: newProductModel(conn, c),
	}
}

func (m *customProductModel) FindPageList(ctx context.Context, page, size int32) ([]*Product, error) {
	query := "SELECT id, name, description, price, stock, image, created_at, updated_at FROM product ORDER BY created_at DESC LIMIT ?, ?"
	var resp []*Product
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, (page-1)*size, size)
	return resp, err
}

func (m *customProductModel) FindPageListCount(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM product"
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query)
	return count, err
}

func (m *customProductModel) DeductStock(ctx context.Context, id, amount int64) (int64, error) {
	query := "UPDATE product SET stock = stock - ?, updated_at = ? WHERE id = ? AND stock >= ?"
	result, err := m.ExecCtx(ctx, query, amount, time.Now().Unix(), id, amount)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
