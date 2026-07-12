package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ defaultPaymentModel = (*paymentModel)(nil)

type (
	defaultPaymentModel interface {
		Insert(ctx context.Context, data *Payment) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Payment, error)
		FindPageByUser(ctx context.Context, userId int64, page, size int32) ([]*Payment, error)
		FindPageByUserCount(ctx context.Context, userId int64) (int64, error)
		UpdateStatus(ctx context.Context, id int64, status int64, transactionId string, updatedAt int64) error
	}

	paymentModel struct {
		conn  sqlx.SqlConn
		cache cache.CacheConf
	}
)

func newPaymentModel(conn sqlx.SqlConn, c cache.CacheConf) defaultPaymentModel {
	return &paymentModel{conn: conn, cache: c}
}

func (m *paymentModel) Insert(ctx context.Context, data *Payment) (sql.Result, error) {
	query := "INSERT INTO payment (user_id, order_id, order_no, amount, method, status, transaction_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return m.conn.ExecCtx(ctx, query, data.UserId, data.OrderId, data.OrderNo, data.Amount, data.Method, data.Status, data.TransactionId, data.CreatedAt, data.UpdatedAt)
}

func (m *paymentModel) FindOne(ctx context.Context, id int64) (*Payment, error) {
	query := "SELECT id, user_id, order_id, order_no, amount, method, status, transaction_id, created_at, updated_at FROM payment WHERE id = ?"
	var resp Payment
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *paymentModel) FindPageByUser(ctx context.Context, userId int64, page, size int32) ([]*Payment, error) {
	query := "SELECT id, user_id, order_id, order_no, amount, method, status, transaction_id, created_at, updated_at FROM payment WHERE user_id = ? ORDER BY created_at DESC LIMIT ?, ?"
	var resp []*Payment
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, (page-1)*size, size)
	return resp, err
}

func (m *paymentModel) FindPageByUserCount(ctx context.Context, userId int64) (int64, error) {
	query := "SELECT COUNT(*) FROM payment WHERE user_id = ?"
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, userId)
	return count, err
}

func (m *paymentModel) UpdateStatus(ctx context.Context, id int64, status int64, transactionId string, updatedAt int64) error {
	query := "UPDATE payment SET status = ?, transaction_id = ?, updated_at = ? WHERE id = ?"
	_, err := m.conn.ExecCtx(ctx, query, status, transactionId, updatedAt, id)
	return err
}
