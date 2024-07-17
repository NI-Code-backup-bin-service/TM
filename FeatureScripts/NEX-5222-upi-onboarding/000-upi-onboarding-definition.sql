--multiline
CREATE PROCEDURE temp_upiOnboarding()
BEGIN
    DROP TEMPORARY TABLE IF EXISTS upi_population_acquirers;
    CREATE TEMPORARY TABLE upi_population_acquirers (acquirerId INT PRIMARY KEY, acquirerIIN TEXT, forwardingIIN TEXT, countryCode TEXT);
    INSERT INTO upi_population_acquirers (acquirerId, acquirerIIN, forwardingIIN, countryCode)
    VALUES (-1234, 'acqIIN_1', 'forwardingIIN_1', 'countryCode_1'),
           (-5678, 'acqIIN_2', 'forwardingIIN_2', 'countryCode_2');

    DROP TEMPORARY TABLE IF EXISTS upi_population_mids;
    CREATE TEMPORARY TABLE upi_population_mids (
        -- These columns are to be filled in manually
        mid VARCHAR(255) PRIMARY KEY, categoryCode TEXT,
        -- These columns only to be used internally by this script
        siteProfileId INT UNIQUE, siteId INT UNIQUE);
    INSERT INTO upi_population_mids (mid, categoryCode)
    VALUES ('mid1', 'CatCode1'),
           ('mid2', 'CatCode1');

    SET @siteProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'site');
    SET @chainProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'chain');
    SET @acquirerProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'acquirer');
    SET @globalProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'global');

    SET @upiDataGroupId = (select data_group_id from data_group where `name` = 'upiQr');
    SET @storeDataGroupId = (select data_group_id from data_group where `name` = 'store');
    SET @modulesDataGroupId = (select data_group_id from data_group where `name` = 'modules');

    SET @merchantNoDeId = (select data_element_id from data_element where data_group_id = @storeDataGroupId AND `name` = 'merchantNo');
    SET @acquirerIinDeId = (select data_element_id from data_element where data_group_id = @upiDataGroupId AND `name` = 'acquirerIIN');
    SET @forwardingIinDeId = (select data_element_id from data_element where data_group_id = @upiDataGroupId AND `name` = 'forwardingIIN');
    SET @countryCodeDeId = (select data_element_id from data_element where data_group_id = @upiDataGroupId AND `name` = 'countryCode');
    SET @categoryCodeDeId = (select data_element_id from data_element where data_group_id = @upiDataGroupId AND `name` = 'categoryCode');
    SET @activeModulesDeId = (select data_element_id from data_element where data_group_id = @modulesDataGroupId AND `name` = 'active');

    SET @chainPriority = (select priority from profile_type where name = 'chain');

    -- Populate the site profile ID for each site in our temp table
    update upi_population_mids upm
    inner join profile_data pd on pd.datavalue = upm.mid AND pd.data_element_id = @merchantNoDeId
    set siteProfileId = pd.profile_id;

    -- Populate the site ID for each site in our temp table
    update upi_population_mids upm
    inner join site_profiles sp on sp.profile_id = upm.siteProfileId
    set upm.siteId = sp.site_id;

    -- Enable the UPI data group ('upiQr') for each given acquirer
    insert ignore into profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by)
    select p.profile_id, @upiDataGroupId, 1, NOW(), 'system', NOW(), 'system'
    from profile p
    inner join profile_type pt on p.profile_type_id = pt.profile_type_id
    inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    and pt.name = 'acquirer';

    -- Populate/update the acquirerIIN value for each acquirer
    update profile_data pd
    inner join profile p on pd.profile_id = p.profile_id
    inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    set pd.datavalue = upa.acquirerIIN,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @acquirerProfileTypeId
      and pd.data_element_id = @acquirerIinDeId
      and pd.datavalue != upa.acquirerIIN;

    insert ignore into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    select p.profile_id, @acquirerIinDeId, upa.acquirerIIN, 1, NOW(), 'system', NOW(), 'system', 1, 0, 0
    from profile p
    inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    where p.profile_type_id = @acquirerProfileTypeId;

    -- Populate/update the forwardingIIN value for each acquirer
    update profile_data pd
        inner join profile p on pd.profile_id = p.profile_id
        inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    set pd.datavalue = upa.forwardingIIN,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @acquirerProfileTypeId
      and pd.data_element_id = @forwardingIinDeId
      and pd.datavalue != upa.forwardingIIN;

    insert ignore into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    select p.profile_id, @forwardingIinDeId, upa.forwardingIIN, 1, NOW(), 'system', NOW(), 'system', 1, 0, 0
    from profile p
    inner join profile_type pt on p.profile_type_id = pt.profile_type_id
    inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    where p.profile_type_id = @acquirerProfileTypeId;

    -- Populate/update the countryCode value for each acquirer
    update profile_data pd
        inner join profile p on pd.profile_id = p.profile_id
        inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    set pd.datavalue = upa.countryCode,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @acquirerProfileTypeId
      and pd.data_element_id = @countryCodeDeId
      and pd.datavalue != upa.countryCode;

    insert ignore into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    select p.profile_id, @countryCodeDeId, upa.countryCode, 1, NOW(), 'system', NOW(), 'system', 1, 0, 0
    from profile p
    inner join upi_population_acquirers upa on upa.acquirerId = p.profile_id
    where p.profile_type_id = @acquirerProfileTypeId;

    -- Populate/update the categoryCode value for each MID
    update profile_data pd
        inner join profile p on pd.profile_id = p.profile_id
        inner join upi_population_mids upm on upm.siteProfileId = p.profile_id
    set pd.datavalue = upm.categoryCode,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @siteProfileTypeId
      and pd.data_element_id = @categoryCodeDeId
      and pd.datavalue != upm.categoryCode;

    insert ignore into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    select p.profile_id, @categoryCodeDeId, upm.categoryCode, 1, NOW(), 'system', NOW(), 'system', 1, 0, 0
    from profile p
    inner join upi_population_mids upm on upm.siteProfileId = p.profile_id
    where p.profile_type_id = @siteProfileTypeId;

    -- Update any site-level overrides for acquirer IIN
    update profile_data pd
    inner join profile p on pd.profile_id = p.profile_id
    inner join upi_population_mids upm on upm.siteProfileId = p.profile_id
    inner join site_profiles sp on sp.site_id = upm.siteId
    inner join profile p_acq on sp.profile_id = p_acq.profile_id and p_acq.profile_type_id = @acquirerProfileTypeId
    inner join upi_population_acquirers upa on upa.acquirerId = p_acq.profile_id
    set pd.datavalue = upa.acquirerIIN,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @siteProfileTypeId
      and pd.data_element_id = @acquirerIinDeId
      and pd.datavalue != upa.acquirerIIN;

