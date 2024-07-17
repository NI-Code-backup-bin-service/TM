--multiline
CREATE PROCEDURE `create_new_site`(
IN merchantId varchar(255),
IN siteName varchar(255),
IN addressLine1 varchar(255),
IN addressLine2 varchar(255),
IN acquirerId int,
IN chainId int,
IN useOldSiteStore bool
)
this_proc:BEGIN
DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
set @merchantIdTest = (SELECT profile_data_id FROM profile_data where data_element_id = 1 AND datavalue = @merchantId);
set @acquirerProfileType = (SELECT pt.name	from profile_type pt left join profile p on p.profile_type_id = pt.profile_type_id WHERE p.profile_id = @acquirerId);
set @chainProfileType = (SELECT pt.name	from profile_type pt left join profile p on p.profile_type_id = pt.profile_type_id WHERE p.profile_id = @chainId);
IF (@merchantIdTest is not null) THEN
	SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'Merchant ID provided already exists';
    LEAVE this_proc;
END IF;
IF (@acquirerProfileType != "acquirer") THEN
	SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'Acquirer ID provided is not a valid Acquirer';
	LEAVE this_proc;
END IF;
IF (@chainProfileType != "chain") THEN
	SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'Chain ID provided is not a valid Chain';
	LEAVE this_proc;
END IF;
call profile_store(-1, 4, @siteName, 1, "Feature Script");
set @profileId = (select profile_id from profile where name = @siteName limit 1);
IF useOldSiteStore = true THEN
	call site_store(-1, @siteName, 1, "Feature Script");
ELSE
	call site_store(-1, 1, "Feature Script");
END IF;
set @siteId = (select last_insert_id());
call site_profiles_store(-1, @siteId, @profileId, 1, "Feature Script");
call site_profiles_store(-1, @siteId, 1, 1, "Feature Script"); -- Add Global
call site_profiles_store(-1, @siteId, @acquirerId, 1, "Feature Script"); -- Add Acquirer ID]
call site_profiles_store(-1, @siteId, @chainId, 1, "Feature Script"); -- Add Chain ID
 -- Add data groups here by altering 3rd element
call profile_data_group_store(-1, @profileId, 1, -1, "Feature Script"); -- Store
call profile_data_group_store(-1, @profileId, 2, -1, "Feature Script"); -- Logging
call profile_data_group_store(-1, @profileId, 3, -1, "Feature Script"); -- Receipt
call profile_data_group_store(-1, @profileId, 4, -1, "Feature Script"); -- Reversal
call profile_data_group_store(-1, @profileId, 5, -1, "Feature Script"); -- EOD
call profile_data_group_store(-1, @profileId, 6, -1, "Feature Script"); -- EMV
call profile_data_group_store(-1, @profileId, 7, -1, "Feature Script"); -- Modules
call profile_data_group_store(-1, @profileId, 8, -1, "Feature Script"); -- User Management
call profile_data_group_store(-1, @profileId, 9, -1, "Feature Script"); -- Core
call profile_data_group_store(-1, @profileId, 10, -1, "Feature Script"); -- Alipay
-- Alter 2nd element to change data element and 3rd element for datavalue
call store_profile_data(@profileId, 1, merchantId, "Feature Script", 1, 1); -- Merchant ID
call store_profile_data(@profileId, 3, @siteName, "Feature Script", 1, 1); -- Site Name
call store_profile_data(@profileId, 4, addressLine1, "Feature Script", 1, 1); -- Address Line 1
call store_profile_data(@profileId, 5, addressLine2, "Feature Script", 1, 1); -- Address Line 2
END