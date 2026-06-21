-- 为messages表添加is_deleted字段，用于实现消息撤回功能
ALTER TABLE `messages` ADD COLUMN `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已撤回：0-未撤回，1-已撤回' AFTER `image_url`;
