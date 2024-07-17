--multiline
CREATE PROCEDURE GET_MID_BY_MID(IN mid varchar(20), primaryDataID INT, secondaryDataID INT)
BEGIN
    /*
     This procedure returns matching primary or secondary MIDs for a given MID and indicates the primary mid if the
     returned mid is a secondaryMid
     */
    SELECT pdSecondary.datavalue as 'mid', 'secondaryMid' as 'midType', pdPrimary.datavalue 'primaryMid'
    FROM profile_data pdSecondary
             INNER JOIN profile_data pdPrimary ON pdSecondary.profile_id = pdPrimary.profile_id
    WHERE pdSecondary.data_element_id = secondaryDataID AND pdPrimary.data_element_id = primaryDataID AND (pdPrimary.datavalue = mid OR pdSecondary.datavalue = mid)

    UNION ALL

    SELECT pd.datavalue as 'mid', 'primaryMid' as 'midType', 0
    FROM profile_data pd
    WHERE pd.datavalue = mid AND pd.data_element_id = primaryDataID;

END;