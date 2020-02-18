CREATE TABLE `test` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `test` text COMMENT '测试',
  `demo` varchar(11) DEFAULT NULL COMMENT 'DEMO',
  `isdel` tinyint(11) DEFAULT NULL COMMENT '是否删除',
  `created_time` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;