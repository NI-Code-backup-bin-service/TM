--multiline;
CREATE PROCEDURE `get_tid_velocity_limits`(
    IN tidID INT,
    IN `level` INT
)
BEGIN
    SELECT
        vl.velocity_limit_id,
        COALESCE(s.scheme_name, ''),
        vl.transaction_limit_daily,
        vl.transaction_limit_batch,
        vl.single_transaction_limit
    FROM velocity_limits vl
    LEFT JOIN schemes s
    ON vl.scheme = s.scheme_id
    WHERE vl.tid_id = tidID AND vl.limit_level = `level`;
END