SET @v_tids := '';
SET @v_newValue := '';
SELECT @v_data_element_id := data_element_id FROM data_element WHERE name = 'preAuthMax' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'modules');
UPDATE profile_data SET datavalue = @v_newValue, updated_at = NOW(), updated_by = 'system' WHERE data_element_id = @v_data_element_id AND profile_id IN(SELECT tid_profile_id FROM tid_site WHERE FIND_IN_SET(tid_id, @v_tids));