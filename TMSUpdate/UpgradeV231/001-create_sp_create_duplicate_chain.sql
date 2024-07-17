--multiline
CREATE PROCEDURE `save_duplicate_chain_profile`(IN profileId int,IN chainProfileId int,IN createdBy varchar(255))
BEGIN
    INSERT INTO profile_data (profile_id,approved,created_by,data_element_id,datavalue,is_encrypted,not_overridable,overriden,updated_by,version)  SELECT profileId,approved,createdBy,data_element_id,datavalue,is_encrypted,not_overridable,overriden,createdBy,version FROM profile_data WHERE profile_id = chainProfileId;
    INSERT INTO profile_data_group(profile_id,data_group_id,version,updated_by,created_by) Select profileId,data_group_id,1,createdBy,createdBy from profile_data_group where profile_id= chainProfileId;
END