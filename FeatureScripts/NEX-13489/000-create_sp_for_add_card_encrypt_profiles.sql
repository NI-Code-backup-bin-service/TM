--multiline
CREATE PROCEDURE update_card_encrypt_into_profiles(IN profileID INT)
BEGIN
    UPDATE profile_data SET datavalue = '["ENOC"]' where data_element_id= (SELECT data_element_id FROM data_element WHERE name = 'encryptionEnabled' and data_group_id = (select data_group_id from data_group where name = 'core')) AND profile_id = profileID;
END;