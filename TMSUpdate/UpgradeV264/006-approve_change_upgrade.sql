--multiline;
CREATE PROCEDURE `approve_change`(IN approval_id int, IN approval_user varchar(256))
BEGIN
    DECLARE approval_type,globaloverrideval,dataloc int;
    DECLARE newVal MEDIUMTEXT;
    DECLARE profType varchar(30);
    DECLARE oldVal MEDIUMTEXT;
    DECLARE profileId INT;
    DECLARE elementId INT;
    DECLARE is_encrypted INT;
    DECLARE terminalId INT;
    DECLARE tidUpdateId INT;
    SET @profileId = (SELECT profile_id from approvals a WHERE a.approval_id = approval_id);
    SET @elementId = (SELECT data_element_id from approvals a WHERE a.approval_id = approval_id);
    SET @approval_type = (SELECT change_type FROM approvals a WHERE a.approval_id = approval_id);
    SET @is_encrypted = (SELECT a.is_encrypted FROM approvals a WHERE a.approval_id = approval_id);

    IF @approval_type = 1 OR  @approval_type = 2 THEN
        SET @newVal = (SELECT new_value from approvals a WHERE a.approval_id = approval_id);
        SET @oldVal = (SELECT current_value from approvals a WHERE a.approval_id = approval_id);
        SET @terminalId=(select name from profile where profile_id=@profileId and profile_type_id=5);
        SET @tidUpdateId=(select tid_update_id from tid_updates where tid_id= @terminalId order by update_date desc limit 1);
        IF @oldVal='Third Party Application' AND @newVal!='Third Party Application' THEN
            UPDATE tid_updates SET third_party_apk="[]", update_date=NOW()
            WHERE tid_id=@terminalId and tid_update_id=@tidUpdateId;
        END IF;

        IF EXISTS (SELECT profile_data_id FROM profile_data pd WHERE pd.profile_id = @profileId AND pd.data_element_id = @elementId) THEN
            UPDATE profile_data pd SET pd.datavalue = @newVal, pd.updated_at=NOW(), pd.updated_by = approval_user, pd.is_encrypted = @is_encrypted
            WHERE pd.profile_id = @profileId AND pd.data_element_id = @elementId;
        ELSE
            SET @profType = (select p1.name from profile_type p1 inner join profile p2 on p1.profile_type_id = p2.profile_type_id where profile_id = @profileId);
            SET @dataloc = (select COUNT(d2.location_name) from data_element_locations_data_element d1 inner join data_element_locations d2 on d1.location_id = d2.location_id where data_element_id = @elementId and d2.location_name = 'acquirer_configuration');
            SET @globaloverrideval = (select not_overridable from profile_data where profile_id =1 and data_element_id =@elementId);
            IF (@profType = 'acquirer' and @dataloc = 1 and @globaloverrideval=1) THEN
                insert into profile_data( profile_id, data_element_id, datavalue, version,  updated_at,updated_by, created_at, created_by, approved, overriden, is_encrypted,not_overridable)
                values (@profileId,
                        @elementId,
                        @newVal,
                        1,
                        NOW(),
                        approval_user,
                        NOW(),
                        approval_user,
                        1,
                        0,
                        @is_encrypted,
                        1);
            ELSE
                insert into profile_data( profile_id, data_element_id, datavalue, version,  updated_at,updated_by, created_at, created_by, approved, overriden, is_encrypted)
                values (@profileId,
                        @elementId,
                        @newVal,
                        1,
                        current_timestamp,
                        approval_user,
                        current_timestamp,
                        approval_user,
                        1,
                        CASE @approval_type WHEN 2 THEN 1 ELSE 0 END,
                        @is_encrypted); -- override if change type is overriden
            END IF;
        END IF;
        IF @terminalId IS NOT NULL THEN
            UPDATE approvals a SET a.tid_id=@terminalId WHERE a.approval_id = approval_id;
        END IF;
    ELSEIF @approval_type=4 THEN
        DELETE pd FROM profile_data pd WHERE pd.profile_id = @profileId AND pd.data_element_id = @elementId;
        UPDATE `site` s
            LEFT JOIN site_profiles sp ON sp.site_id = s.site_id
        SET s.updated_at = NOW()
        WHERE sp.profile_id = @profileId;
    ELSEIF @approval_type=6 THEN
        UPDATE tid
        SET serial = (SELECT a.new_value from approvals a WHERE a.approval_id = approval_id)
        WHERE tid_id = (SELECT a.tid_id from approvals a WHERE a.approval_id = approval_id);
    ELSEIF @approval_type=11 THEN
        SET @paymentServiceGroupName = (SELECT name from payment_service_group where group_id = (SELECT group_id from payment_service where service_id = (SELECT current_value from approvals a WHERE a.approval_id = approval_id)));
        SET @paymentServiceName = (SELECT name from payment_service where service_id = (SELECT current_value from approvals a WHERE a.approval_id = approval_id));
        UPDATE approvals a SET approved = 1, approved_by = approval_user, approved_at = NOW(),a.current_value = @paymentServiceName,a.new_value = @paymentServiceGroupName  WHERE a.approval_id = approval_id;
    END IF;
    UPDATE approvals a SET approved = 1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approval_id;
end