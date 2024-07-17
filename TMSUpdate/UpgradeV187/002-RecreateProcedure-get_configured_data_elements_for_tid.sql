--multiline
CREATE PROCEDURE get_configured_data_elements_for_tid(IN p_tid INT, IN p_lastchecked BIGINT)
begin
  IF ( (SELECT updated_at
        FROM   tid_site
        WHERE  tid_id = p_tid) > From_unixtime(p_lastchecked) )
      OR ( (SELECT s.updated_at
            FROM   tid_site ts
                   INNER JOIN site s
                           ON ts.site_id = s.site_id
            WHERE  ts.tid_id = p_tid) > From_unixtime(p_lastchecked) ) THEN
    SET @v_updatesfromdate = from_unixtime(0);
ELSE
    SET @v_updatesfromdate = from_unixtime(p_lastchecked);
end IF;

SELECT dg.name data_group_name,
       de.name data_element_name,
       de.datatype,
       pt.priority,
       pd.datavalue,
       IF(pd.is_encrypted = 0
              OR pd.is_encrypted = ''
              OR pd.is_encrypted IS NULL, 0, 1)
FROM   (SELECT data.profile_id AS profile_id,
               data.tid_id     AS tid_id,
               data.site_id    AS site_id
        FROM   (SELECT ts.tid_profile_id AS profile_id,
                       ts.tid_id         AS tid_id,
                       ts.site_id        AS site_id
                FROM   (tid_site ts
                    LEFT JOIN site_profiles sp
                        ON (( ts.site_id = sp.site_id )))
                WHERE  ( ts.tid_id = p_tid )
                UNION
                SELECT sp.profile_id AS profile_id,
                       ts.tid_id     AS tid_id,
                       sp.site_id    AS site_id
                FROM   (tid_site ts
                    LEFT JOIN site_profiles sp
                        ON (( ts.site_id = sp.site_id )))
                WHERE  ( ts.tid_id = p_tid )) data) tsp
           INNER JOIN profile p
                      ON tsp.profile_id = p.profile_id
           INNER JOIN profile_type pt
                      ON p.profile_type_id = pt.profile_type_id
           LEFT JOIN profile_data pd
                     ON p.profile_id = pd.profile_id
           INNER JOIN data_element de
                      ON pd.data_element_id = de.data_element_id
           INNER JOIN data_group dg
                      ON de.data_group_id = dg.data_group_id
WHERE  tsp.tid_id = p_tid
  AND ( pd.updated_at IS NULL
    OR pd.updated_at >= @v_updatesfromdate )
ORDER  BY dg.data_group_id;
end