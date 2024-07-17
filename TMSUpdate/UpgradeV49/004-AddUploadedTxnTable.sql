--multiline;
CREATE TABLE IF NOT EXISTS `uploaded_txns` (
  `filename` varchar(255) NOT NULL,
  `checksum` varchar(255) NOT NULL,
  PRIMARY KEY (`filename`),
  UNIQUE KEY `filename_UNIQUE` (`filename`),
  UNIQUE KEY `checksum_UNIQUE` (`checksum`),
  KEY `checksum_Asc` (`checksum`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
