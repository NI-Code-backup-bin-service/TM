--multiline;
CREATE PROCEDURE `get_opi_tid_map`()
BEGIN
    SELECT
        ts.tid_id,
        ts.site_id,
        t.serial,
        spd1.datavalue AS "merchantNo",
        spd2.datavalue AS "propertyId",
        COALESCE(tpd.datavalue, spd3.datavalue) AS "workstationNumber"

    FROM
        tid_site AS ts

        LEFT JOIN
            tid AS t
            ON ts.tid_id = t.tid_id

        LEFT JOIN
            profile_data AS tpd
            ON tpd.profile_id = ts.tid_profile_id
            AND tpd.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "workstationNumber")

        LEFT JOIN
            site_profiles AS sp
            ON sp.site_id = ts.site_id

        JOIN
            profile_data AS spd1
            ON spd1.profile_id = sp.profile_id
            AND spd1.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "merchantNo")

        JOIN
            profile_data AS spd2
            ON spd2.profile_id = sp.profile_id
            AND spd2.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "propertyId")

        JOIN
            profile_data AS spd3
            ON spd3.profile_id = sp.profile_id
            AND spd3.data_element_id = (SELECT data_element_id FROM data_element WHERE name = "workstationNumber");
END