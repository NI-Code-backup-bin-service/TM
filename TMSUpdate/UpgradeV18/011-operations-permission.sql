--multiline
CREATE TABLE `operations_permission` (
  `permission_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`permission_id`),
  UNIQUE KEY `name` (`name`)
);