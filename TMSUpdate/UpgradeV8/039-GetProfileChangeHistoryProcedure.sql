--multiline
CREATE PROCEDURE get_profile_change_history(IN profileId INT)
BEGIN
  SELECT
    de.name,
    pd2.datavalue as original_value,
    pd.datavalue,
    pd.updated_by,
    pd.updated_at,
    pd.approved
  FROM profile_data pd
  LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
  LEFT JOIN profile_data pd2 on pd2.profile_id = pd.profile_id
  AND pd2.data_element_id = pd.data_element_id
  AND pd2.version = (SELECT MAX(d.version) FROM profile_data d
  WHERE d.data_element_id = pd2.data_element_id AND d.profile_id = pd2.profile_id AND d.approved != 2 AND d.version < pd.version)
  WHERE pd.profile_id = profileId
  ORDER BY pd.updated_at DESC;
END