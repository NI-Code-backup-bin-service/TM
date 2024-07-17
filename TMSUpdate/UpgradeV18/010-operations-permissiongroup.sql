--multiline
CREATE TABLE `operations_permissiongroup` (
  `group_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `default_group` int(11) DEFAULT '0',
  PRIMARY KEY (`group_id`),
  UNIQUE KEY `name` (`name`)
);