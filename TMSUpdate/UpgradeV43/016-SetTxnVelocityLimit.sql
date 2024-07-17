--multiline;
CREATE PROCEDURE `set_txn_velocity_limit`(
    IN txnID VARCHAR(255),
    IN velocityID VARCHAR(255),
    IN limitType VARCHAR(255),
    IN txnType VARCHAR(255),
    IN `value` INT
)
BEGIN
    IF (SELECT count(txn_limit_id) FROM velocity_limits_txn
        WHERE txn_limit_id = txnID) > 0
    THEN
        UPDATE velocity_limits_txn
            SET
                velocity_limit_id = velocityID,
                limit_type = (SELECT limit_type_id FROM txn_limit_types WHERE limit_type=limitType),
                txn_type = (SELECT txn_type_id FROM txn_types WHERE txn_type=txnType),
                limit_value = `value`
            WHERE txn_limit_id = txnID;
    ELSE
        INSERT INTO velocity_limits_txn (txn_limit_id, velocity_limit_id, limit_type, txn_type, limit_value)
        VALUES (
        txnID,
        velocityID,
        (SELECT limit_type_id FROM txn_limit_types WHERE limit_type=limitType),
        (SELECT txn_type_id FROM txn_types WHERE txn_type=txnType),
        `value`
        );
    END IF;
END