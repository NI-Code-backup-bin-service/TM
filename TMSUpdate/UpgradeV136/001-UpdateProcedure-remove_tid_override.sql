--multiline
CREATE PROCEDURE remove_tid_override(IN p_tid_profile_id int)
BEGIN
    # Delete all the data values from profile data which are not related to fraud. Fraud must be handled
    # separately as it is independent of the main configuration.
	DELETE pd.*
	FROM profile_data pd
	INNER JOIN data_element de
		ON pd.data_element_id = de.data_element_id
	INNER JOIN data_element_locations_data_element delde
		ON de.data_element_id = delde.data_element_id
	INNER JOIN data_element_locations del
		ON delde.location_id = del.location_id
	WHERE pd.profile_id = p_tid_profile_id
    	AND del.location_name != 'fraud';

    # Delete all the rows in profile_data_groups where there is no corresponding records set in profile_data as
    # these rows are no longer required. Other remaining rows could be for fraud data.
    DELETE
    FROM profile_data_group
    WHERE profile_id = p_tid_profile_id
      AND data_group_id NOT IN (
        SELECT de.data_group_id
        FROM profile_data pd
                 INNER JOIN data_element de on pd.data_element_id = de.data_element_id
        WHERE pd.profile_id = p_tid_profile_id
    );

    # Enumerate the remaining rows for the given profile in profile_data
    SELECT COUNT(*) INTO @v_remaining_elements_count
    FROM profile_data pd
    WHERE pd.profile_id = p_tid_profile_id;

    # If there's no remaining rows in profile data, then we can just delete the override profile altogether. We have to
    # check for remaining data as there could be fraud overrides in place.
    IF @v_remaining_elements_count = 0 THEN
        DELETE FROM approvals WHERE profile_id = p_tid_profile_id AND approved = 0;
        UPDATE tid_site SET tid_profile_id = NULL, updated_at = NOW() WHERE tid_profile_id = p_tid_profile_id;
    END IF;
END;
