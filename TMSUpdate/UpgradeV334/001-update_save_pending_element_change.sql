--multiline;
CREATE PROCEDURE `save_pending_element_change`(
    IN profile_id int,
    IN data_element_id int,
    IN change_type INT,
    IN dataValue MEDIUMTEXT,
    IN updated_by varchar(255),
    IN is_password BOOLEAN,
    IN is_encrypted BOOLEAN,
    IN currentValue TEXT,
    IN data_element_id_user_fraud INT
)
BEGIN
    DECLARE current_value MEDIUMTEXT;
    DECLARE profileType INT;
    SET profileType = (SELECT profile_type_id FROM `profile` p WHERE p.profile_id = profile_id);
    IF profileType = 4 THEN -- Handle sites
        SET @siteId = (SELECT site_id FROM site_profiles sp WHERE sp.profile_id = profile_id);
        SET current_value = (SELECT sd.datavalue from site_data sd
                             WHERE sd.data_element_id = data_element_id AND
                                     sd.site_id = @siteId
                             ORDER BY sd.priority
                             LIMIT 1);
    ELSEIF profileType = 5 THEN
        SET current_value = (SELECT pd.datavalue FROM profile_data pd
                             WHERE pd.profile_id = profile_id
                               AND pd.data_element_id = data_element_id
                             ORDER BY pd.version DESC
                             LIMIT 1);
    ELSEIF profileType = 3 THEN
        SET current_value = (SELECT cd.datavalue from chain_data cd
                             WHERE cd.data_element_id = data_element_id AND
                                 cd.profile_id = profile_id
                             LIMIT 1);
    ELSE
        SET current_value = (SELECT pd.datavalue FROM profile_data pd
                             WHERE pd.profile_id = profile_id
                               AND pd.data_element_id = data_element_id
                             ORDER BY pd.version DESC
                             LIMIT 1);
    END IF;

    if profileType = 2 then
        set @acquirer = (select p.name from profile p where p.profile_id = profile_id);
    elseif profileType = 3 then
        set @acquirer = (select distinct p2.name from profile p
                                                          join chain_profiles cp on cp.chain_profile_id = profile_id
                                                          join profile p2 on p2.profile_id = cp.acquirer_id);
    elseif profileType = 4 then
        set @acquirer = (select distinct p4.name from profile p
                                                          LEFT JOIN (site_profiles tp4
            join profile p4 on p4.profile_id = tp4.profile_id
            join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = @siteId);
    elseif profileType = 5 then
        set @acquirer = (select distinct p4.name from profile p
                                                          join tid_site ts on ts.tid_profile_id = profile_id
                                                          join site t on t.site_id = ts.site_id
                                                          JOIN (site_profiles tp4
            join profile p4 on p4.profile_id = tp4.profile_id
            join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = t.site_id
                         WHERE p.profile_type_id = (select profile_type_id from profile_type where profile_type.name = "tid"));
    end if;

    IF (current_value = "" OR isnull(current_value)) AND data_element_id_user_fraud = data_element_id THEN -- Handle users/fraud previous value
		SET current_value = currentValue;
    END IF;

    INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, created_by, approved, acquirer, is_password, is_encrypted)
    VALUES (profile_id,data_element_id,change_type, current_value, dataValue,NOW(),updated_by,0, @acquirer, is_password, is_encrypted);
END