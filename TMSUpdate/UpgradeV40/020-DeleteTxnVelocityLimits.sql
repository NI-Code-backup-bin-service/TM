--multiline;
CREATE PROCEDURE `delete_txn_velocity_limit`(
    IN siteID INT,
    IN `level` INT,
    IN tidID INT
)
BEGIN
    DELETE FROM velocity_limits_txn
    WHERE velocity_limit_id IN (SELECT velocity_limit_id FROM velocity_limits WHERE site_id = siteID AND limit_level = `level` AND tid_id = tidID);
END