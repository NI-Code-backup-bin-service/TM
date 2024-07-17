-- --multiline;
CREATE PROCEDURE `update_or_set_thirdPartyPackageList_and_insert_into_approval`(IN profileId INT,  IN oldValue varchar(255),IN newValue varchar(255),IN updatedBy varchar(255), tidId INT,IN changeType INT)
BEGIN
SET @dataElementID=(SELECT data_element_id from data_element where name='thirdPartyPackageList' and data_group_id=(select data_group_id from data_group where name='thirdParty'));
SET @siteId=(SELECT site_id from tid_site WHERE tid_id =tidId);
SET @acqName=(select distinct p4.name from profile p
				LEFT JOIN (site_profiles tp4
				join profile p4 on p4.profile_id = tp4.profile_id
				join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id 
				and pt4.priority = 4) on tp4.site_id = @siteId);
SET @dataValue =(SELECT datavalue FROM profile_data WHERE data_element_id = @dataElementID AND profile_id = profileId);
	IF @dataValue IS NOT NULL THEN
		UPDATE profile_data SET datavalue = newValue, updated_at = NOW(), updated_by = updatedBy WHERE data_element_id = @dataElementID AND profile_id = profileId;
	ELSE
		INSERT INTO profile_data 
    							(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden)
    							VALUES
    							(profileId,@dataElementID, newValue, 1, NOW(), updatedBy, NOW(), updatedBy, 1, 0);
END IF;
INSERT INTO approvals (profile_id,data_element_id, change_type, current_value, new_value, created_at,approved_at, approved, created_by,approved_by, tid_id,acquirer)
				   VALUE
				   (profileId,@dataElementID, changeType, oldValue, newValue, NOW(),NOW(), 1, updatedBy,updatedBy, tidId,@acqName);
END