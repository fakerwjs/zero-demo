package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ defaultOrderModel = (*orderModel)(nil)

type (
	defaultOrderModel interface {
		Insert(ctx context.Context, data *Order) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Order, error)
		FindOneByOrderNo(ctx context.Context, orderNo string) (*Order, error)
		FindPageByUser(ctx context.Context, userId int64, page, size int32) ([]*Order, error)
		FindPageByUserCount(ctx context.Context, userId int64) (int64, error)
		UpdateStatus(ctx context.Context, id int64, status int64, updatedAt int64) error
	}

	orderModel struct {
		conn  sqlx.SqlConn
		cache cache.CacheConf
	}
)

func newOrderModel(conn sqlx.SqlConn, c cache.CacheConf) defaultOrderModel {
	return &orderModel{conn: conn, cache: c}
}

func (m *orderModel) Insert(ctx context.Context, data *Order) (sql.Result, error) {
	query := "INSERT INTO `order` (user_id, order_no, total_price, status, items, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	return m.conn.ExecCtx(ctx, query, data.UserId, data.OrderNo, data.TotalPrice, data.Status, data.Items, data.CreatedAt, data.UpdatedAt)
}

func (m *orderModel) FindOne(ctx context.Context, id int64) (*Order, error) {
	query := "SELECT id, user_id, order_no, total_price, status, items, created_at, updated_at FROM `order` WHERE id = ?"
	var resp Order
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *orderModel) FindOneByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	query := "SELECT id, user_id, order_no, total_price, status, items, created_at, updated_at FROM `order` WHERE order_no = ?"
	var resp Order
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderNo)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *orderModel) FindPageByUser(ctx context.Context, userId int64, page, size int32) ([]*Order, error) {
	query := "SELECT id, user_id, order_no, total_price, status, items, created_at, updated_at FROM `order` WHERE user_id = ? ORDER BY created_at DESC LIMIT ?, ?"
	var resp []*Order
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, (page-1)*size, size)
	return resp, err
}

func (m *orderModel) FindPageByUserCount(ctx context.Context, userId int64) (int64, error) {
	query := "SELECT COUNT(*) FROM `order` WHERE user_id = ?"
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, userId)
	return count, err
}

func (m *orderModel) UpdateStatus(ctx context.Context, id int64, status int64, updatedAt int64) error {
	query := "UPDATE `order` SET status = ?, updated_at = ? WHERE id = ?"
	_, err := m.conn.ExecCtx(ctx, query, status, updatedAt, id)
	return err
}
