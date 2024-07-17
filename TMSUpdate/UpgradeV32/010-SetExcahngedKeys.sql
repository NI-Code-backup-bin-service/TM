--multiline;
CREATE PROCEDURE `SetKeysExchanged`(IN keySerial varchar(50))
BEGIN
	Update nikeys nk SET nk.Exchanged = 1 WHERE nk.`Serial` = keySerial;
END