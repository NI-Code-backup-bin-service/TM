--multiline;

CREATE PROCEDURE `get_element_value`(
IN profile_id int, 
IN element_id int
)
BEGIN
	select datavalue 
    from profile_data pd 
    where pd.data_element_id = element_id and
    pd.profile_id = profile_id
    order by pd.version desc
    Limit 1;
END ;
