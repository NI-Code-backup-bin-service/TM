--multiline;
CREATE PROCEDURE `save_pending_profile_change`(
    IN profile_id int,
    IN change_type INT,
    IN dataValue TEXT,
    IN updated_by varchar(255),
    IN tidId TEXT,
    IN approved int)
BEGIN
    SET @profileType = (SELECT profile_type_id FROM `profile` p WHERE p.profile_id = profile_id);

    if @profileType = 4 then
        SET @siteId = (SELECT site_id FROM site_profiles sp WHERE sp.profile_id = profile_id);
        set @acquirer = (select distinct p4.name from profile p
                                                          LEFT JOIN (site_profiles tp4
            join profile p4 on p4.profile_id = tp4.profile_id
            join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = @siteId);
    elseif @profileType = 5 then
        set @acquirer = (select distinct p4.name from profile p
                                                          join tid_site ts on ts.tid_profile_id = profile_id
                                                          join site t on t.site_id = ts.site_id
                                                          JOIN (site_profiles tp4
            join profile p4 on p4.profile_id = tp4.profile_id
            join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = t.site_id
                         WHERE p.profile_type_id = (select profile_type_id from profile_type where profile_type.name = "tid"));					
    end if;

    INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved, created_by, approved_by, approved_at, tid_id, acquirer)
    VALUES (profile_id, 1, change_type, current_value, dataValue, NOW(), approved, updated_by, updated_by, CASE WHEN approved > 0 THEN NOW() ELSE NULL END, tidId, @acquirer);
END