--multiline
CREATE PROCEDURE `get_profile_change_approval_history`(IN typeId INT)
BEGIN
	SELECT
		pd.profile_data_id,
		p.name as profileName,
        de.name,
        pd.datavalue,
        pd.updated_by,
        pd.updated_at
    FROM profile_data pd
    LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
    LEFT JOIN profile p ON p.profile_id = pd.profile_id
    WHERE profile_type_id = typeId AND pd.approved = 0
    ORDER BY pd.updated_at DESC;
END