--multiline
CREATE PROCEDURE GET_TID_BY_TID(IN tid BIGINT)
BEGIN
    /*
     This procedure returns matching primary or secondary TIDs for a given TID and indicates the primary tid if the
     returned tid is a secondaryTid
     */
    SELECT datavalue         AS 'tid',
           'secondaryTid'    AS 'tidType',
           ts.tid_id         AS 'primaryTid',
           ts.tid_profile_id AS 'profileId'
    FROM   profile_data pd
               INNER JOIN data_element de
                          ON pd.data_element_id = de.data_element_id
               INNER JOIN tid_site ts
                          ON pd.profile_id = ts.tid_profile_id
    WHERE  de.NAME = 'secondaryTid'
      AND pd.datavalue = tid
    UNION ALL
    SELECT t.tid_id,
           'primaryTid',
           0,
           IFNULL(ts.tid_profile_id, -1) AS 'profileId'
    FROM   tid t
               INNER JOIN tid_site ts
                          ON t.tid_id = ts.tid_id
    WHERE  t.tid_id = tid;
END;