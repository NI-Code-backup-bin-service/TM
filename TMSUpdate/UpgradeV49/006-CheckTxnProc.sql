--multiline;
CREATE PROCEDURE `checkUploadedTxn`(IN new_checksum VARCHAR(255))
BEGIN
	SELECT filename FROM uploaded_txns WHERE `checksum` = new_checksum;
END