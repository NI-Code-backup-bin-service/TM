--multiline
CREATE PROCEDURE get_profile_change_approval_history(IN typeId INT)
BEGIN
	SELECT
		pd.profile_data_id,
		p.name as profileName,
        de.name,
        pd2.datavalue as original_value,
        pd.datavalue as updated_value,
        pd.updated_by,
        pd.updated_at
    FROM profile_data pd
    LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
    LEFT JOIN profile p ON p.profile_id = pd.profile_id
    LEFT JOIN profile_data pd2 on pd2.profile_id = pd.profile_id 
    AND pd2.data_element_id = pd.data_element_id
    AND pd2.version = (SELECT MAX(d.version) FROM profile_data d 
		WHERE d.data_element_id = pd2.data_element_id AND d.profile_id = pd2.profile_id AND d.approved = 1)      
    WHERE profile_type_id = typeId AND pd.approved = 0
    ORDER BY pd.updated_at DESC;
END