package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	OrderModel interface {
		Insert(ctx context.Context, data *Order) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Order, error)
		FindOneByOrderNo(ctx context.Context, orderNo string) (*Order, error)
		FindPageByUser(ctx context.Context, userId int64, page, size int32) ([]*Order, error)
		FindPageByUserCount(ctx context.Context, userId int64) (int64, error)
		UpdateStatus(ctx context.Context, id int64, status int64, updatedAt int64) error
	}

	customOrderModel struct {
		defaultOrderModel
	}

	Order struct {
		Id         int64  `db:"id"`
		UserId     int64  `db:"user_id"`
		OrderNo    string `db:"order_no"`
		TotalPrice int64  `db:"total_price"`
		Status     int64  `db:"status"`
		Items      string `db:"items"`
		CreatedAt  int64  `db:"created_at"`
		UpdatedAt  int64  `db:"updated_at"`
	}
)

func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c),
	}
}
