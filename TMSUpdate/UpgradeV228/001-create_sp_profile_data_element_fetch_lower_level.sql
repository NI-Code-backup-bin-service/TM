--multiline;
CREATE PROCEDURE `profile_data_element_fetch_lower_level`(IN siteId int,IN dataElementName varchar(255), IN dataGroupName varchar(255))
BEGIN


WITH profileData AS (
    SELECT pd.datavalue, pd.data_element_id,pt.priority,
           ROW_NUMBER() OVER(PARTITION BY pd.data_element_id ORDER BY pt.priority DESC ) AS RowNum
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
    WHERE  dg.name= dataGroupName
      AND de.name = dataElementName
      AND sp.site_id = siteId
)

SELECT datavalue,priority FROM profileData;
END;