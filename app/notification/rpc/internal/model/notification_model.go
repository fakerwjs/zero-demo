package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NotificationModel = (*customNotificationModel)(nil)

type (
	// NotificationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotificationModel.
	NotificationModel interface {
		notificationModel
		// ListByUserId 按用户分页查询通知
		ListByUserId(ctx context.Context, userId int64, offset, limit int32) ([]*Notification, error)
		// CountByUserId 统计用户通知总数
		CountByUserId(ctx context.Context, userId int64) (int64, error)
	}

	customNotificationModel struct {
		*defaultNotificationModel
	}
)

// NewNotificationModel returns a model for the database table.
func NewNotificationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NotificationModel {
	return &customNotificationModel{
		defaultNotificationModel: newNotificationModel(conn, c, opts...),
	}
}

func (m *customNotificationModel) ListByUserId(ctx context.Context, userId int64, offset, limit int32) ([]*Notification, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by id desc limit ?, ?", notificationRows, m.table)
	var resp []*Notification
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, offset, limit)
	return resp, err
}

func (m *customNotificationModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	var total int64
	err := m.QueryRowNoCacheCtx(ctx, &total, query, userId)
	return total, err
}
