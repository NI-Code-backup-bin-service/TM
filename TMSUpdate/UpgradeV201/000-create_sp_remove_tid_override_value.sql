--multiline;
create procedure remove_tid_override_value(IN TID TEXT, IN DataGroupName TEXT, IN ElementName TEXT)
BEGIN
    delete
        pd
    from
        tid t
        inner join tid_site ts on
            ts.tid_id = t.tid_id
        INNER join profile p on
            p.profile_id = ts.tid_profile_id
        inner JOIN profile_data pd on
            pd.profile_id = p.profile_id
        inner JOIN data_element de on
            de.data_element_id = pd.data_element_id
        inner join data_group dg on
            dg.data_group_id = de.data_group_id
    WHERE
        t.tid_id = TID and dg.name = DataGroupName and de.name = ElementName;
END