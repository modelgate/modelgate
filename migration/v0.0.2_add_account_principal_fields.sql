ALTER TABLE `accounts`
  ADD COLUMN `principal_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `id`,
  ADD COLUMN `principal_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' AFTER `principal_type`,
  DROP INDEX `uk_name`,
  DROP COLUMN `nickname`,
  DROP COLUMN `name`;
