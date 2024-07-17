--multiline;
CREATE PROCEDURE `store_Activation`(IN tid int)
BEGIN
	UPDATE tid t SET t.ActivationDate = TIME(NOW())
    WHERE t.tid_id  = tid;
END;