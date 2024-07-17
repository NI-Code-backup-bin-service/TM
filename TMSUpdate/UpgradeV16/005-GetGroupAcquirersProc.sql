--multiline
CREATE PROCEDURE `get_group_acquirers`(IN groupId int)
BEGIN
  select acquirer_profile_id
  from permissiongroup_acquirer pga
  where pga.permissiongroup_id = groupId;
END