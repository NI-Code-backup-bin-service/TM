--multiline;
CREATE PROCEDURE save_tid_to_site_and_pending_profile_change(
    IN profile_id INT,
    IN change_type INT,
    IN dataValue MEDIUMTEXT,
    IN updated_by varchar(255),
    IN tidId INT,
    IN approved INT,
	IN serialNumber varchar(50), 
	IN site INT,  
	IN eodAuto BOOLEAN, 
	IN autoTime varchar(45)
)
BEGIN
    SET @profileType = (SELECT profile_type_id FROM `profile` p WHERE p.profile_id = profile_id);
	INSERT INTO tid (tid_id, serial, flag_status,flagged_date, eod_auto, auto_time) VALUES (tidId, serialNumber, 1, now(), eodAuto, autoTime);
    INSERT INTO tid_site VALUES (tidId, site, NULL, NOW());
    IF @profileType = 4 THEN
        SET @siteId = (SELECT site_id FROM site_profiles sp WHERE sp.profile_id = profile_id);
        SET @acquirer = (SELECT DISTINCT p4.name FROM site_profiles tp4
            JOIN profile p4 ON p4.profile_id = tp4.profile_id
            JOIN profile_type pt4 ON pt4.profile_type_id = p4.profile_type_id AND pt4.priority = 4 WHERE tp4.site_id = @siteId);
    ELSEIF @profileType = 5 THEN
        SET @acquirer = (SELECT DISTINCT p4.name FROM profile p JOIN tid_site ts ON ts.tid_profile_id = profile_id
		  JOIN site t ON t.site_id = ts.site_id JOIN (site_profiles tp4
          JOIN profile p4 ON p4.profile_id = tp4.profile_id
          JOIN profile_type pt4 ON pt4.profile_type_id = p4.profile_type_id AND pt4.priority = 4) ON tp4.site_id = t.site_id
          WHERE p.profile_type_id = (SELECT profile_type_id FROM profile_type WHERE profile_type.name = "tid"));
	END IF;
		SET @createChangeType = (SELECT approval_type_id FROM approval_type WHERE approval_type_name = "Create" LIMIT 1);
		SET @unapprovedRequestCount = (SELECT COUNT(*) FROM approvals a WHERE a.profile_id = profile_id AND a.change_type = change_type AND a.new_value = dataValue AND a.approved_at IS NULL AND a.approved = approved AND a.tid_id = tidId);
		IF  @unapprovedRequestCount < 1 OR change_type = @createChangeType THEN
			INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved, created_by, approved_by, approved_at, tid_id, acquirer)
			VALUES (profile_id, 1, change_type, current_value, dataValue, NOW(), approved, updated_by, updated_by, CASE WHEN approved > 0 THEN NOW() ELSE NULL END, tidId, @acquirer);
	END IF;
END;