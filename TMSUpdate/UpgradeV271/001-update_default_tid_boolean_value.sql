--multiline;
INSERT IGNORE INTO profile_data (
    profile_id,
    datavalue,
    data_element_id,
    version,
    updated_at,
    updated_by,
    approved
)
SELECT
    (SELECT profile_id FROM profile WHERE name = 'global'),
    'false',
    data_element_id,
    1,
    NOW(),
    'system',
    1
FROM
    data_element
WHERE
    datatype = 'BOOLEAN';