-- Update any site-level overrides for forwarding IIN
    update profile_data pd
    inner join profile p on pd.profile_id = p.profile_id
    inner join upi_population_mids upm on upm.siteProfileId = p.profile_id
    inner join site_profiles sp on sp.site_id = upm.siteId
    inner join profile p_acq on sp.profile_id = p_acq.profile_id and p_acq.profile_type_id = @acquirerProfileTypeId
    inner join upi_population_acquirers upa on upa.acquirerId = p_acq.profile_id
    set pd.datavalue = upa.forwardingIIN,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @siteProfileTypeId
      and pd.data_element_id = @forwardingIinDeId
      and pd.datavalue != upa.forwardingIIN;

    -- Update any site-level overrides for country code
    update profile_data pd
        inner join profile p on pd.profile_id = p.profile_id
        inner join upi_population_mids upm on upm.siteProfileId = p.profile_id
        inner join site_profiles sp on sp.site_id = upm.siteId
        inner join profile p_acq on sp.profile_id = p_acq.profile_id and p_acq.profile_type_id = @acquirerProfileTypeId
        inner join upi_population_acquirers upa on upa.acquirerId = p_acq.profile_id
    set pd.datavalue = upa.countryCode,
        pd.version = pd.version + 1,
        pd.updated_at = NOW(),
        pd.updated_by = 'system'
    where p.profile_type_id = @siteProfileTypeId
      and pd.data_element_id = @countryCodeDeId
      and pd.datavalue != upa.countryCode;

    -- Need to append 'upi' module to all the sites.
    -- First, need to add 'upi' to the active modules where it is directly set at the site level.
    UPDATE profile_data pd
    JOIN profile p on pd.profile_id = p.profile_id AND pd.data_element_id = @activeModulesDeId
    JOIN profile_type pt on p.profile_type_id = pt.profile_type_id
    JOIN upi_population_mids upm ON p.profile_id = upm.siteProfileId
    SET pd.datavalue = REPLACE(pd.datavalue, ']', ',"upi"]')
    WHERE pd.datavalue NOT LIKE '%upi%'
    AND pt.profile_type_id = @siteProfileTypeId;
    -- Next, need to add 'upi' as a new entry at site level where it is using chain/acquirer/global values
    INSERT IGNORE INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    SELECT upm.siteProfileId, @activeModulesDeId, REPLACE(sd.datavalue, ']', ',"upi"]'), 1, NOW(), 'system', NOW(), 'system', 1, 0, 0
    FROM NextGen_TMS.site_data sd
    INNER JOIN (
        SELECT site_id, min(priority) as 'priority'
        FROM site_data
        WHERE data_element_id = @activeModulesDeId
        GROUP BY site_id) sdm -- sdm short for 'site data min'
        ON sd.site_id = sdm.site_id and sd.priority = sdm.priority
    INNER JOIN upi_population_mids upm on sd.site_id = upm.siteId
    WHERE data_element_id = @activeModulesDeId
      AND sd.priority = sdm.priority
      AND sd.priority >= @chainPriority;

    -- Add UPI module to TIDs
    update profile_data pd
    inner join profile p on pd.profile_id = p.profile_id and pd.data_element_id = @activeModulesDeId
    inner join tid_site ts on ts.tid_profile_id = p.profile_id
    inner join upi_population_mids upm on upm.siteId = ts.site_id
    SET pd.datavalue = REPLACE(pd.datavalue, ']', ',"upi"]')
    WHERE pd.datavalue NOT LIKE '%upi%';

    -- Update all updated times to ensure delta checks work and configs get downloaded
    UPDATE profile_data pd
    SET pd.updated_at = NOW()
    WHERE pd.datavalue LIKE '%upi%'
    AND pd.data_element_id = (SELECT data_element_id FROM data_element de WHERE de.name = 'active');

    DROP TEMPORARY TABLE IF EXISTS upi_population_acquirers;
    DROP TEMPORARY TABLE IF EXISTS upi_population_mids;
END
