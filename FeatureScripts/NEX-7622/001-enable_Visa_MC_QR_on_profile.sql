--multiline
CREATE PROCEDURE `enable_Visa_MC_QR_on_profile`(in MID text, MPAN text, CatergoryCode text,	dataGroup text)
BEGIN
    SET @dataGroupName = dataGroup;
    SET @MPANName = "mpan";
    SET @categoryCodeName = "categoryCode";
    SET @activeName = "active";
	
    SET @profile_id = (Select profile_id from profile_data pd where pd.data_element_id = (Select data_element_id from data_element where name = "merchantNo") AND 
						pd.datavalue = MID);
    
    SET @site_id = (Select site_id from site_profiles where profile_id = @profile_id);
    
    SET @chain_id = (SELECT p.profile_id FROM site_profiles sp 
						LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
						LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
						WHERE sp.site_id = @site_id and pt.name = "chain");
                        
	SET @acquirer_id = (SELECT p.profile_id FROM site_profiles sp 
						LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
						LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
						WHERE sp.site_id = @site_id and pt.name = "acquirer");
    
    SET @dataGroupExists := EXISTS(Select * FROM profile_data_group 
		WHERE profile_id = @profile_id AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = @dataGroupName));
    
    # Enable the desired QR DataGroup
    IF @dataGroupExists = 0 THEN
		INSERT ignore into profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) 
		values (@profile_id, (SELECT data_group_id from data_group where name = @dataGroupName), 1, NOW(), 'NISuper', NOW(), 'NISuper');
    END IF;
    
    # Set the MPAN value 
    INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) 
		values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @MPANName AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = @dataGroupName)), MPAN, 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) 
	ON DUPLICATE KEY 
		UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
        
	# Set the CategoryCode value
	INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) 
		values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @categoryCodeName AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = @dataGroupName)), CatergoryCode, 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) 
	ON DUPLICATE KEY 
		UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
    
	# Check if active modules exists on the site already as an override otherwise we need to override it first
    SET @overrideExists := EXISTS(Select * FROM profile_data 
		WHERE profile_id = @profile_id AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active"));
        
	IF @overrideExists = 0 THEN
		# set the override initially to the same as chain / acquirer so we don't lose any default modules
        SET @activeModulesValue := (Select datavalue FROM profile_data 
			WHERE profile_id = @chain_id AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active"));
            
		# If there isn't any active modules set at chain check acquirer
		IF @activeModulesValue IS NULL THEN
			SET @activeModulesValue := (Select datavalue FROM profile_data 
				WHERE profile_id = @acquirer_id AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active"));
		END IF;
        
		# add the override
		INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) 
			values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = "active" AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = "modules")), @activeModulesValue, 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) 
		ON DUPLICATE KEY 
			UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);     
	END IF;
    
	# Update the Site Active Modules to add selected QR module
    SET @datavalue = CONCAT(',\"', @dataGroupName, '\"]');
    UPDATE profile_data pd
	SET datavalue = REPLACE(pd.datavalue, "]", @datavalue), updated_at = CURRENT_TIMESTAMP, updated_by = "system"
	WHERE pd.profile_id = @profile_id
	AND pd.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active")
	AND datavalue NOT LIKE CONCAT('%', @dataGroupName, '%');
    
    # Find any TID overrides and add QR to the active modules
	UPDATE profile_data pd
    JOIN tid_site ts on ts.tid_profile_id = pd.profile_id
	SET pd.datavalue = REPLACE(pd.datavalue, "]", @datavalue), pd.updated_at = CURRENT_TIMESTAMP, pd.updated_by = "system"    
    WHERE pd.profile_id = ts.tid_profile_id 
	AND pd.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "active")
    AND ts.site_id = @site_id AND ts.tid_profile_id IS NOT NULL
	AND pd.datavalue NOT LIKE CONCAT('%', @dataGroupName, '%');
END