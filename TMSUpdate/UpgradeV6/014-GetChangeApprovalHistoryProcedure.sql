--multiline
CREATE PROCEDURE get_change_approval_history()
BEGIN
	SELECT
	      pd.profile_data_id,
	      sp.site_id,
        de.name,
        pd.datavalue,
        pd.updated_by,
        pd.updated_at
    FROM profile_data pd
    RIGHT JOIN site_profiles sp ON sp.profile_id = pd.profile_id
    LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
    LEFT JOIN profile p ON p.profile_id = pd.profile_id
    LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
    WHERE pt.priority = 1 OR pt.priority = 2 AND pd.approved = 0
    ORDER BY pd.updated_at DESC;
END;