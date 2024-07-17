--multiline
CREATE PROCEDURE `Enable_IPP_DG_Update_DE_on_profile`(in dataGroupName text, in softUiDataGroupName text, MID text, bankUserId text, paymentLocationId text,	merchantTag text, participantBankCode text, participantGroupCode text, mainMenuConfiguration text)
BEGIN
    SET @BankUserID = "bankUserId";
    SET @PaymentLocationID = "paymentLocationId";
    SET @MerchantTag = "merchantTag";
    SET @ParticipantBankCode = "participantBankCode";
    SET @ParticipantGroupCode = "participantGroupCode";
    SET @MainMenuConfiguration = "mainMenuConfiguration";

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
    
    SET @dataGroupExists := EXISTS(Select * FROM profile_data_group WHERE profile_id = @profile_id AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName));
    
    SET @SoftUiDataGroupExists := EXISTS(Select * FROM profile_data_group WHERE profile_id = @profile_id AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = softUiDataGroupName));
    
    # Enable the desired DataGroup
    IF @dataGroupExists = 0 THEN
		INSERT ignore into profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) 
		values (@profile_id, (SELECT data_group_id from data_group where name = dataGroupName), 1, NOW(), 'system', NOW(), 'system');
    END IF;
    
    # Set the bankUserId value
    INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @BankUserID AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), bankUserId, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
    ON DUPLICATE KEY
        UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
        
	# Set the PaymentLocationID value
	INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) 
		values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @PaymentLocationID AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), paymentLocationId, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
	ON DUPLICATE KEY 
		UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

    # Set the MerchantTag value
    INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @MerchantTag AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), merchantTag, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
    ON DUPLICATE KEY
        UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

    # Set the ParticipantBankCode value
    INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @ParticipantBankCode AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), participantBankCode, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
    ON DUPLICATE KEY
        UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

    # Set the ParticipantGroupCode value
    INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @ParticipantGroupCode AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), participantGroupCode, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
    ON DUPLICATE KEY
        UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

    # Enable the softUi DataGroup
    IF @SoftUiDataGroupExists = 0 THEN
		INSERT ignore into profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) 
		values (@profile_id, (SELECT data_group_id from data_group where name = softUiDataGroupName), 1, NOW(), 'system', NOW(), 'system');
    END IF;

    # Set the MainMenuConfiguration value
    INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @MainMenuConfiguration AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = softUiDataGroupName)), mainMenuConfiguration, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
    ON DUPLICATE KEY
        UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
END