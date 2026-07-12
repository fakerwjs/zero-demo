package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ defaultProductModel = (*productModel)(nil)

type (
	defaultProductModel interface {
		Insert(ctx context.Context, data *Product) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Product, error)
		Update(ctx context.Context, data *Product) error
		Delete(ctx context.Context, id int64) error
		QueryRowsNoCacheCtx(ctx context.Context, v any, sql string, args ...any) error
		QueryRowNoCacheCtx(ctx context.Context, v any, sql string, args ...any) error
		ExecCtx(ctx context.Context, sql string, args ...any) (sql.Result, error)
	}

	productModel struct {
		conn  sqlx.SqlConn
		cache cache.CacheConf
	}
)

func newProductModel(conn sqlx.SqlConn, c cache.CacheConf) defaultProductModel {
	return &productModel{conn: conn, cache: c}
}

func (m *productModel) Insert(ctx context.Context, data *Product) (sql.Result, error) {
	query := "INSERT INTO product (name, description, price, stock, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	return m.conn.ExecCtx(ctx, query, data.Name, data.Description, data.Price, data.Stock, data.Image, data.CreatedAt, data.UpdatedAt)
}

func (m *productModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	query := "SELECT id, name, description, price, stock, image, created_at, updated_at FROM product WHERE id = ?"
	var resp Product
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *productModel) Update(ctx context.Context, data *Product) error {
	query := "UPDATE product SET name = ?, description = ?, price = ?, stock = ?, image = ?, updated_at = ? WHERE id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.Description, data.Price, data.Stock, data.Image, data.UpdatedAt, data.Id)
	return err
}

func (m *productModel) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM product WHERE id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *productModel) QueryRowsNoCacheCtx(ctx context.Context, v any, sql string, args ...any) error {
	return m.conn.QueryRowsCtx(ctx, v, sql, args...)
}

func (m *productModel) QueryRowNoCacheCtx(ctx context.Context, v any, sql string, args ...any) error {
	return m.conn.QueryRowCtx(ctx, v, sql, args...)
}

func (m *productModel) ExecCtx(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	return m.conn.ExecCtx(ctx, sql, args...)
}
