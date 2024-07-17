--multiline
create view approved_profile_data as
select p.profile_id       AS profile_id,
       pd.data_element_id AS data_element_id,
       pt.name            AS source,
       pd.datavalue       AS datavalue,
       pd.overriden       AS overriden,
       pd.is_encrypted    AS is_encrypted,
       pd.not_overridable AS not_overridable
from ((NextGen_TMS.profile_data pd join NextGen_TMS.profile p on ((p.profile_id = pd.profile_id)))
         join NextGen_TMS.profile_type pt on ((p.profile_type_id = pt.profile_type_id)))
where (pd.approved = 1);

