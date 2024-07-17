--multiline;
CREATE PROCEDURE `get_velocity_limits`(IN siteID INT, IN tidID INT)
BEGIN
    DECLARE schemeLevel, siteLevel INT;

    IF (SELECT COUNT(*) FROM velocity_limits WHERE site_id = siteID AND tid_id = tidID AND limit_level = 2) > 0  THEN
        SELECT 2 INTO schemeLevel;
    ELSE
        SELECT 1 INTO schemeLevel;
    END IF;

    IF (SELECT COUNT(*) FROM velocity_limits WHERE site_id = siteID AND tid_id = tidID AND limit_level = 4) > 0  THEN
        SELECT 4 INTO siteLevel;
    ELSE
        SELECT 3 INTO siteLevel;
    END IF;

    SELECT COALESCE(s.scheme_name, 'terminal'),vl.transaction_limit_daily,vl.transaction_limit_batch,vl.single_transaction_limit,vl.cumulative_daily,vl.cumulative_batch,b.limit_type,c.txn_type,a.limit_value
    FROM velocity_limits vl
    LEFT JOIN schemes s
        ON vl.scheme = s.scheme_id
    LEFT JOIN velocity_limits_txn a
        ON a.velocity_limit_id = vl.velocity_limit_id
    LEFT JOIN txn_limit_types b
        ON a.limit_type = b.limit_type_id
    LEFT JOIN txn_types c
        ON a.txn_type = c.txn_type_id
    WHERE vl.site_id = siteID AND vl.tid_id = tidID AND vl.limit_level IN (schemeLevel ,siteLevel);
END;