--multiline;
CREATE PROCEDURE `record_site_to_history`(
    IN profile_id int,
    IN change_type INT,
    IN dataValue MEDIUMTEXT,
    IN updated_by varchar(255),
    IN approved int,
    IN merchantId varchar(45))
BEGIN
    SET @siteId = (SELECT site_id FROM site_profiles sp WHERE sp.profile_id = profile_id);
    set @acquirer = (select distinct p4.name from profile p
                                                      LEFT JOIN (site_profiles tp4
        join profile p4 on p4.profile_id = tp4.profile_id
        join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = @siteId);
    INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved_at, created_by, approved_by, approved, tid_id, merchant_id, acquirer)
    VALUES (profile_id, 1, change_type, current_value, dataValue, NOW(), NOW(), updated_by, updated_by, approved, null, merchantId, @acquirer);
END