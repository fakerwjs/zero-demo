CREATE TABLE IF NOT EXISTS `payment` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '支付ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `order_id` BIGINT NOT NULL COMMENT '订单ID',
    `order_no` VARCHAR(64) NOT NULL COMMENT '订单编号',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '支付金额（分）',
    `method` INT NOT NULL DEFAULT 1 COMMENT '支付方式：1支付宝 2微信 3银行卡',
    `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1待支付 2支付成功 3支付失败 4已退款',
    `transaction_id` VARCHAR(128) DEFAULT NULL COMMENT '第三方交易号',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_order_no` (`order_no`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付表';