INSERT IGNORE INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted, not_overridable) SELECT p.profile_id, de.data_element_id, 0, 1, NOW(), 'system', NOW(), 'system', 1, 0, 0, 0 FROM profile p JOIN profile_data_group pdg ON p.profile_id=pdg.profile_id JOIN data_group dg ON pdg.data_group_id=dg.data_group_id JOIN data_element de ON dg.data_group_id=de.data_group_id WHERE p.profile_type_id = 4 AND dg.name = 'nol' AND de.name = 'driverId';