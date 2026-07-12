CREATE TABLE IF NOT EXISTS `order` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `order_no` VARCHAR(64) NOT NULL UNIQUE COMMENT '订单编号',
    `total_price` BIGINT NOT NULL DEFAULT 0 COMMENT '总价（分）',
    `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1待支付 2已支付 3已发货 4已完成 5已取消',
    `items` TEXT NOT NULL COMMENT '订单项JSON',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';