--multiline;
CREATE PROCEDURE add_data_groups_to_tid_profiles()
INSERT ignore INTO profile_data_group (
  profile_id, data_group_id, version,
  updated_at, updated_by, created_at,
  created_by
)
SELECT
    all_groups.tid_profile_id,
    all_groups.data_group_id,
    1,
    CURRENT_TIMESTAMP,
    NULL,
    CURRENT_TIMESTAMP,
    NULL
FROM
    (
        SELECT
            dg.data_group_id,
            ts.tid_profile_id
        from
            data_group dg cross
                              join tid_site ts
        where
            ts.tid_profile_id is not null
    ) as all_groups
        left join (
        select
            pdg.data_group_id,
            p.profile_id
        from
            profile_data_group pdg
                inner join profile p on p.profile_id = pdg.profile_id
                inner join profile_type pt on pt.profile_type_id = p.profile_type_id
        where
                pt.name = 'tid'
    ) as already_allocated on all_groups.data_group_id = already_allocated.data_group_id
        and all_groups.tid_profile_id = already_allocated.profile_id
where
    already_allocated.data_group_id is null;