--multiline;
CREATE PROCEDURE get_profile_change_history(IN profileId int)
BEGIN
    DECLARE siteCount INT;
SELECT count(*) into siteCount from profile p inner join profile_type pt on p.profile_type_id=pt.profile_type_id where pt.name ='site' and profile_id = profileId;
IF siteCount > 0 THEN
SELECT
    CONCAT(
            dg.name, '/',
        -- NEX-7085 - in some versions of the db, displayname_en can be empty instead of just null
        -- this IF now checks for both null and empty, fixing the display problem occuring previously
            IF(de.displayname_en IS NULL OR de.displayname_en = '', de.name, de.displayname_en)
        ),
    a.change_type as change_type,
    a.current_value as original_value,
    a.new_value,
    a.created_by as updated_by,
    a.created_at as updated_at,
    a.approved,
    a.tid_id,
    a.is_password,
    a.is_encrypted
FROM approvals a
         LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
         LEFT JOIN data_group dg ON dg.data_group_id = de.data_group_id
WHERE a.profile_id IN (
    SELECT profileId
    UNION
    select ts.tid_profile_id from site s inner join tid_site ts on ts.site_id=s.site_id inner join site_profiles sp on sp.site_id=s.site_id where sp.profile_id=profileId
)
ORDER BY a.created_at DESC;
ELSE
SELECT
    CONCAT(
            dg.name, '/',
        -- NEX-7085 - in some versions of the db, displayname_en can be empty instead of just null
        -- this IF now checks for both null and empty, fixing the display problem occuring previously
            IF(de.displayname_en IS NULL OR de.displayname_en = '', de.name, de.displayname_en)
        ),
    a.change_type as change_type,
    a.current_value as original_value,
    a.new_value,
    a.created_by as updated_by,
    a.created_at as updated_at,
    a.approved,
    a.tid_id,
    a.is_password,
    a.is_encrypted
FROM approvals a
         LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
         LEFT JOIN data_group dg ON dg.data_group_id = de.data_group_id
WHERE a.profile_id = profileId
ORDER BY a.created_at DESC;
END IF;
END;