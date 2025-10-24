-- Check if 'demo' database exists; if not, create it
CREATE DATABASE IF NOT EXISTS promo CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Use the 'demo' database
USE promo;

CREATE TABLE IF NOT EXISTS `promo_code` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(255) DEFAULT NULL COMMENT 'Promo Name',
  `description` varchar(255) DEFAULT NULL COMMENT 'Promo Description',
  `termCondition` varchar(255) DEFAULT NULL COMMENT 'Promo Term & Condition',
  `startDate` int(10) unsigned NULL COMMENT 'Promo Start Date',
  `endDate` int(10) unsigned NULL COMMENT 'Promo End Date',
  `banner` varchar(255) DEFAULT NULL COMMENT 'Promo Banner',
  `rules` ENUM('otg', 'sp', 'pw') NULL COMMENT 'Promo Rules e.g Order Total is Greater than, Spesifict Product, Payment With',
  `amount` decimal(13, 2) NULL COMMENT 'Promo Amount Rules OTG(Order Total is Greater than)',
  `product_slug` varchar(255) DEFAULT NULL COMMENT 'Promo Product Slug Rules SP(Spesifict Product)',
  `paymentMethod` decimal(13, 2) NULL COMMENT 'Promo Paymkent Method Rules PW(Payment With)',
  `promoAction` ENUM('fixed', 'percent') NOT NULL COMMENT 'Promo Action',
  `promoType` varchar(255) NOT NULL COMMENT 'Promo Type e.g pax, total',
  `type` varchar(255) NOT NULL COMMENT 'Type e.g tour, flight, hotel, etc',
  `promoCode` varchar(255) NOT NULL COMMENT 'Promo code',
  `customerLimit` int(10) NULL COMMENT 'Customer Limit',
  `newCustomer` tinyint(3) unsigned NOT NULL COMMENT 'Promo New Customer',
  `quantity` int(10) NOT NULL COMMENT 'Promo Quantity',
  `status` varchar(255) DEFAULT NULL COMMENT 'Promo Status',
  `specialPromo` int(10) unsigned DEFAULT NULL COMMENT 'Promo Special',
  `isDisplay` tinyint(3) unsigned NOT NULL COMMENT 'Is Display',
  `platform` varchar(50) DEFAULT NULL COMMENT 'Platform name',
  `created_at`   bigint(20) UNSIGNED                            NOT NULL,
  `created_by`   bigint(20) UNSIGNED                            NOT NULL,
  `updated_at`   bigint(20) UNSIGNED                            NOT NULL,
  `updated_by`   bigint(20) UNSIGNED                            NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `promoCode_unique` (`promoCode`),
  KEY `idx_promoStatus` (`status`) USING HASH,
  KEY `idx_promoId` (`id`) USING HASH,
  KEY `idx_platform` (`platform`) USING HASH,
  KEY `idx_combine_v1` (
    `id`,
    `status`,
    `type`,
    `platform`
  ) USING HASH,
  KEY `idx_promoType` (`type`) USING HASH,
  FULLTEXT KEY `idx_fulltext_name` (`name`),
  INDEX `idx_promo_id` (`id` ASC)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `strike_throught_price` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `description` varchar(255) DEFAULT NULL COMMENT 'Promo Description',
  `product_slug` varchar(255) DEFAULT NULL COMMENT 'Product Slug',
  `departure` INT DEFAULT NULL COMMENT 'Product Departure (opsional)',
  `termCondition` varchar(255) DEFAULT NULL COMMENT 'Promo Term & Condition',
  `startDate` int(10) unsigned NULL COMMENT 'Promo Start Date',
  `endDate` int(10) unsigned NULL COMMENT 'Promo End Date',
  `promoAction` ENUM('fixed', 'percent') NOT NULL COMMENT 'Promo Action',
  `type` varchar(255) NOT NULL COMMENT 'Type e.g tour, flight, hotel, etc',
  `status` varchar(255) DEFAULT NULL COMMENT 'Promo Status',
  `platform` varchar(50) DEFAULT NULL COMMENT 'Platform name',
  `created_at`   bigint(20) UNSIGNED                            NOT NULL,
  `created_by`   bigint(20) UNSIGNED                            NOT NULL,
  `updated_at`   bigint(20) UNSIGNED                            NOT NULL,
  `updated_by`   bigint(20) UNSIGNED                            NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_promoStatus` (`status`) USING HASH,
  KEY `idx_promoId` (`id`) USING HASH,
  KEY `idx_platform` (`platform`) USING HASH,
  KEY `idx_combine_v1` (
    `id`,
    `status`,
    `type`,
    `platform`
  ) USING HASH,
  KEY `idx_promoType` (`type`) USING HASH,
  FULLTEXT KEY `idx_fulltext_name` (`name`),
  INDEX `idx_promo_id` (`id` ASC)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;