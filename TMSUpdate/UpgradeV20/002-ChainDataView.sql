--multiline
CREATE VIEW chain_data AS
select
	profile_ids.chain_profile_id as profile_id,
    de.data_element_id,
    coalesce(apd_chain.source, apd_acquirer.source, apd_global.source) as source,
    coalesce(apd_chain.datavalue, apd_acquirer.datavalue, apd_global.datavalue) as datavalue,
    coalesce(apd_chain.overriden, apd_acquirer.overriden, apd_global.overriden) as overriden
from (
	select
		chain_profile_id,
		acquirer_id as acquirer_profile_id,
        1 as global_profile_id
	from chain_profiles cp) profile_ids
join data_element de
join data_group dg
on dg.data_group_id = de.data_group_id
left join approved_profile_data apd_chain
on apd_chain.profile_id = profile_ids.chain_profile_id
and apd_chain.data_element_id = de.data_element_id
left join approved_profile_data apd_acquirer
on apd_acquirer.profile_id = profile_ids.acquirer_profile_id
and apd_acquirer.data_element_id = de.data_element_id
left join approved_profile_data apd_global
on apd_global.profile_id = profile_ids.global_profile_id
and apd_global.data_element_id = de.data_element_id