--multiline
CREATE PROCEDURE update_dci_scheme_into_profiles(IN profileID INT)
BEGIN
    SET @ErrorMsg = (SELECT CONCAT('Online refund not enabled to give profileID = ',profileID));
    SET @IsOnlineRefundEnabled = (SELECT datavalue FROM profile_data WHERE data_element_id= (SELECT data_element_id FROM data_element WHERE name = 'onlineRefund') AND datavalue = 'true' AND profile_id = profileID);
    IF @IsOnlineRefundEnabled = 'true' THEN
        UPDATE profile_data SET datavalue = '["MASTER","VISA","DINERS"]' where data_element_id= (SELECT data_element_id FROM data_element WHERE name = 'onlineRefundSchemes') AND profile_id = profileID;
    ELSE
        SELECT @ErrorMsg;
    END IF;
END;