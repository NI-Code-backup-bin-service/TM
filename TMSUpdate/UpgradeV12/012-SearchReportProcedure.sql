--multiline;
CREATE PROCEDURE `search_report_data`( IN id int, IN profileType text, IN updatesFrom BIGINT )
BEGIN
  DECLARE updatesFromDate DATETIME ;
  SET @updatesFromDate = FROM_UNIXTIME(updatesFrom);
  SET SESSION group_concat_max_len = 1000000;
  set @id = id;
  select concat("coalesce(" , group_concat("v",l.priority,".datavalue"), ")")
  from (select priority from profile_type where name = profileType order by priority) l
       into @datavalues;
  select concat("coalesce(" , group_concat("v",l.priority,".updated_at"), ")")
  from (select priority from profile_type where name = profileType order by priority) l
       into @updatedAt;
  select concat("coalesce(" , group_concat("v",l.priority,".level"), ")")
  from (select priority from profile_type where name = profileType order by priority) l
       into @levels;
  select concat("coalesce(" , group_concat("v",l.priority,".overriden"), ")")
  from (select priority from profile_type where name = profileType order by priority) l
       into @overrides;
  select group_concat(concat("left join site_data v", l.priority, " on v",l.priority,".site_id = e.site_id and v",l.priority,".data_element_id = e.data_element_id and v", l.priority, ".priority = ",l.priority
                        ," and v", l.priority,".version = (select MAX(d",l.priority,".version) from site_data d",l.priority," where d",l.priority,".site_id = e.site_id and d",l.priority,".data_element_id = e.data_element_id and d", l.priority, ".priority = ",l.priority,")") SEPARATOR ' ')
  from (select priority from profile_type where name = profileType order by priority) l
       into @joins;
  select concat_ws(" ",
                   "select
                   e.site_id,
                   e.data_group_id,
                   dg.name as 'data_group',
                   e.data_element_id,
                   e.name,",
                   @levels, " as `source`,",
                   @datavalues, "as `datavalue`,",
                   @overrides, "as `overriden`,
                   e.datatype,
                   e.is_allow_empty,
                   e.max_length,
                   e.validation_expression,
                   e.validation_message,
                   e.front_end_validate,
                   e.options,
                   e.display_name",
                   " from site_data_elements e",
                   "join data_group dg on dg.data_group_id = e.data_group_id",
                   @joins,
                   " where e.site_id = ", @id,
                   "AND ",
                   @updatedAt, " >= '", @updatesFromDate,
                   "' ") into @sql;
  PREPARE stmt1 FROM @sql;
  EXECUTE stmt1;
  DEALLOCATE PREPARE stmt1;
END