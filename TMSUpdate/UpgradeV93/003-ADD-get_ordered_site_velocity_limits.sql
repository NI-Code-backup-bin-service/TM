--multiline
CREATE PROCEDURE `get_ordered_site_velocity_limits`(
    IN siteID INT,
    IN tidID INT,
    IN `level` INT
)
BEGIN
SELECT
    vl.velocity_limit_id,
    COALESCE(s.scheme_name, ''),
    vl.transaction_limit_daily,
    vl.transaction_limit_batch,
    vl.single_transaction_limit,
    vl.cumulative_daily,
    vl.cumulative_batch,
    vl.limit_index
FROM velocity_limits vl
         LEFT JOIN schemes s
                   ON vl.scheme = s.scheme_id
WHERE vl.site_id = siteID AND vl.tid_id = tidID AND vl.limit_level = `level`;
END