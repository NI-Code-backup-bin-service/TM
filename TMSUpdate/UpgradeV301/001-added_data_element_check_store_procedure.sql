--multiline
CREATE PROCEDURE check_data_element_exists( IN tidId INT,  IN data_element_id INT)
BEGIN
    SELECT DISTINCT (de.data_element_id) from profile_data_group AS pdg
        Left Join data_element AS de ON de.data_group_id = pdg.data_group_id
        Left Join data_group as dg ON dg.data_group_id = pdg.data_group_id
        Left Join site_profiles as sp ON sp.profile_id = pdg.profile_id
        Left Join tid_site_profiles as tsp ON tsp.site_id = sp.site_id
        Left Join profile as p ON p.profile_id = sp.profile_id
        Left Join profile_type as pt ON pt.profile_type_id = p.profile_type_id and pt.profile_type_id != (select profile_type_id from profile_type where name = 'global')
    WHERE tsp.tid_id = tidId and de.tid_overridable = 1 and de.data_element_id = data_element_id;
END;