--multiline
CREATE VIEW approved_profile_data AS
select
    p.profile_id,
    pd.data_element_id,
    pt.name as source,
    pd.datavalue,
    pd.overriden
from profile_data pd
join profile p
on p.profile_id = pd.profile_id
join profile_type pt
on p.profile_type_id = pt.profile_type_id
where pd.approved = 1