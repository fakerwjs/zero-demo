CREATE TABLE IF NOT EXISTS `product` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `name` VARCHAR(255) NOT NULL COMMENT '商品名称',
    `description` TEXT DEFAULT NULL COMMENT '商品描述',
    `price` BIGINT NOT NULL DEFAULT 0 COMMENT '价格（分）',
    `stock` BIGINT NOT NULL DEFAULT 0 COMMENT '库存',
    `image` VARCHAR(500) DEFAULT NULL COMMENT '商品图片URL',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';