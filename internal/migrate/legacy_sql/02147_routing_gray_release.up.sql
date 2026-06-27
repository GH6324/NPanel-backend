CREATE TABLE IF NOT EXISTS `routing_gray_release` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Routing gray release ID',
  `profile_code` varchar(128) NOT NULL COMMENT 'Routing profile code',
  `name` varchar(255) NOT NULL COMMENT 'Gray release name',
  `status` varchar(32) NOT NULL DEFAULT 'draft' COMMENT 'draft/running/paused/completed/rolled_back',
  `batch_no` bigint NOT NULL DEFAULT 0 COMMENT 'Current gray batch number',
  `target_type` varchar(32) NOT NULL DEFAULT 'user' COMMENT 'user/user_subscribe/subscribe/node',
  `target_ids_json` text NOT NULL COMMENT 'Target IDs JSON array',
  `operator` varchar(128) NOT NULL DEFAULT '' COMMENT 'Operator',
  `rollback_reason` text COMMENT 'Rollback reason',
  `started_at` datetime NULL COMMENT 'Started at',
  `ended_at` datetime NULL COMMENT 'Ended at',
  `release_json` text NOT NULL COMMENT 'Release metadata JSON',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
  PRIMARY KEY (`id`),
  KEY `idx_routing_gray_release_profile` (`profile_code`, `status`),
  KEY `idx_routing_gray_release_updated` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Routing gray release batches';
