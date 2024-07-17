--multiline;
CREATE PROCEDURE `required_software_version_fetch`( IN site_id int, IN updatesFrom BIGINT )
BEGIN
  SET @updatesFromDate = FROM_UNIXTIME(updatesFrom);
  SET SESSION group_concat_max_len = 1000000;
  set @site_id = site_id;
  select concat("coalesce(" , group_concat("v",l.priority,".datavalue"), ")")
  from (select priority from profile_type order by priority) l
       into @datavalues;
  select concat("coalesce(" , group_concat("v",l.priority,".updated_at"), ")")
  from (select priority from profile_type order by priority) l
       into @updatedAt;
  select concat("coalesce(" , group_concat("v",l.priority,".level"), ")")
  from (select priority from profile_type order by priority) l
       into @levels;
  select concat("coalesce(" , group_concat("v",l.priority,".overriden"), ")")
  from (select priority from profile_type order by priority) l
       into @overrides;
  select group_concat(concat("left join site_data v", l.priority, " on v",l.priority,".site_id = e.site_id and v",l.priority,".data_element_id = e.data_element_id and v", l.priority, ".priority = ",l.priority
                        ," and v", l.priority,".version = (select MAX(d",l.priority,".version) from site_data d",l.priority," where d",l.priority,".site_id = e.site_id and d",l.priority,".data_element_id = e.data_element_id and d", l.priority, ".priority = ",l.priority,")") SEPARATOR ' ')
  from (select priority from profile_type order by priority) l
       into @joins;
  select concat_ws(" ",
                   "select ",
                   @datavalues, "as `datavalue`",
                   " from site_data_elements e",
                   "join data_group dg on dg.data_group_id = e.data_group_id",
                   "left join data_element de on de.name = \"RequiredSoftwareVersion\"",
                   @joins,
                   " where e.site_id = ", @site_id,
                   "AND ",
                   " e.data_element_id = de.data_element_id ",
                   "AND ",
                   "( ",  @updatedAt, " IS NULL", "or ", @updatedAt, " >= '", @updatesFromDate, "' ",
                   ") limit 1"
           ) into @sql;
  PREPARE stmt1 FROM @sql;
  EXECUTE stmt1;
  DEALLOCATE PREPARE stmt1;
END