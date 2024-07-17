--multiline;
CREATE TABLE IF NOT EXISTS `approvals` (
  `approval_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `profile_id` int(11) DEFAULT NULL,
  `data_element_id` int(11) DEFAULT NULL,
  `change_type` int(11) DEFAULT NULL,
  `current_value` TEXT DEFAULT NULL,
  `new_value` TEXT DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `approved_at` datetime DEFAULT NULL,
  `approved` int(11) DEFAULT NULL,
  `created_by` varchar(45) DEFAULT NULL,
  `approved_by` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`approval_id`),
  KEY `FK_approval_type_idx` (`change_type`),
  KEY `FK_data_element_idx` (`data_element_id`),
  KEY `FK_profile_id_idx` (`profile_id`),
  CONSTRAINT `FK_approval_type` FOREIGN KEY (`change_type`) REFERENCES `approval_type` (`approval_type_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_data_element` FOREIGN KEY (`data_element_id`) REFERENCES `data_element` (`data_element_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_profile_id` FOREIGN KEY (`profile_id`) REFERENCES `profile` (`profile_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
);
