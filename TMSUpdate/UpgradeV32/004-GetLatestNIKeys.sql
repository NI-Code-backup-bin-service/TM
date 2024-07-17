--multiline;
CREATE PROCEDURE `getLatestNIKey`(in keytype varchar(50))
BEGIN
	SELECT `Serial`,PublicKey,PrivateKey FROM nikeys nk WHERE nk.`type` = keytype AND nk.Exchanged = 1  ORDER BY StartDate DESC LIMIT 1;
END