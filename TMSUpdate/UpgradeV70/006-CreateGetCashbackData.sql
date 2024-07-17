--multiline;
CREATE PROCEDURE `get_cashback_data` (profileID INT)
BEGIN
	SELECT `cashback_data` FROM `cashback` WHERE `profile_id` = profileID;
END