--multiline;
CREATE PROCEDURE update_required_software_version(IN target_version varchar(45))
begin
    update profile_data set datavalue = target_version, updated_by = 'system', created_at = now() Where profile_data_id = 336;
end;