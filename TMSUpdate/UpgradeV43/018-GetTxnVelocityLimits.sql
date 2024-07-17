--multiline;
CREATE PROCEDURE `get_txn_velocity_limits`(
    IN limitID VARCHAR(255)
)
BEGIN
    SELECT
        a.txn_limit_id,
        b.limit_type,
        c.txn_type,
        c.display_name,
        a.limit_value
    FROM velocity_limits_txn a
    LEFT JOIN txn_limit_types b
    ON a.limit_type = b.limit_type_id
    LEFT JOIN txn_types c
    ON a.txn_type = c.txn_type_id
    WHERE a.velocity_limit_id = limitID;
END