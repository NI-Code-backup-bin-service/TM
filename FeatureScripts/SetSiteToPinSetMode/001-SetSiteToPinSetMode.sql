--multiline
CREATE PROCEDURE `set_site_pinset` (in profileId int)
BEGIN
    set @modeElementId = (select de.data_element_id from data_element de inner join data_group dg on dg.name = "modules" where de.name = "mode" limit 1);
    set @eodAutoElementId = (select de.data_element_id from data_element de inner join data_group dg on dg.name = "endOfDay"where de.name = "auto" limit 1);
    set @eodHardElementId = (select de.data_element_id from data_element de inner join data_group dg on dg.name = "endOfDay"where de.name = "hardLimit" limit 1);
    set @eodSoftElementId = (select de.data_element_id from data_element de inner join data_group dg on dg.name = "endOfDay"where de.name = "softLimit" limit 1);

    if (not exists (select * from profile_data where profile_id = profileId and data_element_id = @modeElementId)) then insert into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) values (profileId, @modeElementId, "pinset", 1, CURDATE(), "admin", CURDATE(), "admin", 1, 0);
    else update profile_data set datavalue = "pinset" where data_element_id = @modeElementId and profile_id = profileId;
    end if;

    if (not exists (select * from profile_data where profile_id = profileId and data_element_id = @eodAutoElementId)) then insert into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) values (profileId, @eodAutoElementId, "false", 1, CURDATE(), "admin", CURDATE(), "admin", 1, 1);
    else update profile_data set datavalue = "false" where data_element_id = @eodAutoElementId and profile_id = profileId;
    end if;

    if (not exists (select * from profile_data where profile_id = profileId and data_element_id = @eodHardElementId)) then insert into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) values (profileId, @eodHardElementId, "0", 1, CURDATE(), "admin", CURDATE(), "admin", 1, 1);
    else update profile_data set datavalue = "0" where data_element_id = @eodHardElementId and profile_id = profileId;
    end if;

    if (not exists (select * from profile_data where profile_id = profileId and data_element_id = @eodSoftElementId)) then insert into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) values (profileId, @eodSoftElementId, "0", 1, CURDATE(), "admin", CURDATE(), "admin", 1, 1);
    else update profile_data set datavalue = "0" where data_element_id = @eodSoftElementId and profile_id = profileId;
    end if;
END