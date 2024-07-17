--multiline
CREATE PROCEDURE GET_TID_BY_TID(IN tid BIGINT)
BEGIN
    /*
     This procedure returns matching primary or secondary TIDs for a given TID and indicates the primary tid if the
     returned tid is a secondaryTid
     */
    SELECT datavalue as 'tid', 'secondaryTid' as 'tidType', ts.tid_id as 'primaryTid'
    FROM profile_data pd
             INNER JOIN data_element de ON
            pd.data_element_id = de.data_element_id
             INNER JOIN tid_site ts ON
            pd.profile_id = ts.tid_profile_id
    WHERE de.name = 'secondaryTid' AND pd.datavalue = tid
    UNION ALL
    SELECT tid_id, 'primaryTid', 0
    FROM tid t
    WHERE t.tid_id = tid;
END;