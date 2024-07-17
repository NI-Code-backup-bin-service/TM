--multiline
delete pd from profile_data pd
left join data_element de on de.data_element_id = pd.data_element_id
left join data_group dg on dg.data_group_id = de.data_group_id
left join profile p on p.profile_id = pd.profile_id
where dg.name = 'multiDeviceEmv' and p.name = 'global';