--multiline;
CREATE TABLE `nikeys` (
  `Serial` varchar(50) NOT NULL,
  `PublicKey` varchar(2000) DEFAULT NULL,
  `PrivateKey` varchar(2000) DEFAULT NULL,
  `Type` varchar(5) DEFAULT NULL,
  `StartDate` datetime DEFAULT NULL,
  `Exchanged` int(11) DEFAULT NULL,
  PRIMARY KEY (`Serial`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
