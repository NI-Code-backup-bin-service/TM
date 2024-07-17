--multiline
INSERT IGNORE INTO profile_data
(profile_id,
 data_element_id,
 datavalue,
 version,
 updated_at,
 updated_by,
 created_at,
 created_by,
 approved,
 overriden,
 is_encrypted)
SELECT profile_id,
       (SELECT data_element_id FROM data_element WHERE name = 'autoCutOver'),
       'false',
       1,
       CURRENT_TIMESTAMP,
       'system',
       CURRENT_TIMESTAMP,
       'system',
       1,
       0,
       0
FROM (
         SELECT DISTINCT(profile_id) FROM site_data WHERE level = 'site' AND profile_id IN (SELECT DISTINCT(profile_id) FROM profile_data_group WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay') and site_id NOT IN (SELECT DISTINCT(site_id) FROM site_data WHERE level = 'chain' AND profile_id IN (SELECT DISTINCT(profile_id) FROM profile_data_group WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay'))))
     ) AS subquery;