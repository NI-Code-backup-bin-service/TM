--multiline;
INSERT INTO approvals (profile_id, data_element_id, change_type, current_value, new_value, created_at, created_by, approved_at, approved_by, approved)
SELECT pd.profile_id, pd.data_element_id, 1, pd2.datavalue, pd.datavalue, pd.created_at, pd.created_by, pd.updated_at, pd.updated_by, 1 FROM profile_data pd  
  LEFT JOIN data_element de ON de.data_element_id = pd.data_element_id
  LEFT JOIN profile_data pd2 on pd2.profile_id = pd.profile_id
  AND pd2.data_element_id = pd.data_element_id
  AND pd2.version = (SELECT MAX(d.version) FROM profile_data d
  WHERE d.data_element_id = pd2.data_element_id AND d.profile_id = pd2.profile_id AND d.approved != 2 AND d.version < pd.version)
  ORDER BY pd.updated_at DESC;
