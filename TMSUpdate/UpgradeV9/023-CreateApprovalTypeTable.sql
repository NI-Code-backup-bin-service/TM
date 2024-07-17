--multiline;
CREATE TABLE IF NOT EXISTS `approval_type` (
  `approval_type_id` int(11) NOT NULL,
  `approval_type_name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`approval_type_id`)
);
