--multiline
create procedure get_chain_data(IN $profile_id int, IN $acquirer_id int)
BEGIN
    SELECT
        de.data_group_id,
        dg.name as data_group,
        de.data_element_id,
        de.name,
        source,
        datavalue,
        overriden,
        de.datatype,
        de.is_allow_empty,
        de.max_length,
        de.validation_expression,
        de.validation_message,
        de.front_end_validate,
        de.options,
        de.displayname_en,
        cd.is_encrypted,
        de.is_password,
        IFNULL(cd.not_overridable, 0)
    FROM chain_data cd
             JOIN data_element de ON cd.data_element_id = de.data_element_id
             JOIN data_group dg ON dg.data_group_id = de.data_group_id
             JOIN profile_data_group pdg ON pdg.data_group_id = dg.data_group_id AND (pdg.profile_id = $profile_id OR pdg.profile_id = $acquirer_id)
    WHERE cd.profile_id = $profile_id;
END;

