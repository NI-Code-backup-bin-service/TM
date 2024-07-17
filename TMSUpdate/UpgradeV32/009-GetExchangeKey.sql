--multiline;
CREATE PROCEDURE `GetExchangeKey`(IN keytype varchar(5))
BEGIN
	SELECT `Serial`,PublicKey,PrivateKey FROM nikeys nk WHERE nk.`type` = keytype AND nk.Exchanged = 0 ORDER BY StartDate DESC LIMIT 1;
END