--multiline
INSERT IGNORE INTO profile_data_group(
      profile_id,
      data_group_id,
      version,
      updated_at,
      updated_by,
      created_at,
      created_by
    )
SELECT
    profile_id AS profile_id,
    (SELECT data_group_id FROM data_group WHERE NAME = "receipt") AS data_group_id,
    1 AS version,
    current_timestamp AS updated_at,
    "system" AS updated_by,
    current_timestamp AS created_at,
    "system" AS created_by
FROM profile
WHERE profile_type_id = (SELECT profile_type_id FROM profile_type WHERE NAME = "tid");