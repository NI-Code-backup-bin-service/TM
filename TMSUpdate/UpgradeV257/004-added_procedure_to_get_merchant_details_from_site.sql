--multiline
CREATE PROCEDURE `get_merchant_value_by_key`(IN tid int, IN data_element_name varchar(50))
BEGIN
SELECT distinct sd.datavalue FROM site_data_elements as e
    INNER JOIN data_element de ON de.data_element_id = e.data_element_id
    INNER JOIN data_group dg ON dg.data_group_id = e.data_group_id
    LEFT JOIN site_data sd ON sd.site_id = e.site_id AND sd.data_element_id = e.data_element_id
    WHERE e.site_id = (SELECT site_id from tid_site where tid_id = tid)
    AND e.data_group_id IN (SELECT data_group_id FROM data_group WHERE name='pullpayments')
    AND e.name=data_element_name AND e.location_name IN ('site_configuration','tid_override')
    order by e.sort_order_in_group;
END