--multiline;
CREATE PROCEDURE get_site_change_history(IN siteId INT)
BEGIN
	SELECT
      de.name,
      pd.datavalue,
      pd.updated_by,
      pd.updated_at,
      pd.approved
    FROM profile_data pd
    RIGHT JOIN site_profiles sp ON sp.profile_id = pd.profile_id 
    LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
    LEFT JOIN profile p ON p.profile_id = pd.profile_id
    LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
    WHERE sp.site_id = siteId AND pt.priority = 2
    ORDER BY pd.updated_at DESC;
END