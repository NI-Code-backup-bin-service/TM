--multiline;
CREATE PROCEDURE `SaveNIKeys`(
IN KeySerial varchar(50),
IN PubKey varchar(250),
IN PrivKey varchar(250),
IN keyType varchar(5)
)
BEGIN
	INSERT INTO nikeys(`Serial`, PublicKey,PrivateKey, `type`, StartDate) VALUES (KeySerial, PubKey,PrivKey,keyType, NOW());
END