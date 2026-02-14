DROP TABLE IF EXISTS `account_api_keys`;

CREATE TABLE `account_api_keys` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `account_id` bigint unsigned NOT NULL DEFAULT '0',
  `key_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `key_prefix` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `key_suffix` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `key_hash` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` enum('enabled','disabled','revoked') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `scope` json DEFAULT NULL,
  `quote_limit` bigint unsigned DEFAULT NULL,
  `quote_used` bigint unsigned NOT NULL DEFAULT '0',
  `rate_limit` int unsigned DEFAULT NULL,
  `last_used_at` datetime(3) DEFAULT NULL,
  `expired_at` datetime DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key_hash` (`key_hash`),
  KEY `idx_key_prefix` (`key_prefix`),
  KEY `idx_account` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `accounts`;

CREATE TABLE `accounts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `balance` bigint unsigned NOT NULL DEFAULT '0',
  `status` enum('enabled','disabled') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=10003 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `ledgers`;

CREATE TABLE `ledgers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `account_id` bigint unsigned NOT NULL DEFAULT '0',
  `type` enum('consume','refund','charge','adjust') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'consume',
  `amount` bigint NOT NULL DEFAULT '0',
  `balance_after` bigint unsigned NOT NULL DEFAULT '0',
  `request_id` bigint unsigned NOT NULL DEFAULT '0',
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_account_type_time` (`created_at`,`account_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=103753 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `menus`;

CREATE TABLE `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `pid` bigint unsigned NOT NULL DEFAULT '0',
  `name` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `type` tinyint unsigned NOT NULL DEFAULT '1',
  `route_name` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `route_path` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `component` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `i18n_key` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `order` bigint unsigned NOT NULL DEFAULT '0',
  `icon_type` tinyint unsigned NOT NULL DEFAULT '1',
  `icon` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` enum('enabled','disabled') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `keep_alive` tinyint unsigned NOT NULL DEFAULT '0',
  `constant` tinyint unsigned NOT NULL DEFAULT '0',
  `href` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `hide_in_menu` tinyint unsigned NOT NULL DEFAULT '0',
  `active_menu` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `multi_tab` tinyint unsigned NOT NULL DEFAULT '0',
  `fixed_index_in_tab` bigint unsigned NOT NULL DEFAULT '0',
  `query` json DEFAULT NULL,
  `buttons` json DEFAULT NULL,
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_router_name` (`route_name`)
) ENGINE=InnoDB AUTO_INCREMENT=10037 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



DROP TABLE IF EXISTS `model_pricings`;

