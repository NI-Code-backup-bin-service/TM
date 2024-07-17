--multiline;
CREATE PROCEDURE `get_users_value_by_site_id`(
    IN siteId int,
    IN profileId int,
    IN profileType TEXT,
    IN tId int
)
BEGIN
    IF profileType = "dailyTxnCleanseTime" THEN
        SELECT datavalue FROM profile_data WHERE data_element_id=(SELECT data_element_id FROM data_element WHERE name=profileType) AND profile_id = profileId;
    ELSE
        IF profileType = "siteVelocity" THEN
            SELECT
                vl.velocity_limit_id,
                COALESCE(s.scheme_name, ''),
                vl.transaction_limit_daily,
                vl.transaction_limit_batch,
                vl.single_transaction_limit,
                vl.cumulative_daily,
                vl.cumulative_batch,
                vl.limit_index
            FROM velocity_limits vl LEFT JOIN schemes s ON vl.scheme = s.scheme_id
            WHERE vl.site_id = siteId AND vl.tid_id = tId;
        ELSE
            SELECT user_id, Username, PIN , Modules from site_level_users WHERE site_id = siteId;
        END IF;
    END IF;
END