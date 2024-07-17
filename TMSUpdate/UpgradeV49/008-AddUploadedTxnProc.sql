--multiline;
CREATE PROCEDURE `add_uploaded_txn`(IN fname VARCHAR(255), IN c_sum VARCHAR(255))
BEGIN
	INSERT INTO uploaded_txns(filename, `checksum`) VALUES (fname,c_sum);
END