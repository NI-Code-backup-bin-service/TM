--multiline
CREATE PROCEDURE `remove_site_override`(
IN merchantid varchar(50), IN data_group_name Varchar(100), IN element_name varchar(50)
)

-- --------------------------------------------------------------------------------
-- Routine DDL
-- CREATED BY : SURENDER GUSAIN
-- MODDED BY : ABHINANDAN NARAYANAN
-- CREATED ON : 30/11/2022
-- MODDED ON : 13/04/2023
-- PURPOSE	  : Remove Override element for site level only
-- --------------------------------------------------------------------------------
BEGIN
DECLARE profileID, siteId, datagroupid, element_id, approvalid, overr INT;
DECLARE originalvalue, updatedvalue, acquirer mediumtext;


#--get site id from merchant id
SET @siteId=  (SELECT site_id FROM site_profiles WHERE profile_id = (
    SELECT profile_id FROM profile_data WHERE datavalue = merchantid
      AND data_element_id = (
        SELECT data_element_id FROM data_element WHERE name = 'merchantNo'
      )
  ));

SET @profileId = (SELECT
					  sp.profile_id
				  FROM site_profiles sp
						   LEFT JOIN profile p ON p.profile_id = sp.profile_id
						   LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
				  WHERE sp.site_id = @siteId
				  ORDER BY pt.priority
				  LIMIT 1);


#--GET DATA GROUP ID FROM GROUP NAME
SET @datagroupid=(SELECT data_group_id from data_group where displayname_en=data_group_name LIMIT 1);
SET @element_id=(SELECT data_element_id from data_element where displayname_en=element_name and data_group_id=@datagroupid LIMIT 1);
#--GET OVERRIDDEN VALUE FOR THE DATA ELEMENT
SET @overr=(SELECT overriden from profile_data where profile_id = @profileID and data_element_id = @element_id);

#select @siteId,@profileId,@datagroupid,@element_id,@overr;

IF (@siteId IS NULL  OR @profileId=0) THEN
    SELECT 'site or profile not found.' As StatusInfo;
ELSEIF (@element_id IS NULL  OR @element_id=0) THEN
	SELECT 'Data element not found.' As StatusInfo;
ELSEIF (@overr = 0 OR @overr IS NULL) THEN
    SELECT 'Data element cannot be removed as it is not overridden' As StatusInfo;
ELSE
    #--GET ORIGINAL VALUE
    SET @originalvalue =(
    SELECT datavalue from site_data sd
    WHERE sd.site_id = @siteId AND sd.data_element_id = @element_id AND priority > 2 ORDER BY priority
    LIMIT 1);

    #--GET UPDATED VALUE
    SET @updatedvalue =(
    SELECT datavalue from site_data sd
    WHERE sd.site_id = @siteId AND sd.data_element_id = @element_id AND priority = 2
    LIMIT 1);

    #--Get Acquirer Name- As For Site level priority is 4
    SET @acquirer = (select distinct p4.name from profile p
				    LEFT JOIN (site_profiles tp4
				    join profile p4 on p4.profile_id = tp4.profile_id
				    join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id
				    and pt4.priority = 4) on tp4.site_id = @siteId);

    #--Record for history to change details
    INSERT INTO approvals
    (
    profile_id, data_element_id, change_type, current_value, new_value, created_at,
    created_by, approved, acquirer, is_password, is_encrypted
    )
    VALUES
    (
    @profileId, @element_id, 4, @originalvalue, @updatedvalue, NOW(),
    'System', 0, @acquirer, false, false
    );

    SET @approvalid = LAST_INSERT_ID();

    CALL `approve_change`(@approvalid, 'System');

    SELECT 'Site Override removed successfully.-->' As StatusInfo;

END IF;

END
