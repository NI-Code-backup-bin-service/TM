--multiline
CREATE PROCEDURE GET_MID_BY_MID(IN mid varchar(20))
BEGIN
    /*
     This procedure returns matching primary or secondary MIDs for a given MID and indicates the primary mid if the
     returned mid is a secondaryMid
     */
    SELECT pdSecondary.datavalue as 'mid', 'secondaryMid' as 'midType', pdPrimary.datavalue 'primaryMid'
    FROM profile_data pdSecondary
             INNER JOIN data_element deSecondary ON pdSecondary.data_element_id = deSecondary.data_element_id
             INNER JOIN profile_data pdPrimary ON pdSecondary.profile_id = pdPrimary.profile_id
             INNER JOIN data_element dePrimary ON pdPrimary.data_element_id = dePrimary.data_element_id
    WHERE deSecondary.name = 'secondaryMid' AND dePrimary.name = 'MerchantNo' AND (pdPrimary.datavalue = mid OR pdSecondary.datavalue = mid)
    UNION ALL
    SELECT pd.datavalue as 'mid', 'primaryMid' as 'midType', 0
    FROM profile_data pd
             INNER JOIN data_element de ON pd.data_element_id = de.data_element_id
    WHERE de.name = 'MerchantNo' AND pd.datavalue = mid;
END;