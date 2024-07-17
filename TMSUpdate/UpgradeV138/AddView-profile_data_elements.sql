--multiline
CREATE VIEW profile_data_elements AS
SELECT
    dg.name data_group_name,
    de.name data_element_name,
    pd.profile_id,
    pt.priority profile_type_priority,
    pd.datavalue value
FROM profile_data pd
         INNER JOIN data_element de on pd.data_element_id = de.data_element_id
         INNER JOIN data_group dg on de.data_group_id = dg.data_group_id
         INNER JOIN profile p on pd.profile_id = p.profile_id
         INNER JOIN profile_type pt on p.profile_type_id = pt.profile_type_id;

