--multiline
CREATE PROCEDURE approve_change(IN approval_id int, IN approval_user varchar(256))
BEGIN
    DECLARE approval_type int;
    DECLARE newVal MEDIUMTEXT;
    DECLARE oldVal MEDIUMTEXT;
    DECLARE profileId INT;
    DECLARE elementId INT;
    DECLARE is_encrypted INT;
    DECLARE terminalId INT;
    DECLARE tidUpdateId INT;
    SET profileId = (SELECT profile_id from approvals a WHERE a.approval_id = approval_id);
    SET elementId = (SELECT data_element_id from approvals a WHERE a.approval_id = approval_id);
    SET approval_type = (SELECT change_type FROM approvals a WHERE a.approval_id = approval_id);
    SET is_encrypted = (SELECT a.is_encrypted FROM approvals a WHERE a.approval_id = approval_id);

    IF approval_type = 1 OR  approval_type = 2 THEN
        SET newVal = (SELECT new_value from approvals a WHERE a.approval_id = approval_id);
        SET oldVal = (SELECT current_value from approvals a WHERE a.approval_id = approval_id);
        SET terminalId=(select name from profile where profile_id=profileId and profile_type_id=5);
        SET tidUpdateId=(select tid_update_id from tid_updates where tid_id= terminalId order by update_date desc limit 1);
        IF oldVal='Third Party Application' AND newVal!='Third Party Application' THEN
            UPDATE tid_updates SET third_party_apk="[]", update_date=NOW()
            WHERE tid_id=terminalId and tid_update_id=tidUpdateId;
        END IF;

        IF EXISTS (SELECT profile_data_id FROM profile_data pd WHERE pd.profile_id = profileId AND pd.data_element_id = elementId) THEN
            UPDATE profile_data pd SET pd.datavalue = newVal, pd.updated_at=NOW(), pd.updated_by = approval_user, pd.is_encrypted = is_encrypted
            WHERE pd.profile_id = profileId AND pd.data_element_id = elementId;
        ELSE
            insert into profile_data( profile_id, data_element_id, datavalue, version,  updated_at,updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (profileId,
                    elementId,
                    newVal,
                    1,
                    current_timestamp,
                    approval_user,
                    current_timestamp,
                    approval_user,
                    1,
                    CASE approval_type WHEN 2 THEN 1 ELSE 0 END,
                    is_encrypted); -- override if change type is overriden
        END IF;
    ELSEIF approval_type=4 THEN
        DELETE pd FROM profile_data pd WHERE pd.profile_id = profileId AND pd.data_element_id = elementId;
        UPDATE `site` s
            LEFT JOIN site_profiles sp ON sp.site_id = s.site_id
        SET s.updated_at = NOW()
        WHERE sp.profile_id = profileId;
    ELSEIF approval_type=6 THEN
        UPDATE tid
        SET serial = (SELECT a.new_value from approvals a WHERE a.approval_id = approval_id)
        WHERE tid_id = (SELECT a.tid_id from approvals a WHERE a.approval_id = approval_id);
    END IF;
    UPDATE approvals a SET approved = 1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approval_id;
end;

