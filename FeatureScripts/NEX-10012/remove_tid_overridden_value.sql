-- please replace <Data_Element_Name> with the field name that needs to be removed from the TID override panel
UPDATE data_group dg inner join data_element de on dg.data_group_id = de.data_group_id set de.tid_overridable = 0 where de.name in (<Data_Element_Name>);
-- please replace <TID>, <Data_Group_Name>, <Data_Element_Name> with TID, data group name and data element name respectively,
-- so the TID level overridden value will be removed for the given TID
call remove_tid_override_value(<TID>, <Data_Group_Name>, <Data_Element_Name>);