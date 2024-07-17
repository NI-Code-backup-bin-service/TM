--multiline
CREATE PROCEDURE `get_report_data_using_profile`(IN ProfileId int)
BEGIN
select dg.data_group_id,dg.name,de.data_element_id,de.name,pd.datavalue,de.datatype,de.is_allow_empty,de.max_length,de.validation_expression,de.validation_message,de.front_end_validate,de.options,de.displayname_en from profile_data pd
join data_element de on pd.data_element_id=de.data_element_id
join data_group dg on dg.data_group_id=de.data_group_id
where profile_id=ProfileId;
END