// 临时验证脚本：直连 notification.rpc 调 Send / ListByUser
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cli := notification.NewNotificationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send 两条
	for i := 1; i <= 2; i++ {
		resp, err := cli.Send(ctx, &notification.SendReq{
			UserId:  1,
			Title:   fmt.Sprintf("测试通知%d", i),
			Content: "hello from nctest",
			Channel: notification.Channel_CHANNEL_EMAIL,
		})
		fmt.Printf("Send#%d -> resp=%v err=%v\n", i, resp, err)
	}

	// ListByUser
	list, err := cli.ListByUser(ctx, &notification.ListByUserReq{UserId: 1, Page: 1, Size: 10})
	fmt.Printf("ListByUser err=%v total=%d\n", err, list.GetTotal())
	for _, n := range list.GetList() {
		fmt.Printf("  - id=%d title=%s channel=%v isRead=%v\n", n.Id, n.Title, n.Channel, n.IsRead)
	}
}
