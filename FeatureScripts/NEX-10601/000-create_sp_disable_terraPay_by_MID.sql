--multiline
CREATE PROCEDURE `disable_terraPay_by_MID`(in MID text, provider text)
BEGIN
    # Check if merchantID exists.
    SET @merchantExists = EXISTS(Select profile_id from profile_data pd where pd.data_element_id = (Select data_element_id from data_element where name = "merchantNo") AND
            pd.datavalue = MID);

    IF @merchantExists = 1 THEN
        SET @profile_id = (Select profile_id from profile_data pd where pd.data_element_id = (Select data_element_id from data_element where name = "merchantNo") AND
                            pd.datavalue = MID);

        SET @site_id = (Select site_id from site_profiles where profile_id = @profile_id);

        # Check if active modules exists on the site already as an override otherwise we need to override it first
        SET @overrideExists = EXISTS(Select * FROM profile_data
            WHERE profile_id = @profile_id AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active" AND data_group_id = (select data_group_id from data_group where name = 'modules')));

        IF @overrideExists = 1 THEN
            # set the override initially to the same as site so we don't lose any default modules
            SET @activeModulesValue = (Select datavalue FROM profile_data
                WHERE profile_id = @profile_id AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active" AND data_group_id = (select data_group_id from data_group where name = 'modules')));
        END IF;

        # Update the Site Active Modules to remove terraPay.
        set @providervalue:=replace(@activeModulesValue,provider,'');
        UPDATE profile_data pd
        SET datavalue = @providervalue, updated_at = CURRENT_TIMESTAMP, updated_by = "system"
        WHERE pd.profile_id = @profile_id
        AND pd.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active" AND data_group_id = (select data_group_id from data_group where name = 'modules'))
        AND pd.datavalue LIKE CONCAT('%', "terraPay", '%');

        # Find any TID overrides and remove terraPay from the active modules
        UPDATE profile_data pd
        JOIN tid_site ts on ts.tid_profile_id = pd.profile_id
        SET pd.datavalue = REPLACE(pd.datavalue,provider, ''), pd.updated_at = CURRENT_TIMESTAMP, pd.updated_by = "system"
        WHERE pd.profile_id = ts.tid_profile_id
        AND pd.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active" AND data_group_id = (select data_group_id from data_group where name = 'modules'))
        AND ts.site_id = @site_id AND ts.tid_profile_id IS NOT NULL
        AND pd.datavalue LIKE CONCAT('%', "terraPay", '%');
    END IF;
END