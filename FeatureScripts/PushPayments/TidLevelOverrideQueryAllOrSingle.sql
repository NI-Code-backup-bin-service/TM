SET @modulesDataGroupId = (SELECT data_group_id FROM data_group WHERE `name` = 'modules');
SET @activeDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'active' AND data_group_id = @modulesDataGroupId);

SET @tid = NULL;

SELECT t.*
FROM tid t
JOIN tid_site ts ON t.tid_id = ts.tid_id AND ts.tid_profile_id IS NOT NULL
JOIN `profile` p ON ts.tid_profile_id = p.profile_id
JOIN profile_data pd ON p.profile_id = pd.profile_id AND pd.data_element_id = @activeDataElementId
WHERE (@tid IS NULL OR t.tid_id = @tid);