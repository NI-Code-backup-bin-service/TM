--multiline
CREATE PROCEDURE `delete_tid_override`( IN tid VARCHAR(255) )
BEGIN
    SET @tid_profile_id = (SELECT tid_profile_id FROM tid_site WHERE `tid_id` = tid LIMIT 1);
    CALL remove_tid_override(@tid_profile_id);
END