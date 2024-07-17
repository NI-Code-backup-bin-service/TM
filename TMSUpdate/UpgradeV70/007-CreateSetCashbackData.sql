--multiline;
CREATE PROCEDURE `set_cashback_data` (profileID INT, cashbackData TEXT)
BEGIN
	INSERT IGNORE INTO `cashback` (`profile_id`, `cashback_data`) VALUES (profileID, cashbackData) ON DUPLICATE KEY UPDATE `cashback_data` = cashbackData;
END