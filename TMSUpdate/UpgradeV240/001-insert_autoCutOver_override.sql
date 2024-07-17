--multiline
INSERT IGNORE INTO profile_data
(profile_id,
 data_element_id,
 datavalue,
 version,
 updated_by,
 created_by,
 approved,
 overriden,
 is_encrypted)
SELECT profile_id,
       (SELECT data_element_id FROM data_element WHERE name = 'autoCutOver'),
       'true',
       1,
       'system',
       'system',
       1,
       1,
       0
FROM (
         SELECT profile_id from profile_data where data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'time' and data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay')) and LENGTH(datavalue) > 5
     ) AS subquery;