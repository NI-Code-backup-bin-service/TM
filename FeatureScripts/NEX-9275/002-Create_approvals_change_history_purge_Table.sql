--multiline
CREATE TABLE `approvals_change_history_purge` (
                                                  `approval_id` bigint NOT NULL AUTO_INCREMENT,
                                                  `profile_id` int DEFAULT NULL,
                                                  `data_element_id` int DEFAULT NULL,
                                                  `change_type` int DEFAULT NULL,
                                                  `current_value` mediumtext,
                                                  `new_value` mediumtext,
                                                  `created_at` datetime DEFAULT NULL,
                                                  `approved_at` datetime DEFAULT NULL,
                                                  `approved` int DEFAULT NULL,
                                                  `created_by` varchar(45) DEFAULT NULL,
                                                  `approved_by` varchar(45) DEFAULT NULL,
                                                  `tid_id` varchar(45) DEFAULT NULL,
                                                  `merchant_id` varchar(45) DEFAULT NULL,
                                                  `acquirer` text,
                                                  `is_encrypted` tinyint NOT NULL DEFAULT '0',
                                                  `is_password` tinyint NOT NULL DEFAULT '0',
                                                  PRIMARY KEY (`approval_id`)
);