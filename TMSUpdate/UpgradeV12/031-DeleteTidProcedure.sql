--multiline;
CREATE PROCEDURE `delete_tid`(IN tid INT, IN site int)
BEGIN
    START TRANSACTION;
        SET @profileID = (SELECT tid_profile_id FROM tid_site WHERE tid_id = tid LIMIT 1);
        DELETE FROM tid_site WHERE tid_id = tid AND site_id = site;
        DELETE FROM tid_updates WHERE tid_id = tid;
        DELETE FROM approvals where profile_id = @profileID and approved = 0;
        IF NOT EXISTS(SELECT * FROM tid_site WHERE tid_id = tid)
        THEN
            DELETE FROM tid WHERE tid_id = tid;
        END IF;
    COMMIT;
END