CREATE TABLE `model_pricings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `provider_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `model_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `currency` enum('USD','CNY','POINT') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'USD',
  `points_per_currency` bigint unsigned NOT NULL DEFAULT '1',
  `token_num` bigint unsigned NOT NULL DEFAULT '1000000',
  `input_price` decimal(10,4) unsigned NOT NULL DEFAULT '0.0000',
  `input_cache_price` decimal(10,4) unsigned NOT NULL DEFAULT '0.0000',
  `output_price` decimal(10,4) unsigned NOT NULL DEFAULT '0.0000',
  `status` enum('enabled','disabled') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `effective_from` datetime NOT NULL,
  `effective_to` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_provider_model_effective` (`provider_code`,`model_code`,`effective_from`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `models`;

CREATE TABLE `models` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `provider_id` bigint unsigned NOT NULL DEFAULT '0',
  `provider_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `actual_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `priority` bigint NOT NULL DEFAULT '1',
  `weight` bigint NOT NULL DEFAULT '100',
  `status` enum('enabled','disabled','deprecated') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_provider_model` (`provider_code`,`code`),
  KEY `idx_name_status` (`name`(32))
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `permissions`;

CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `code` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `data` json DEFAULT NULL,
  `desc` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `provider_api_keys`;

CREATE TABLE `provider_api_keys` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `provider_id` bigint NOT NULL DEFAULT '0',
  `provider_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `key_prefix` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `key_suffix` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `key_encrypted` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `quote_used` bigint unsigned NOT NULL DEFAULT '0',
  `quote_limit` bigint unsigned DEFAULT NULL,
  `rate_limit` int unsigned DEFAULT NULL,
  `last_used_at` datetime(3) DEFAULT NULL,
  `weight` bigint NOT NULL DEFAULT '100',
  `status` enum('enabled','disabled','cooldown') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `faild_count` bigint NOT NULL DEFAULT '0',
  `cool_down_until` datetime DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_provider_status` (`provider_id`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `providers`;

CREATE TABLE `providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `base_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` enum('enabled','disabled') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_name` (`name`(32))
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `refresh_tokens`;

CREATE TABLE `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增 id',
  `user_id` bigint NOT NULL DEFAULT '0',
  `jti` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `expires_at` datetime NOT NULL,
  `description` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk-user-jti` (`user_id`,`jti`)
) ENGINE=InnoDB AUTO_INCREMENT=10245 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `relay_hourly_usages`;

CREATE TABLE `relay_hourly_usages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `time` datetime NOT NULL,
  `provider_code` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `total_request` bigint NOT NULL DEFAULT '0',
  `total_success` bigint NOT NULL DEFAULT '0',
  `total_failed` bigint NOT NULL DEFAULT '0',
  `total_point` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_time_provider` (`time`,`provider_code`)
) ENGINE=InnoDB AUTO_INCREMENT=533 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `relay_usages`;

CREATE TABLE `relay_usages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `total_request` bigint NOT NULL DEFAULT '0',
  `total_success` bigint NOT NULL DEFAULT '0',
  `total_failed` bigint NOT NULL DEFAULT '0',
  `total_point` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `request_attempts`;

CREATE TABLE `request_attempts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_uuid` binary(16) NOT NULL,
  `attempt_no` tinyint unsigned NOT NULL DEFAULT '0',
  `actual_model` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `account_id` bigint unsigned NOT NULL DEFAULT '0',
  `account_api_key_id` bigint unsigned NOT NULL DEFAULT '0',
  `provider_id` bigint unsigned NOT NULL DEFAULT '0',
  `provider_api_key_id` bigint unsigned NOT NULL DEFAULT '0',
  `model_id` bigint unsigned NOT NULL DEFAULT '0',
  `prompt_tokens` bigint unsigned NOT NULL DEFAULT '0',
  `completion_tokens` bigint unsigned NOT NULL DEFAULT '0',
  `total_tokens` bigint unsigned NOT NULL DEFAULT '0',
  `status` enum('pending','success','failed','cancelled') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'pending',
  `completed_at` datetime(3) DEFAULT NULL,
  `error_code` int unsigned NOT NULL DEFAULT '0',
  `error_message` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_request_uuid_attempt_no` (`request_uuid`,`attempt_no`)
) ENGINE=InnoDB AUTO_INCREMENT=103750 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `requests`;

CREATE TABLE `requests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_uuid` binary(16) NOT NULL,
  `actual_model` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `account_id` bigint unsigned NOT NULL DEFAULT '0',
  `account_api_key_id` bigint unsigned NOT NULL DEFAULT '0',
  `provider_id` bigint unsigned NOT NULL DEFAULT '0',
  `provider_api_key_id` bigint unsigned NOT NULL DEFAULT '0',
  `model_id` bigint unsigned NOT NULL DEFAULT '0',
  `prompt_tokens` bigint unsigned NOT NULL DEFAULT '0',
  `completion_tokens` bigint unsigned NOT NULL DEFAULT '0',
  `total_tokens` bigint unsigned NOT NULL DEFAULT '0',
  `status` enum('pending','success','failed','cancelled') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'pending',
  `completed_at` datetime(3) DEFAULT NULL,
  `error_code` int unsigned NOT NULL DEFAULT '0',
  `error_message` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_request_uuid` (`request_uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=103687 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `code` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT '0',
  `permission` json DEFAULT NULL,
  `description` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `status` enum('enabled','disabled') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=10001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增 id',
  `username` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `nickname` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `email` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `phone` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `gender` enum('male','female','unknown') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'unknown',
  `roles` varchar(2000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `password` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `avatar_url` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `description` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` enum('enabled','disabled') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'enabled',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=10001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


SET NAMES utf8mb4;

LOCK TABLES `menus` WRITE;

INSERT INTO `menus` (`id`, `pid`, `name`, `type`, `route_name`, `route_path`, `component`, `i18n_key`, `order`, `icon_type`, `icon`, `status`, `keep_alive`, `constant`, `href`, `hide_in_menu`, `active_menu`, `multi_tab`, `fixed_index_in_tab`, `query`, `buttons`, `created_at`, `updated_at`)
VALUES
  (10009,0,'首页',1,'home','/home','layout.base$view.home','route.home',0,1,'mdi:monitor-dashboard','enabled',0,0,'',0,'',0,0,'[]','[]','2025-05-24 19:43:09.000','2026-01-11 17:52:19.000'),
  (10011,0,'系统管理',1,'manage','/manage','layout.base','route.manage',5,1,'carbon:cloud-service-management','enabled',0,0,'',0,'',0,0,'[]','[]','2025-05-25 00:37:23.000','2025-05-25 00:37:23.000'),
  (10012,10011,'用户管理',2,'manage_user','/manage/user','view.manage_user','route.manage_user',1,1,'ic:round-manage-accounts','enabled',0,0,'',0,'',0,0,'[]','[]','2025-05-25 00:38:06.000','2025-05-25 00:38:06.000'),
  (10013,10011,'角色管理',2,'manage_role','/manage/role','view.manage_role','route.manage_role',2,1,'carbon:user-role','enabled',0,0,'',0,'',0,0,'[]','[]','2025-05-25 00:38:53.000','2025-05-25 00:44:17.000'),
  (10014,10011,'菜单管理  ',2,'manage_menu','/manage/menu','view.manage_menu','route.manage_menu',3,1,'material-symbols:route','enabled',0,0,'',0,'',0,0,'[]','[]','2025-05-25 00:39:27.000','2025-05-25 00:40:49.000'),
  (10015,0,'登录',2,'login','/login/:module(pwd-login|code-login|register|reset-pwd|bind-wechat)?','layout.blank$view.login','route.login',0,2,'activity','enabled',0,1,'',1,'',0,0,'[]','[]','2025-05-25 16:09:31.000','2025-05-25 16:09:31.000'),
  (10022,0,'模型管理',1,'relay','/relay','layout.base','route.relay',0,1,'icon-park-solid:layers','enabled',0,0,'',0,'',0,0,'[]','[]','2025-12-23 22:31:57.000','2026-01-08 01:13:39.000'),
  (10023,10022,'供应商',2,'relay_provider','/relay/provider','view.relay_provider','route.relay_provider',1,1,'icon-park-solid:peoples-two','enabled',0,0,'',0,'',0,0,'[]','[]','2025-12-23 22:32:35.000','2026-01-08 01:29:10.000'),
  (10025,10022,'模型价格',2,'relay_model-pricing','/relay/model-pricing','view.relay_model-pricing','route.relay_model-pricing',4,1,'icon-park-solid:financing-two','enabled',0,0,'',0,'',0,0,'[]','[]','2025-12-23 22:47:12.000','2026-01-06 03:31:54.000'),
  (10026,10022,'模型 APIkey',2,'relay_provider-api-key','/relay/provider-api-key','view.relay_provider-api-key','route.relay_provider-api-key',3,1,'icon-park-solid:open-one','enabled',0,0,'',0,'',0,0,'[]','[]','2025-12-23 22:51:29.000','2026-01-06 03:31:47.000'),
  (10027,10022,'模型列表',2,'relay_model','/relay/model','view.relay_model','route.relay_model',2,1,'icon-park-solid:robot-one','enabled',0,0,'',0,'',0,0,'[]','[]','2025-12-23 22:54:21.000','2026-01-08 01:16:58.000'),
  (10028,0,'用户管理',1,'user','/user','layout.base','route.user',0,1,'icon-park-solid:user','enabled',0,0,'',0,'',0,0,'[]','[]','2026-01-10 14:25:27.000','2026-01-10 14:30:33.000'),
  (10031,10028,'账号管理',2,'user_account','/user/account','view.user_account','route.user_account',0,1,'icon-park-solid:id-card-h','enabled',0,0,'',0,'',0,0,'[]','[]','2026-01-11 14:58:45.000','2026-01-11 15:06:40.000'),
  (10032,10028,'用户APIKey',2,'user_api-key','/user/api-key','view.user_api-key','route.user_api-key',0,1,'icon-park-solid:open-one','enabled',0,0,'',0,'',0,0,'[{\"key\": \"key12\", \"value\": \"value12\"}]','[{\"code\": \"test12\", \"desc\": \"1112\"}, {\"code\": \"test2\", \"desc\": \"222\"}]','2026-01-11 14:59:32.000','2026-01-11 22:19:39.000'),
  (10033,0,'用量管理',1,'usage','/usage','layout.base','route.usage',0,1,'icon-park-solid:workbench','enabled',0,0,'',0,'',0,0,'[]','[]','2026-01-14 00:57:27.000','2026-01-14 23:05:19.000'),
  (10034,10033,'请求管理',2,'usage_request','/usage/request','view.usage_request','route.usage_request',0,1,'icon-park-solid:trend','enabled',0,0,'',0,'',0,0,'[]','[]','2026-01-14 01:04:11.000','2026-01-14 23:05:49.000'),
  (10035,10033,'账单管理',2,'usage_ledger','/usage/ledger','view.usage_ledger','route.usage_ledger',0,1,'icon-park-solid:financing-two','enabled',0,0,'',0,'',0,0,'[]','[]','2026-01-14 22:47:12.000','2026-01-14 22:49:48.000'),
  (10036,10011,'接口权限',2,'manage_permission','/manage/permission','view.manage_permission','route.manage_permission',4,1,'icon-park-solid:tree-diagram','enabled',0,0,'',0,'',0,0,'[]','[{\"code\": \"system:permission:list\", \"desc\": \"列表\"}, {\"code\": \"system:permission:edit\", \"desc\": \"编辑\"}, {\"code\": \"system:permission:delete\", \"desc\": \"删除\"}]','2026-01-25 00:16:48.082','2026-01-26 02:35:06.896');

UNLOCK TABLES;

LOCK TABLES `model_pricings` WRITE;

INSERT INTO `model_pricings` (`id`, `provider_code`, `model_code`, `currency`, `points_per_currency`, `token_num`, `input_price`, `input_cache_price`, `output_price`, `status`, `effective_from`, `effective_to`, `created_at`)
VALUES
  (1,'openai','gpt-4o','USD',7200000,1000000,2.5000,1.2500,10.0000,'enabled','2025-12-01 00:00:00','2026-12-31 00:00:00','2025-12-20 15:38:38'),
  (2,'openai','gpt-5.2','USD',7200000,1000000,1.7500,0.1750,14.0000,'enabled','2025-12-01 00:00:00','2026-12-31 00:00:00','2025-12-20 15:47:45'),
  (3,'openai','gpt-5.1-codex','USD',7200000,1000000,1.2500,0.1250,10.0000,'enabled','2025-12-01 00:00:00','2026-12-31 00:00:00','2025-12-20 15:49:32'),
  (5,'openai','gpt-5-codex','USD',7200000,1000000,1.2500,0.1250,10.0000,'enabled','2025-12-01 00:00:00','2026-12-31 00:00:00','2025-12-20 15:50:07'),
  (6,'openai','gpt-5-codex-max','USD',7200000,1000000,1.2500,0.1250,10.0000,'enabled','2025-12-15 00:00:00','2026-12-31 00:00:00','2025-12-20 15:50:43'),
  (7,'deepseek','deepseek-chat','CNY',1000000,1000000,2.0000,0.2000,3.0000,'enabled','2025-12-15 00:00:00','2026-12-31 00:00:00','2025-12-22 01:40:42'),
  (8,'openai','text-embedding-3-small','USD',7200000,1000000,0.0200,0.5000,8.0000,'enabled','2025-12-15 00:00:00','2026-12-31 00:00:00','2025-12-22 11:42:47'),
  (9,'openai','gpt-4.1','USD',7200000,1000000,2.0000,0.0000,0.0000,'enabled','2025-12-15 00:00:00','2026-12-31 00:00:00','2025-12-22 15:44:30'),
  (10,'zhipu','glm-4.7','CNY',1000000,1000000,2.0000,0.4000,8.0000,'enabled','2025-12-30 00:00:00','2026-12-31 00:00:00','2025-12-30 13:41:16'),
  (11,'anthropic','claude-sonnet-4-5-20250929','USD',7200000,1000000,3.0000,0.0000,15.0000,'enabled','2025-12-30 00:00:00','2030-12-31 00:00:00','2025-12-30 18:02:21'),
  (12,'anthropic','claude-sonnet-4','USD',7200000,1000000,3.0000,0.0000,15.0000,'enabled','2025-12-30 00:00:00','2026-12-31 00:00:00','2025-12-30 18:02:21'),
  (15,'anthropic','claude-haiku-4-5-20251001','USD',7200000,1000000,1.0000,0.0000,5.0000,'enabled','2026-01-01 23:37:48','2030-01-01 23:37:56','2026-01-13 23:37:59'),
  (16,'zhipu-claude','glm-4.7','CNY',1000000,1000000,2.0000,0.4000,8.0000,'enabled','2026-01-01 16:32:15','2030-01-31 16:32:29','2026-01-14 16:32:31'),
  (17,'zhipu-claude','glm-4.5-air','CNY',1000000,1000000,2.0000,0.4000,8.0000,'enabled','2026-01-01 17:43:30','2030-01-31 17:43:36','2026-01-14 17:43:38');

UNLOCK TABLES;


LOCK TABLES `models` WRITE;

INSERT INTO `models` (`id`, `provider_id`, `provider_code`, `code`, `actual_code`, `name`, `priority`, `weight`, `status`, `created_at`, `updated_at`)
VALUES
  (1,1,'openai','gpt-4o','','GPT-4o',1,100,'enabled','2025-12-20 14:39:28.000','2025-12-20 14:39:28.000'),
  (2,2,'deepseek','deepseek-chat','','DeepSeek-Chat',1,100,'enabled','2025-12-20 14:40:10.000','2025-12-20 14:40:10.000'),
  (3,1,'openai','text-embedding-3-small','','text-embedding-3-small',1,100,'enabled','2025-12-22 11:41:12.000','2025-12-22 11:41:12.000'),
  (4,1,'openai','gpt-4.1','','gpt-4.1',1,100,'enabled','2025-12-22 15:43:45.000','2025-12-22 15:43:45.000'),
  (5,3,'zhipu','glm-4.7','','GLM-4.7',1,100,'enabled','2025-12-30 13:00:00.000','2026-01-14 02:24:12.000'),
  (7,6,'anthropic','claude-sonnet-4-5-20250929','','claude sonnet 4.5',1,100,'enabled','2025-12-30 13:00:00.000','2026-01-13 23:23:46.000'),
  (9,6,'anthropic','claude-sonnet-4','','claude sonnet 4',1,100,'enabled','2025-12-30 13:00:00.000','2026-01-09 23:52:55.000'),
  (10,6,'anthropic','claude-haiku-4-5-20251001','','claude-haiku-4-5',1,100,'enabled','2026-01-13 23:36:20.000','2026-01-13 23:36:20.000'),
  (13,9,'zhipu-claude','glm-4.5-air','','GLM 4.5 Air',10,100,'enabled','2026-01-14 15:42:25.000','2026-01-14 17:42:59.000'),
  (14,9,'zhipu-claude','glm-4.7','','GLM 4.7',10,100,'enabled','2026-01-14 15:42:58.000','2026-01-14 17:42:43.000');

UNLOCK TABLES;



LOCK TABLES `permissions` WRITE;

INSERT INTO `permissions` (`id`, `created_at`, `updated_at`, `name`, `code`, `data`, `desc`)
VALUES
  (1,'2026-01-25 18:44:24.808','2026-01-25 19:47:06.442','Relay 编辑账号','relay:account:edit','[{\"path\": \"admin.v1.RelayService/CreateAccount\", \"method\": \"POST\"}, {\"path\": \"admin.v1.RelayService/UpdateAccount\", \"method\": \"POST\"}]','包括新增、修改账号信息权限'),
  (2,'2026-01-25 19:02:14.090','2026-01-25 19:46:43.847','Relay 删除账号','relay:account:delete','[{\"path\": \"admin.v1.RelayService/DeleteAccounts\", \"method\": \"POST\"}]','删除账号'),
  (3,'2026-01-25 19:49:49.838','2026-01-25 19:50:17.868','Realy 编辑用户 ApiKey','relay:account-api-key:edit','[{\"path\": \"admin.v1.RelayService/CreateAccountApiKey\", \"method\": \"POST\"}, {\"path\": \"admin.v1.RelayService/UpdateAccountApiKey\", \"method\": \"POST\"}, {\"path\": \"admin.v1.RelayService/GetAccountList\", \"method\": \"POST\"}]','新增、编辑'),
  (4,'2026-01-25 19:50:56.254','2026-01-25 19:50:56.254','Relay 删除用户 ApiKey','relay:account-api-key:delete','[{\"path\": \"admin.v1.RelayService/DeleteAccountApiKeys\", \"method\": \"POST\"}]',''),
  (5,'2026-01-25 19:52:00.238','2026-01-25 19:52:00.238','Relay 账号列表','relay:account:list','[{\"path\": \"admin.v1.RelayService/GetAccountList\", \"method\": \"POST\"}]',''),
  (6,'2026-01-25 19:54:29.939','2026-01-25 19:54:29.939','Relay 账号APiKey 列表','relay:account-api-key:list','[{\"path\": \"admin.v1.RelayService/GetAccountApiKeyList\", \"method\": \"POST\"}]',''),
  (7,'2026-01-26 01:13:48.631','2026-01-26 02:08:05.210','Auth 获得登录用户详情','auth:user:info','[{\"path\": \"admin.v1.AuthService/GetUserInfo\", \"method\": \"POST\"}]','获得登录用户详情'),
  (8,'2026-01-26 02:08:37.247','2026-01-26 02:08:37.247','System 用户列表','system.user.list','[{\"path\": \"admin.v1.SystemService/GetUserList\", \"method\": \"POST\"}]',''),
  (9,'2026-01-26 02:09:13.232','2026-01-26 02:09:13.232','System 编辑用户','system.user.edit','[{\"path\": \"admin.v1.SystemService/CreateUser\", \"method\": \"POST\"}, {\"path\": \"admin.v1.SystemService/UpdateUser\", \"method\": \"POST\"}]',''),
  (10,'2026-01-26 02:09:46.543','2026-01-26 02:09:46.543','System 删除用户','system.user.delete','[{\"path\": \"admin.v1.SystemService/DeleteUsers\", \"method\": \"POST\"}]',''),
  (11,'2026-01-26 02:18:24.570','2026-01-26 02:20:03.208','System 角色列表','system:role:list','[{\"path\": \"admin.v1.SystemService/GetRoleList\", \"method\": \"POST\"}]',''),
  (12,'2026-01-26 02:18:58.104','2026-01-26 02:18:58.104','System 编辑角色','system:role:edit','[{\"path\": \"admin.v1.SystemService/CreateRole\", \"method\": \"POST\"}, {\"path\": \"admin.v1.SystemService/UpdateRole\", \"method\": \"POST\"}]',''),
  (13,'2026-01-26 02:19:52.020','2026-01-26 02:19:52.020','System 删除角色','system:role:delete','[{\"path\": \"admin.v1.SystemService/DeleteRoles\", \"method\": \"POST\"}]',''),
  (14,'2026-01-26 02:20:48.113','2026-01-26 02:20:48.113','System 菜单列表','system:menu:list','[{\"path\": \"admin.v1.SystemService/GetMenuList\", \"method\": \"POST\"}]',''),
  (15,'2026-01-26 02:21:24.393','2026-01-26 02:21:24.393','System 编辑菜单','system:menu:edit','[{\"path\": \"admin.v1.SystemService/CreateMenu\", \"method\": \"POST\"}, {\"path\": \"admin.v1.SystemService/UpdateMenu\", \"method\": \"POST\"}]',''),
  (16,'2026-01-26 02:21:53.981','2026-01-26 02:21:53.981','System 删除菜单','system:menu:delete','[{\"path\": \"admin.v1.SystemService/DeleteMenus\", \"method\": \"POST\"}]',''),
  (17,'2026-01-26 02:22:23.623','2026-01-26 02:22:23.623','System 权限列表','system:permission:list','[{\"path\": \"admin.v1.SystemService/GetPermissionList\", \"method\": \"POST\"}]',''),
  (18,'2026-01-26 02:23:18.294','2026-01-26 02:23:43.527','System 编辑权限','system:permission:edit','[{\"path\": \"admin.v1.SystemService/CreatePermission\", \"method\": \"POST\"}, {\"path\": \"admin.v1.SystemService/UpdatePermission\", \"method\": \"POST\"}]',''),
  (19,'2026-01-26 02:25:00.817','2026-01-26 02:25:00.817','System 删除权限','system:permission:delete','[{\"path\": \"admin.v1.SystemService/DeletePermissions\", \"method\": \"POST\"}]',''),
  (20,'2026-01-26 02:26:04.029','2026-01-26 02:26:04.029','Relay 请求列表','relay:request:list','[{\"path\": \"admin.v1.RelayService/GetRequestList\", \"method\": \"POST\"}]',''),
  (21,'2026-01-26 02:26:25.579','2026-01-26 02:26:25.579','Relay 删除请求','relay:request:delete','[{\"path\": \"admin.v1.RelayService/DeleteRequests\", \"method\": \"POST\"}]',''),
  (22,'2026-01-26 02:26:50.252','2026-01-26 02:26:50.252','Relay 账单列表','relay:ledger:list','[{\"path\": \"admin.v1.RelayService/GetLedgerList\", \"method\": \"POST\"}]',''),
  (23,'2026-01-26 02:27:11.805','2026-01-26 02:27:11.805','Relay 删除账单','relay:ledger:delete','[{\"path\": \"admin.v1.RelayService/DeleteLedgers\", \"method\": \"POST\"}]','');

UNLOCK TABLES;


LOCK TABLES `providers` WRITE;

INSERT INTO `providers` (`id`, `code`, `name`, `base_url`, `status`, `updated_at`, `created_at`)
VALUES
  (1,'openai','OpenAI','https://api.openai.com','enabled','2025-12-20 14:38:49.000','2025-12-20 14:38:49.000'),
  (2,'deepseek','DeepSeek','https://api.deepseek.com','enabled','2026-01-18 22:56:36.642','2025-12-20 14:39:11.000'),
  (3,'zhipu','智谱','https://open.bigmodel.cn','enabled','2026-01-14 02:24:12.000','2025-12-25 00:00:00.000'),
  (6,'anthropic','anthropic','https://api.anthropic.com','enabled','2025-12-30 17:58:48.000','2025-12-30 17:58:48.000'),
  (9,'zhipu-claude','智谱-Claude','https://open.bigmodel.cn','enabled','2026-01-14 16:28:53.000','2026-01-14 16:28:53.000');

UNLOCK TABLES;


LOCK TABLES `roles` WRITE;

INSERT INTO `roles` (`id`, `name`, `code`, `is_super_admin`, `permission`, `description`, `status`, `created_at`, `updated_at`)
VALUES
  (10000,'开发人员','developer',1,'{\"home\": \"home\", \"buttons\": null, \"menu_ids\": [10009]}','开发人员','enabled','2025-05-23 00:01:02.000','2026-01-23 02:14:41.765');

UNLOCK TABLES;

LOCK TABLES `users` WRITE;

INSERT INTO `users` (`username`, `nickname`, `email`, `phone`, `gender`, `roles`, `password`, `avatar_url`, `description`, `status`, `created_at`, `updated_at`)
VALUES
	('admin', '管理员', 'admin@modelgate.io', '13812345678', 'male', '10000', '$2a$10$DSoYy8I0/cpbXUypPVyTE.HWjEnfD0dDoYAeCpXreKsobEPUv0fXK', '', ' 超级管理员', 'enabled', '2025-04-26 00:07:33.000', '2025-05-25 22:54:52.000');

UNLOCK TABLES;

LOCK TABLES `relay_usages` WRITE;

INSERT INTO `relay_usages` (`id`, `total_request`, `total_success`, `total_failed`, `total_point`, `created_at`, `updated_at`)
VALUES
	(1, 0, 0, 0, 0, '2026-02-01 00:00:00.000', '2026-02-01 00:00:00.000');

UNLOCK TABLES;