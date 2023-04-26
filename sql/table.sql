-- template.`user` definition

CREATE TABLE `user` (
                        `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键id',
                        `address` varchar(100) NOT NULL COMMENT '用户地址',
                        `balance` varchar(100) DEFAULT NULL COMMENT '用户代币余额',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `user_id_IDX` (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';