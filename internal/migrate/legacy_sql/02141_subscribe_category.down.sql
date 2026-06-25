ALTER TABLE `subscribe` DROP INDEX IF EXISTS `idx_category_id`;
ALTER TABLE `subscribe` DROP COLUMN IF EXISTS `category_id`;
DROP TABLE IF EXISTS `subscribe_category`;
