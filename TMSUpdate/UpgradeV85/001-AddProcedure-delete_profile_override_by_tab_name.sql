--multiline
CREATE PROCEDURE delete_profile_override_by_tab_name(IN p_profileId INT, IN p_tab_name VARCHAR(255))
BEGIN
    DELETE FROM profile_data WHERE profile_data_id IN (
        SELECT pd.profile_data_id
        FROM (SELECT * FROM profile_data) pd
                 INNER JOIN profile p ON
                pd.profile_id = p.profile_id
                 INNER JOIN data_element_locations_data_element delde ON
                pd.data_element_id = delde.data_element_id
                 INNER JOIN data_element_locations del ON
                    delde.location_id = del.location_id
                AND
                    del.profile_type_id = p.profile_type_id
        WHERE pd.profile_id = p_profileId AND location_name = p_tab_name);
END;