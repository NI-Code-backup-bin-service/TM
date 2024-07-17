--multiline;
CREATE PROCEDURE `add_wcp_fields_to_site`(
    IN mcc VARCHAR(255),
    IN bcc VARCHAR(255),
    IN phone VARCHAR(255),
    IN wcp_name VARCHAR(255),
    IN email VARCHAR(255),
    IN submerchid VARCHAR(255),
    IN MID VARCHAR(255)
)
BEGIN
    DECLARE wcp_data_group_id INT DEFAULT (SELECT data_group_id FROM data_group WHERE name = 'weChatPay' LIMIT 1);
    DECLARE mid_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'merchantNo' LIMIT 1);

    DECLARE site_profile_id INT DEFAULT (
        SELECT profile_id FROM profile_data
        WHERE 
            data_element_id = mid_data_element_id
            AND datavalue = MID
        LIMIT 1
    );

    /* IDs of the data elements that need to be added to the site */
    DECLARE wcpMcc_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'wcpMcc' LIMIT 1);
    DECLARE wcpBcc_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'wcpBcc' LIMIT 1);
    DECLARE phone_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'wcpContactPhone' LIMIT 1);
    DECLARE name_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'wcpContactName' LIMIT 1);
    DECLARE email_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'wcpContactEmail' LIMIT 1);
    DECLARE submerchid_data_element_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = 'wcpSubMerchantId' LIMIT 1);


    /* Check if the site has the data group set - add it to the site if not */
    SET @data_group_count = (SELECT COUNT(*) FROM profile_data_group WHERE profile_id = site_profile_id AND data_group_id = wcp_data_group_id);

    IF @data_group_count < 1 THEN
        INSERT INTO profile_data_group (profile_id, data_group_id, `version`, updated_at, updated_by, created_at, created_by)
            VALUES (site_profile_id, wcp_data_group_id, 1, NOW(), 'admin', NOW(), 'admin')
            ON DUPLICATE KEY UPDATE updated_at = NOW();
    END IF;


    SET @site_id = (SELECT sp.site_id FROM site_profiles sp WHERE sp.profile_id = site_profile_id LIMIT 1);

    /* Check if the site's chain or acquirer have the weChatPay data group set */
    SET @acquirer_id = (SELECT p.profile_id FROM site_profiles sp
            LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
            LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
            WHERE sp.site_id = @site_id AND pt.name = 'acquirer');

    SET @chain_id = (SELECT p.profile_id FROM site_profiles sp
            LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
            LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
            WHERE sp.site_id = @site_id AND pt.name = 'chain');


    SET @override_count = (SELECT COUNT(*) FROM profile_data_group
            WHERE data_group_id = wcp_data_group_id AND (profile_id = @acquirer_id OR profile_id = @chain_id)
        );

        
    IF @override_count > 0 THEN

        /* Add the wechat data as overriden since it exists at the acquirer or chain */

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, wcpMcc_data_element_id, mcc, 1, NOW(), 'admin', NOW(), 'admin', 1, 1, 0)
            ON DUPLICATE KEY UPDATE datavalue = mcc, updated_at = NOW();
        
        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, wcpBcc_data_element_id, bcc, 1, NOW(), 'admin', NOW(), 'admin', 1, 1, 0)
            ON DUPLICATE KEY UPDATE datavalue = bcc, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, phone_data_element_id, phone, 1, NOW(), 'admin', NOW(), 'admin', 1, 1, 0)
            ON DUPLICATE KEY UPDATE datavalue = phone, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, name_data_element_id, wcp_name, 1, NOW(), 'admin', NOW(), 'admin', 1, 1, 0)
            ON DUPLICATE KEY UPDATE datavalue = wcp_name, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, email_data_element_id, email, 1, NOW(), 'admin', NOW(), 'admin', 1, 1, 0)
            ON DUPLICATE KEY UPDATE datavalue = email, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, submerchid_data_element_id, submerchid, 1, NOW(), 'admin', NOW(), 'admin', 1, 1, 0)
            ON DUPLICATE KEY UPDATE datavalue = submerchid, updated_at = NOW();

    ELSE

        /* Add the wechat data as not overriden since it doesn't exist at acquirer or chain */

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, wcpMcc_data_element_id, mcc, 1, NOW(), 'admin', NOW(), 'admin', 1, 0, 0)
            ON DUPLICATE KEY UPDATE datavalue = mcc, updated_at = NOW();
        
        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, wcpBcc_data_element_id, bcc, 1, NOW(), 'admin', NOW(), 'admin', 1, 0, 0)
            ON DUPLICATE KEY UPDATE datavalue = bcc, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, phone_data_element_id, phone, 1, NOW(), 'admin', NOW(), 'admin', 1, 0, 0)
            ON DUPLICATE KEY UPDATE datavalue = phone, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, name_data_element_id, wcp_name, 1, NOW(), 'admin', NOW(), 'admin', 1, 0, 0)
            ON DUPLICATE KEY UPDATE datavalue = wcp_name, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, email_data_element_id, email, 1, NOW(), 'admin', NOW(), 'admin', 1, 0, 0)
            ON DUPLICATE KEY UPDATE datavalue = email, updated_at = NOW();

        INSERT INTO profile_data (profile_id, data_element_id, datavalue, `version`, updated_at, updated_by, created_at, created_by, approved, overriden, not_overridable)
            VALUES (site_profile_id, submerchid_data_element_id, submerchid, 1, NOW(), 'admin', NOW(), 'admin', 1, 0, 0)
            ON DUPLICATE KEY UPDATE datavalue = submerchid, updated_at = NOW();
    END IF;
END;