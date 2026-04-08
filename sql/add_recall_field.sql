
-- 添加消息撤回功能的数据库迁移脚本
-- 执行此脚本以更新数据库表结构

USE go_gin_chat;

-- 为 messages 表添加 is_recalled 字段
ALTER TABLE `messages` ADD COLUMN `is_recalled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已撤回 0-否 1-是' AFTER `image_url`;

-- 添加索引（可选，用于提高查询效率）
ALTER TABLE `messages` ADD INDEX `idx_is_recalled` (`is_recalled`);

