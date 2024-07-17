--multiline
CREATE PROCEDURE `store_PIN`(IN tid int, IN PIN varchar(5), IN timeout int, IN is_encrypted BOOLEAN)
BEGIN
    UPDATE tid t SET t.PIN = PIN, t.ExpiryDate = DATE_ADD(NOW(), INTERVAL timeout MINUTE), t.is_encrypted = is_encrypted
    WHERE t.tid_id  = tid;
END