--multiline;
CREATE  PROCEDURE `getLatestUAISKey`(in keytype varchar(50))
BEGIN
	SELECT `Serial`,PublicKey FROM uaiskeys uk WHERE uk.`type` = keytype  ORDER BY StartDate DESC LIMIT 1;
END