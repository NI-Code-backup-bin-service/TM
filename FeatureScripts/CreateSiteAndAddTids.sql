set @siteName = "Feature Script Test Site";
set @merchantID = 654567834567;

call profile_store(-1, 4, @siteName, 1, "Feature Script");
set @profileId = (select profile_id from profile where name = @siteName limit 1);

call site_store(-1, 1, "Feature Script");
set @siteId = (select last_insert_id());

call site_profiles_store(-1, @siteId, @profileId, 1, "Feature Script");
call site_profiles_store(-1, @siteId, 1, 1, "Feature Script"); -- global
call site_profiles_store(-1, @siteId, 2, 1, "Feature Script"); -- NI
call site_profiles_store(-1, @siteId, 1565, 1, "Feature Script"); -- n genius

-- Add data groups here by altering 3rd element
call profile_data_group_store(-1, @profileId, 1, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 2, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 3, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 4, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 5, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 6, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 7, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 8, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 9, -1, "Feature Script");
call profile_data_group_store(-1, @profileId, 10, -1, "Feature Script");

-- Alter 2nd element to change data element and 3rd element for datavalue
call store_profile_data(@profileId, 1, @merchantID, "Feature Script", 1, 1);
call store_profile_data(@profileId, 3, @siteName, "Feature Script", 1, 1);
                                               
-- Add the tids now the site has been made (this can be split into it's own script)
call save_tid(34567457, 765938473625, @siteId);
call save_tid(34567454, 765938473626, @siteId);
call save_tid(34567455, 765938473627, @siteId);