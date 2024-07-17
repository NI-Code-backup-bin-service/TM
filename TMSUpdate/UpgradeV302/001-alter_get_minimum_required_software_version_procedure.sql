--multiline;
CREATE PROCEDURE `get_minimum_required_software_version`(IN tidID INT)
BEGIN
    SELECT pd.datavalue
    FROM profile_data pd
    JOIN data_element de ON pd.data_element_id = de.data_element_id
    JOIN data_group dg ON de.data_group_id = dg.data_group_id
    JOIN site_profiles sp ON pd.profile_id = sp.profile_id
    JOIN profile p ON p.profile_id = sp.profile_id
    JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
    JOIN tid_site ts ON ts.site_id = sp.site_id
    WHERE dg.name = "core"
      AND de.name = "RequiredSoftwareVersion"
      AND ts.tid_id = tidID
    ORDER BY pt.priority
    LIMIT 1;
END;