--multiline;
CREATE TABLE `uaiskeys` (
  `Serial` varchar(50) NOT NULL,
  `PublicKey` varchar(2000) DEFAULT NULL,
  `Type` varchar(5) DEFAULT NULL,
  `StartDate` datetime DEFAULT NULL,
  PRIMARY KEY (`Serial`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
