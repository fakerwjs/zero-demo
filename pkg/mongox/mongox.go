// Package mongox 封装 MongoDB 客户端，用于存储非结构化/日志类文档。
package mongox

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationLog 通知审计日志文档。
type NotificationLog struct {
	NotificationID int64     `bson:"notificationId"`
	UserID         int64     `bson:"userId"`
	Title          string    `bson:"title"`
	Content        string    `bson:"content"`
	Channel        int32     `bson:"channel"`
	CreatedAt      time.Time `bson:"createdAt"`
}

// Client 持有 MongoDB 连接与目标集合。
type Client struct {
	cli  *mongo.Client
	coll *mongo.Collection
}

// New 连接 MongoDB。uri 形如 mongodb://root:pass@127.0.0.1:27017。
func New(uri, database, collection string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := cli.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &Client{
		cli:  cli,
		coll: cli.Database(database).Collection(collection),
	}, nil
}

// InsertNotificationLog 写入一条通知审计日志。
func (c *Client) InsertNotificationLog(ctx context.Context, log *NotificationLog) error {
	_, err := c.coll.InsertOne(ctx, log)
	return err
}

// CountByUser 统计某用户的审计日志数（供验证用）。
func (c *Client) CountByUser(ctx context.Context, userID int64) (int64, error) {
	return c.coll.CountDocuments(ctx, map[string]any{"userId": userID})
}

func (c *Client) Close(ctx context.Context) error {
	return c.cli.Disconnect(ctx)
}
