--multiline;
CREATE PROCEDURE `getNIKey`(in id varchar(50))
BEGIN
	SELECT `Serial`,PublicKey,PrivateKey FROM nikeys WHERE `Serial` = id;
END