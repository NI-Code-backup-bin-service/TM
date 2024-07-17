--multiline;
CREATE TABLE `keystore` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `identifier` varchar(45) NOT NULL,
  `lkey` varchar(256) NOT NULL,
  `rkey` varchar(256) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `identifier_UNIQUE` (`identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
