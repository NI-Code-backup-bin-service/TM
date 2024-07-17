--multiline;
CREATE PROCEDURE `save_pending_element_change`(
 IN profile_id int,
 IN data_element_id int,
 IN change_type INT,
 IN dataValue varchar(255),
 IN updated_by varchar(255))
BEGIN
	DECLARE current_value VARCHAR(256);
    DECLARE profileType INT;
    DECLARE siteID INT;
    
    SET profileType = (SELECT profile_type_id FROM `profile` p WHERE p.profile_id = profile_id);
    
    IF profileType = 4 THEN -- Handle sites
		SET @siteId = (SELECT site_id FROM site_profiles sp WHERE sp.profile_id = profile_id);
        SET current_value = (SELECT sd.datavalue from site_data sd					
							 WHERE sd.data_element_id = data_element_id AND 
                             sd.site_id = @siteId
                             ORDER BY sd.priority
                             LIMIT 1);
        
	ELSE
		SET current_value = (SELECT pd.datavalue FROM profile_data pd 
								WHERE pd.profile_id = profile_id 
                                AND pd.data_element_id = data_element_id
                                ORDER BY pd.version DESC
                                LIMIT 1);
    END IF;
    
	INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, created_by, approved) 
    VALUES (profile_id,data_element_id,change_type, current_value, dataValue,NOW(),updated_by,0);
END