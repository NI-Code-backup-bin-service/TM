--multiline
CREATE PROCEDURE `get_operations_group_acquirers`(IN groupId int)
BEGIN
  select acquirer_profile_id
  from operations_permissiongroup_acquirer pga
  where pga.permissiongroup_id = groupId;
END