--multiline;
CREATE TABLE `site_level_users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `site_id` int(11) DEFAULT NULL,
  `Username` varchar(256) DEFAULT NULL,
  `PIN` varchar(5) DEFAULT NULL,
  `Modules` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `SitePIN_UQ` (`site_id`,`PIN`),
  UNIQUE KEY `SiteUsername_UQ` (`Username`,`site_id`),
  KEY `SiteUser_FK_idx` (`site_id`),
  KEY `PIN_IX` (`PIN`),
  KEY `Site_IX` (`site_id`),
  CONSTRAINT `SiteUser_FK` FOREIGN KEY (`site_id`) REFERENCES `site` (`site_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=latin1;
