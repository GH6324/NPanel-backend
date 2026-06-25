CREATE TABLE IF NOT EXISTS `subscribe_category` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Subscribe Category ID',
    `parent_id` bigint NOT NULL DEFAULT 0 COMMENT 'Parent Category ID',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Category Name',
    `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'Category Description',
    `language` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Language',
    `show` tinyint(1) NOT NULL DEFAULT 1 COMMENT 'Show',
    `sort` int NOT NULL DEFAULT 0 COMMENT 'Sort Order',
    `created_at` datetime(3) DEFAULT NULL COMMENT 'Create Time',
    `updated_at` datetime(3) DEFAULT NULL COMMENT 'Update Time',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_language` (`language`),
    KEY `idx_show_sort` (`show`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Subscribe Categories';

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe' AND COLUMN_NAME = 'category_id');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe` ADD COLUMN `category_id` bigint NOT NULL DEFAULT 0 COMMENT ''Subscribe Category ID'' AFTER `quota`',
    'SELECT ''Column category_id already exists in subscribe table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe' AND INDEX_NAME = 'idx_category_id');
SET @sql = IF(@index_exists = 0,
    'ALTER TABLE `subscribe` ADD INDEX `idx_category_id` (`category_id`)',
    'SELECT ''Index idx_category_id already exists in subscribe table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
