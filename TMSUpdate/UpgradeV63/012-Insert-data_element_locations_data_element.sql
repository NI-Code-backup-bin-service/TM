--multiline
INSERT IGNORE INTO data_element_locations_data_element
(
    location_id,
    data_element_id
)
SELECT data_elements.* FROM (
                                SELECT (SELECT location_id
                                        FROM data_element_locations
                                        WHERE location_name = 'site_configuration' and profile_type_id = 4),
                                       de1.data_element_id
                                from data_element de1
                                where de1.name != 'dailyTxnCleanseTime'
                                UNION
                                SELECT (SELECT location_id
                                        FROM data_element_locations
                                        WHERE location_name = 'fraud' and profile_type_id = 4),
                                       data_element_id
                                from data_element de2
                                where de2.name = 'dailyTxnCleanseTime'
                                UNION
                                SELECT (SELECT location_id
                                        FROM data_element_locations
                                        WHERE location_name = 'fraud' and profile_type_id = 5),
                                       data_element_id
                                from data_element de3
                                where de3.name = 'dailyTxnCleanseTime'
                            ) data_elements