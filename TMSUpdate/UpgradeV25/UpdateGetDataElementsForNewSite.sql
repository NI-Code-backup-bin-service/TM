--multiline;
CREATE PROCEDURE `get_data_elements_for_new_site`(IN acquirer_id INT, IN chain_id INT)
BEGIN
  select
    pdg.data_group_id,
    dg.name as 'data_group',
    e.data_element_id,
    e.name,
    pt.name as `source`,
    pd.datavalue,
    e.datatype,
    e.is_allow_empty,
    e.options,
    e.displayname_en
  from `profile` p
         left join profile_data_group pdg ON pdg.profile_id = p.profile_id
         join data_group dg on dg.data_group_id = pdg.data_group_id
         left join data_element e ON e.data_group_id = dg.data_group_id
         left join profile_data pd ON pd.data_element_id = e.data_element_id AND pd.profile_id = p.profile_id
         left join profile_type pt ON p.profile_type_id = pt.profile_type_id
  where p.profile_id = acquirer_id
  UNION
  select
    pdg.data_group_id,
    dg.name as 'data_group',
    e.data_element_id,
    e.name,
    pt.name as `source`,
    pd.datavalue,
    e.datatype,
    e.is_allow_empty,
    e.options,
    e.displayname_en
  from `profile` p
         left join profile_data_group pdg ON pdg.profile_id = p.profile_id
         join data_group dg on dg.data_group_id = pdg.data_group_id
         left join data_element e ON e.data_group_id = dg.data_group_id
         left join profile_data pd ON pd.data_element_id = e.data_element_id AND pd.profile_id = p.profile_id
         left join profile_type pt ON p.profile_type_id = pt.profile_type_id
  where p.profile_id =  chain_id ;
END