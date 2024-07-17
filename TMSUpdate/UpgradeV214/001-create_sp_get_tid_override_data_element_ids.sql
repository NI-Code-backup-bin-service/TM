--multiline
CREATE PROCEDURE get_tid_override_data_element_ids()
BEGIN
    SELECT delde.data_element_id FROM data_element_locations_data_element delde
    JOIN data_element de
        ON de.data_element_id=delde.data_element_id
    WHERE location_id IN
        (SELECT location_id FROM data_element_locations WHERE location_name IN('site_configuration','tid_override'))
    AND de.tid_overridable=1
    GROUP BY delde.data_element_id
    HAVING COUNT(delde.data_element_id)=1 AND SUM(delde.location_id)=
        (SELECT location_id FROM data_element_locations WHERE location_name = 'tid_override');
END