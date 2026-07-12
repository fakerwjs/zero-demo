USE `zero_demo`;

CREATE TABLE IF NOT EXISTS `notification` (
  `id`         bigint       NOT NULL AUTO_INCREMENT,
  `user_id`    bigint       NOT NULL DEFAULT '0' COMMENT '接收用户id',
  `title`      varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
  `content`    varchar(1024) NOT NULL DEFAULT '' COMMENT '内容',
  `channel`    tinyint      NOT NULL DEFAULT '1' COMMENT '渠道 1站内信 2短信 3邮件 4推送',
  `is_read`    tinyint      NOT NULL DEFAULT '0' COMMENT '是否已读 0否 1是',
  `create_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='通知表';
