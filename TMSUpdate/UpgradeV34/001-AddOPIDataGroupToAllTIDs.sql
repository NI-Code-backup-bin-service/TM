--multiline;
insert into profile_data_group(
      profile_id,
      data_group_id,
      version,
      updated_at,
      updated_by,
      created_at,
      created_by
    )

    select
		profile_id as profile_id,
        (select data_group_id from data_group where name = "opi") as data_group_id,
        1 as version,
        current_timestamp as updated_at,
        "system" as updated_by,
        current_timestamp as created_at,
        "system" as created_by

    from profile

	where profile_type_id = (select profile_type_id from profile_type where name = "tid")