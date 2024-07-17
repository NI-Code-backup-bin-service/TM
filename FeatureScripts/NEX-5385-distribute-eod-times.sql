update data_group dg inner join data_element de on dg.data_group_id = de.data_group_id inner join profile_data pd on de.data_element_id = pd.data_element_id set pd.datavalue = CONCAT('00:', LPAD((pd.profile_data_id % 12) * 5, 2, '0')) where dg.name = 'endOfDay' and de.name = 'time' and pd.datavalue = '00:00';