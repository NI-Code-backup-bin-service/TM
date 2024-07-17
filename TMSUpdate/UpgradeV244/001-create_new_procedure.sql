--multiline;
CREATE PROCEDURE get_data_group_info_for_TID_export(IN data_group_id INT)
BEGIN
SELECT
    dg.name,
    de.name,
    de.displayname_en,
    de.export_display_index
FROM data_group dg
JOIN data_element de ON de.data_group_id = dg.data_group_id and de.export_display_index IS NOT NULL
WHERE dg.data_group_id = data_group_id
ORDER BY de.data_element_id;
END;