--multiline;
CREATE PROCEDURE get_element_id_by_group_element_name(IN group_name varchar(255), IN element_name varchar(255))
BEGIN
    set @group_id = (select data_group_id from data_group dg
                   where dg.name = group_name);
    SELECT data_element_id FROM data_element WHERE data_group_id = @group_id AND `name` = element_name;
END