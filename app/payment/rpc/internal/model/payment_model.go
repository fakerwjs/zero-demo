package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentModel = (*customPaymentModel)(nil)

type (
	PaymentModel interface {
		Insert(ctx context.Context, data *Payment) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Payment, error)
		FindPageByUser(ctx context.Context, userId int64, page, size int32) ([]*Payment, error)
		FindPageByUserCount(ctx context.Context, userId int64) (int64, error)
		UpdateStatus(ctx context.Context, id int64, status int64, transactionId string, updatedAt int64) error
	}

	customPaymentModel struct {
		defaultPaymentModel
	}

	Payment struct {
		Id            int64  `db:"id"`
		UserId        int64  `db:"user_id"`
		OrderId       int64  `db:"order_id"`
		OrderNo       string `db:"order_no"`
		Amount        int64  `db:"amount"`
		Method        int64  `db:"method"`
		Status        int64  `db:"status"`
		TransactionId string `db:"transaction_id"`
		CreatedAt     int64  `db:"created_at"`
		UpdatedAt     int64  `db:"updated_at"`
	}
)

func NewPaymentModel(conn sqlx.SqlConn, c cache.CacheConf) PaymentModel {
	return &customPaymentModel{
		defaultPaymentModel: newPaymentModel(conn, c),
	}
}
