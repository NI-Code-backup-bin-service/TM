--multiline;
CREATE PROCEDURE `create_tid_override`(IN siteId INT, IN profileId INT, IN updatedBy varchar(255))
BEGIN
    INSERT INTO profile_data
    (profile_id, data_element_id, datavalue, version, updated_by, created_by, approved, overriden, is_encrypted)
    (WITH profileData AS (
        SELECT pd.data_element_id, pd.datavalue, pd.is_encrypted,
               ROW_NUMBER() OVER(PARTITION BY pd.data_element_id ORDER BY pt.priority) AS RowNum
        FROM profile_data pd
        JOIN data_element de
            ON de.data_element_id=pd.data_element_id
        JOIN site_profiles sp
            ON pd.profile_id = sp.profile_id
        JOIN profile p
            ON p.profile_id = sp.profile_id
        JOIN profile_type pt
            ON p.profile_type_id = pt.profile_type_id
        WHERE sp.site_id = siteId
        and de.tid_overridable=1
        and de.data_group_id IN (SELECT DISTINCT(de.data_group_id)
            FROM profile_data pd
                JOIN site_profiles sp
                    ON pd.profile_id= sp.profile_id
                JOIN profile p
                    ON p.profile_id=sp.profile_id
                JOIN profile_type pt
                    ON pt.profile_type_id=p.profile_type_id
                JOIN data_element de
                    ON de.data_element_id=pd.data_element_id
            WHERE sp.site_id=siteId )
    )
    SELECT profileId, data_element_id, datavalue, 0, updatedBy, updatedBy, 1, 0, is_encrypted
    FROM profileData
    WHERE RowNum = 1);
END