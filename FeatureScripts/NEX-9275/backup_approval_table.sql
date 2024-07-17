--multiline
INSERT INTO `approvals_change_history_purge`(
    approval_id,
    profile_id,
    data_element_id,
    change_type,
    current_value,
    new_value,
    created_at,
    approved_at,
    approved,
    created_by,
    approved_by,
    tid_id,
    merchant_id,
    acquirer,
    is_encrypted,
    is_password)

SELECT
    approval_id,
    profile_id,
    data_element_id,
    change_type,
    current_value,
    new_value,
    created_at,
    approved_at,
    approved,
    created_by,
    approved_by,
    tid_id,
    merchant_id,
    acquirer,
    is_encrypted,
    is_password
FROM approvals a
WHERE a.created_at <= date_sub(now(), interval 6 month)