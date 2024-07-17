#############
# Migration #
#############
#
# Get the cashback data group and definitions ids into variables to be used
SELECT data_group_id INTO @cashback_datagroup_id FROM data_group WHERE name = 'cashback';
SELECT data_element_id INTO @cashback_definitions_id FROM data_element WHERE name = 'definitions' AND data_group_id = @cashback_datagroup_id;
#
# Enabled the cashback data group for profiles with cashback config set
INSERT INTO profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) select distinct profile_id, @cashback_datagroup_id, 1, NOW(), 'system', NOW(), 'system' FROM cashback;
#
# Migrate all the cashback json into the data elements table to the respective profile
INSERT INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved) select profile_id, @cashback_definitions_id, cashback_data, 1, NOW(), 'system', NOW(), 'system', 1 FROM cashback;
#
#
#
###########
# Cleanup #
###########
#
# Get the cashback permission ID
SELECT permission_id INTO @cashback_user_group_id FROM permission WHERE name = 'Cashback';
#
# Delete the cashback user group rows from user_permissiongroup
DELETE FROM user_permissiongroup WHERE permission_group_id = @cashback_user_group_id;
#
# Delete the permission itself
DELETE FROM permission WHERE permission_id = @cashback_user_group_id;
#
# Finally, drop the cashback table
DROP TABLE cashback;
