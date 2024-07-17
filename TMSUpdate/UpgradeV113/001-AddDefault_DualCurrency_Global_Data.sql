SELECT data_group_id INTO @dc_data_group_id FROM data_group where name = 'dualCurrency';
#
# Insert Global dualCurrency/cardDefinitions from emv/cardDefinitions
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'emv' AND de.name = 'cardDefinitions';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'cardDefinitions';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);
#
# Insert Global dualCurrency/contactApplicationConfigShared from emv/contactApplicationConfigShared
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'emv' AND de.name = 'contactApplicationConfigShared';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'contactApplicationConfigShared';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);
#
# Insert Global dualCurrency/ctlsApplicationConfigShared from emv/ctlsApplicationConfigShared
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'emv' AND de.name = 'ctlsApplicationConfigShared';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'ctlsApplicationConfigShared';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);
#
# Insert Global dualCurrency/terminalCountryCode from emv/terminalCountryCode
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'emv' AND de.name = 'terminalCountryCode';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'terminalCountryCode';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);
#
# Insert Global dualCurrency/dccMinValue from modules/dccMinValue
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'modules' AND de.name = 'dccMinValue';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'dccMinValue';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);
#
# Insert Global dualCurrency/dccMaxValue from modules/dccMaxValue
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'modules' AND de.name = 'dccMaxValue';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'dccMaxValue';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);
#
# Insert Global dualCurrency/dccLocalBins from modules/dccLocalBins
SELECT pd.datavalue INTO @value FROM profile_data pd INNER JOIN data_element de on pd.data_element_id = de.data_element_id INNER JOIN data_group dg on de.data_group_id = dg.data_group_id WHERE pd.profile_id = 1 AND dg.name = 'modules' AND de.name = 'dccLocalBins';
SELECT de.data_element_id INTO @dc_data_element_id FROM data_element de WHERE de.data_group_id = @dc_data_group_id AND de.name = 'dccLocalBins';
INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, @dc_data_element_id, @value, 1, NOW(), 'system', NOW(), 'system', 1, 0);