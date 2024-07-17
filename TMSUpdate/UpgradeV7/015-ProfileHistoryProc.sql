--multiline
CREATE PROCEDURE `get_profile_change_history`(IN profileId INT)
BEGIN
	SELECT
      de.name,
      pd.datavalue,
      pd.updated_by,
      pd.updated_at,
      pd.approved
    FROM profile_data pd
    LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
    WHERE pd.profile_id = profileId
    ORDER BY pd.updated_at DESC;
END