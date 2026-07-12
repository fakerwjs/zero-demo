CREATE DATABASE IF NOT EXISTS `zero_demo` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `zero_demo`;

CREATE TABLE IF NOT EXISTS `user` (
  `id`         bigint       NOT NULL AUTO_INCREMENT,
  `username`   varchar(64)  NOT NULL DEFAULT '' COMMENT '用户名',
  `password`   varchar(128) NOT NULL DEFAULT '' COMMENT 'bcrypt 密码哈希',
  `mobile`     varchar(20)  NOT NULL DEFAULT '' COMMENT '手机号',
  `nickname`   varchar(64)  NOT NULL DEFAULT '' COMMENT '昵称',
  `create_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_username` (`username`),
  UNIQUE KEY `uniq_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';
