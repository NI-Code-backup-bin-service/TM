--multiline
CREATE PROCEDURE `tid_list_fetch`(IN search_term varchar(255), IN acquirers text)
BEGIN
    set @search = upper(concat('%', ifnull(search_term, ''), '%'));
    SELECT JSON_ARRAYAGG(PED_DATA.PED_DATA)
    FROM (
             SELECT JSON_MERGE_PRESERVE(
                            JSON_OBJECT(
                                    'tid', t.tid_id,
                                    'serial', t.serial,
                                    'pin', t.PIN,
                                    'expiryTime', t.ExpiryDate,
                                    'activationTime', t.ActivationDate,
                                    'presence', t.Presence,
                                    'siteId', s.site_id,
                                    'siteName', pdn.datavalue,
                                    'chainId', p3.profile_id,
									'chainName', p3.name,
                                    'merchantId', pd.datavalue,
                                    'appVer', t.software_version,
                                    'firmwareVer', t.firmware_version,
                                    'lastTransaction', date_format(t.last_transaction_time, "%Y-%m-%dT%H:%i:%sZ"),
                                    'lastCheckedTime',
                                    date_format(from_unixtime(t.last_checked_time / 1000), "%Y-%m-%dT%H:%i:%sZ"),
                                    'confirmedTime',
                                    date_format(from_unixtime(t.confirmed_time / 1000), "%Y-%m-%dT%H:%i:%sZ"),
                                    'lastAPKDownload', date_format(t.last_apk_download, "%Y-%m-%dT%H:%i:%sZ")
                                ),
                            JSON_MERGE_PRESERVE(
                                    JSON_REMOVE(
                                            JSON_OBJECTAGG(
                                                    IF(tid_override.overrideFieldType IN ('BOOLEAN', 'INTEGER', 'LONG'),
                                                       tid_override.overrideFieldName, 'null'),
                                                    CASE
                                                        WHEN tid_override.overrideFieldType = 'BOOLEAN' THEN
                                                            if(UPPER(tid_override.overrideFieldValue) = 'TRUE',
                                                               CAST(true AS JSON), CAST(false AS JSON))
                                                        WHEN tid_override.overrideFieldType IN ('INTEGER', 'LONG') THEN
                                                            CAST(CAST(tid_override.overrideFieldValue AS UNSIGNED) AS JSON)
                                                        END), '$.null'),
                                    JSON_REMOVE(
                                            JSON_OBJECTAGG(
                                                    IF(tid_override.overrideFieldType IN ('STRING', 'JSON'),
                                                       tid_override.overrideFieldName, 'null'),
                                                    IF(tid_override.overrideFieldType IN ('STRING', 'JSON'),
                                                       tid_override.overrideFieldValue, null)
                                                ), '$.null')
                                )
                        ) PED_DATA
             from tid t
                      left join tid_site ts on ts.tid_id = t.tid_id
                      left join site s on s.site_id = ts.site_id
                      LEFT JOIN (site_profiles tp2
                 join profile p2 on p2.profile_id = tp2.profile_id
                 join profile_type pt2 on pt2.profile_type_id = p2.profile_type_id and pt2.priority = 2)
                                on tp2.site_id = s.site_id
                      LEFT JOIN profile_data pd ON pd.profile_id = p2.profile_id AND pd.data_element_id = 1
                 AND pd.version = (SELECT MAX(d.version)
                                   FROM profile_data d
                                   WHERE d.data_element_id = 1
                                     AND d.profile_id = p2.profile_id
                                     AND d.approved = 1)
                      LEFT JOIN profile_data pdn ON pdn.profile_id = p2.profile_id AND pdn.data_element_id = 3
                 AND pdn.version = (SELECT MAX(d.version)
                                    FROM profile_data d
                                    WHERE d.data_element_id = 3
                                      AND d.profile_id = p2.profile_id
                                      AND d.approved = 1)
				LEFT JOIN
				(site_profiles tp3
					join profile p3 on p3.profile_id = tp3.profile_id
					join profile_type pt3 on pt3.profile_type_id = p3.profile_type_id and pt3.priority = 3)
				on tp3.site_id = s.site_id

                      LEFT JOIN (site_profiles tp4
                 join profile p4 on p4.profile_id = tp4.profile_id
                 join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4)
                                on tp4.site_id = s.site_id
                 #Generates the TID overrides json object
                      LEFT JOIN (
                 select p.profile_id,
                        concat(dg.displayname_en, '.', de.name) As 'overrideFieldName',
                        pd.datavalue                            As 'overrideFieldValue',
                        de.datatype                             As 'overrideFieldType'
                 from profile p
                          LEFT JOIN profile_data pd ON
                     pd.profile_id = p.profile_id
                          INNER JOIN data_element de on
                     pd.data_element_id = de.data_element_id
                          INNER JOIN data_group dg on de.data_group_id = dg.data_group_id
             ) tid_override ON tid_override.profile_id = ts.tid_profile_id
             where ts.tid_id = t.tid_id
               and FIND_IN_SET(p4.name, acquirers)
               and (upper(t.tid_id) like @search
                 or upper(t.serial) like @search
                 or upper(pdn.datavalue) like @search)
             group by t.tid_id,
                      t.serial,
                      t.PIN,
                      t.ExpiryDate,
                      t.ActivationDate,
                      t.Presence,
                      s.site_id,
                      pdn.datavalue,
                      pd.datavalue,
                      p3.profile_id,
                      p3.name
         ) PED_DATA;
END
