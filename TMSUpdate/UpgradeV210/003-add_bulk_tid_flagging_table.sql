--multiline;
CREATE TABLE IF NOT EXISTS `bulk_tid_flagging` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `filename` varchar(100) NOT NULL UNIQUE,
    `created_by` varchar(45) DEFAULT NULL,
    `approved_by` varchar(45) DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `approved_at` datetime DEFAULT NULL,
    `approved` int(11) DEFAULT 0,
    PRIMARY KEY (`id`)
    );