--multiline;
CREATE PROCEDURE `site_data_fetch`(
 IN site_id int,
 IN updatesFrom BIGINT
)
BEGIN

DECLARE updatesFromDate DATETIME ;
SET @updatesFromDate = FROM_UNIXTIME(updatesFrom);

SET SESSION group_concat_max_len = 1000000;

set @site_id = site_id;

select concat("coalesce(" , group_concat("v",l.priority,".datavalue"), ")")
from (select priority from profile_type order by priority) l
into @datavalues;

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
	      "select
		e.site_id,
		e.data_group_id,
		dg.name as 'data_group',
		e.data_element_id,
		e.name,",
		@levels,
		" as `source`,",
		@datavalues,
		"as `datavalue`,",
        @overrides,
        "as `overriden`,
		e.datatype,
		e.is_allow_empty,
        e.max_length,
        e.validation_expression,
        e.validation_message,
        e.front_end_validate,
        e.options"
	      " from site_data_elements e",
	      "join data_group dg on dg.data_group_id = e.data_group_id",
	      @joins,
	      " where e.site_id = ",
	      @site_id,
          "AND e.updated_at >= ",
          "'",
		  @updatesFromDate,
	      "' ") into @sql;

PREPARE stmt1 FROM @sql;
EXECUTE stmt1;
DEALLOCATE PREPARE stmt1;
END