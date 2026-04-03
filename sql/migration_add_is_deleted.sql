-- 为 messages 表添加 is_deleted 字段，用于消息撤回功能
ALTER TABLE `messages` ADD COLUMN `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已撤回' AFTER `image_url`;
