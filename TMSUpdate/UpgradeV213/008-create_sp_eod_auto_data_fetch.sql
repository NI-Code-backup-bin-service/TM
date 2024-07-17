--multiline;
CREATE PROCEDURE `eod_auto_data_fetch`(IN siteId int)
BEGIN
    WITH profileData AS (
        SELECT pd.datavalue, de.name, pd.data_element_id,
               ROW_NUMBER() OVER(PARTITION BY pd.data_element_id ORDER BY pt.priority) AS RowNum
        FROM profile_data pd
        LEFT JOIN data_element de
            ON de.data_element_id=pd.data_element_id
        JOIN data_group dg
            ON dg.data_group_id=de.data_group_id
        JOIN site_profiles sp
            ON pd.profile_id = sp.profile_id
        JOIN profile p
            ON p.profile_id = sp.profile_id
        JOIN profile_type pt
            ON p.profile_type_id = pt.profile_type_id
        WHERE  dg.name="endOfDay"
        AND de.name IN ("auto", "time")
        AND sp.site_id = siteId
    )
    SELECT *
    FROM profileData
    WHERE RowNum = 1;
END