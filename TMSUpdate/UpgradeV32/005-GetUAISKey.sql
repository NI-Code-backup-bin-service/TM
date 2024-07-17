--multiline;
CREATE  PROCEDURE `getUAISKey`(in id varchar(50))
BEGIN
	SELECT `Serial`,PublicKey FROM uaiskeys WHERE `Serial` = id;
END