--multiline
CREATE VIEW chain_data AS
select profile_ids.chain_profile_id                                                                  AS profile_id,
       de.data_element_id                                                                            AS data_element_id,
       coalesce(apd_chain.source, apd_acquirer.source, apd_global.source)                            AS source,
       coalesce(apd_chain.datavalue, apd_acquirer.datavalue, apd_global.datavalue)                   AS datavalue,
       coalesce(apd_chain.overriden, apd_acquirer.overriden, apd_global.overriden)                   AS overriden,
       coalesce(apd_chain.is_encrypted, apd_acquirer.is_encrypted, apd_global.is_encrypted)          AS is_encrypted,
       coalesce(apd_chain.not_overridable, apd_acquirer.not_overridable, apd_global.not_overridable) AS not_overridable
from (((((((select cp.chain_profile_id AS chain_profile_id,
                   cp.acquirer_id      AS acquirer_profile_id,
                   1                   AS global_profile_id
            from NextGen_TMS.chain_profiles cp)) profile_ids join NextGen_TMS.data_element de) join NextGen_TMS.data_group dg on ((dg.data_group_id = de.data_group_id))) left join NextGen_TMS.approved_profile_data apd_chain on ((
        (apd_chain.profile_id = profile_ids.chain_profile_id) and
        (apd_chain.data_element_id = de.data_element_id)))) left join NextGen_TMS.approved_profile_data apd_acquirer on ((
        (apd_acquirer.profile_id = profile_ids.acquirer_profile_id) and
        (apd_acquirer.data_element_id = de.data_element_id))))
         left join NextGen_TMS.approved_profile_data apd_global
                   on (((apd_global.profile_id = profile_ids.global_profile_id) and
                        (apd_global.data_element_id = de.data_element_id))));