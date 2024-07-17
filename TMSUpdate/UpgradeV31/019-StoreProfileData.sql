--multiline
CREATE PROCEDURE `store_profile_data`(
    IN profile_ident int,
    IN data_element int,
    IN datavalue varchar(255),
    IN updated_by varchar(255),
    IN approved int,
    IN overriden int,
    IN is_encrypted BOOL
)
BEGIN

    declare currentVersion int;

    set currentVersion = (select MAX(version) from profile_data
                          where `data_element_id`= data_element
                            and `profile_id` = profile_ident);

    if currentVersion IS NULL THEN
        set currentVersion = 0;
    END IF;

    insert into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by,created_at, created_by, approved, overriden, is_encrypted)
    values (profile_ident, data_element, datavalue, currentVersion+1, current_timestamp, updated_by, current_timestamp, updated_by, approved, overriden, is_encrypted);
END