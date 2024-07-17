--multiline;
CREATE PROCEDURE `chain_list_fetch`(
    IN search_term varchar(255),
    IN acquirers TEXT
)
BEGIN
    set @search = upper(concat('%', ifnull(search_term,''), '%'));
    select distinct
        p.profile_id as 'chain_profile_id',
        p.name,
        p2.name
    from profile p
             left join chain_profiles cp on cp.chain_profile_id = p.profile_id
             left join profile p2 on p2.profile_id = cp.acquirer_id
    where upper(p.name) like @search
      and FIND_IN_SET(p2.name, acquirers);
END