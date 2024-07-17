--multiline;
CREATE PROCEDURE `profile_list_fetch_by_type_name`(
    IN profile_type_name varchar(50),
    IN acquirers TEXT
)
BEGIN
    if profile_type_name = "acquirer" then
        select
            p.profile_id,
            p.name as 'type_name'
        from profile_type t
            left join profile p on p.profile_type_id = t.profile_type_id
        where FIND_IN_SET(p.name, acquirers)
            and t.name = profile_type_name;
    else
        select distinct
            p.profile_id,
            p.name as 'type_name'
        from profile_type t
        left join profile p on p.profile_type_id = t.profile_type_id
        left join chain_profiles cp on cp.chain_profile_id = p.profile_id
        left join profile p2 on p2.profile_id = cp.acquirer_id
        where FIND_IN_SET(p2.name, acquirers);
    end if;
END