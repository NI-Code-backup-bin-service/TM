ALTER TABLE profile_data ADD UNIQUE INDEX UQ (profile_id ASC, data_element_id ASC);
ALTER TABLE profile_data_group ADD UNIQUE INDEX UQ (profile_id ASC, data_group_id ASC);
