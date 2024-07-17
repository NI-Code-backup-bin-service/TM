--multiline;
CREATE TABLE `tid_user_override` (
  `tid_user_id` int(11) NOT NULL AUTO_INCREMENT,
  `tid_id` int(8) unsigned DEFAULT NULL,
  `Username` varchar(45) DEFAULT NULL,
  `PIN` varchar(45) DEFAULT NULL,
  `Modules` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`tid_user_id`),
  UNIQUE KEY `UQ_Username` (`Username`),
  UNIQUE KEY `UQ_PIN` (`PIN`),
  KEY `tid_tid_users_FK_idx` (`tid_id`),
  KEY `PIN` (`PIN`),
  CONSTRAINT `FX_TIDID` FOREIGN KEY (`tid_id`) REFERENCES `tid` (`tid_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=latin1;
