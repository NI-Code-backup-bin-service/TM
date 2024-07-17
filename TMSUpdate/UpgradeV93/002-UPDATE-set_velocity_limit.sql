--multiline
CREATE PROCEDURE `set_velocity_limit`(
    IN id VARCHAR(255),
    IN siteID INT,
    IN tidID INT,
    IN limitLevel INT,
    IN schemeIN VARCHAR(255),
    IN dailyTL INT,
    IN batchTL INT,
    IN singleTL INT,
    IN dailyLimit INT,
    IN batchLimit INT,
    IN limitIndex INT
)
BEGIN
    IF (SELECT count(velocity_limit_id) FROM velocity_limits
        WHERE velocity_limit_id = id) > 0
    THEN
        UPDATE velocity_limits
        SET
            site_id = siteID,
            tid_id = tidID,
            limit_level = limitLevel,
            scheme = (SELECT scheme_id FROM schemes WHERE scheme_name=schemeIN),
            transaction_limit_daily = dailyTL,
            transaction_limit_batch = batchTL,
            single_transaction_limit = singleTL,
            cumulative_daily = dailyLimit,
            cumulative_batch = batchLimit,
            limit_index = limitIndex
        WHERE velocity_limit_id = id;
    ELSE
        INSERT INTO velocity_limits (velocity_limit_id, site_id, tid_id, limit_level, scheme, transaction_limit_daily, transaction_limit_batch, single_transaction_limit, cumulative_daily, cumulative_batch, limit_index)
        VALUES (
                   id,
                   siteID,
                   tidID,
                   limitLevel,
                   (SELECT scheme_id FROM schemes WHERE scheme_name=schemeIN),
                   dailyTL,
                   batchTL,
                   singleTL,
                   dailyLimit,
                   batchLimit,
                   limitIndex
               );
    END IF;
END