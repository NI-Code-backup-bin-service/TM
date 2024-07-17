--multiline;
CREATE PROCEDURE `SaveUAISKeys`(
IN KeySerial varchar(50),
IN PubKey varchar(2000),
IN keyType varchar(5)
)
BEGIN
	INSERT INTO uaiskeys(`Serial`, PublicKey, `type`, StartDate) VALUES (KeySerial, PubKey,keyType, NOW());
END