--multiline;
CREATE PROCEDURE get_data_groups_from_tid(IN tid int)
BEGIN
	SELECT distinct(dg.name)
	FROM profile_data pd
    JOIN site_profiles sp ON pd.profile_id=sp.profile_Id
    JOIN tid_site ts ON sp.site_id = ts.site_id
    JOIN profile p ON p.profile_id = sp.profile_id
    JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
    JOIN data_element de ON pd.data_element_id=de.data_element_id
    JOIN data_group dg ON dg.data_group_id=de.data_group_id
    WHERE ts.tid_id=tid
      AND pt.name = "site";
END