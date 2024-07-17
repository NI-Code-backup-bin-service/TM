--multiline;
CREATE PROCEDURE `site_store`(
  IN site_id int,
  IN version int,
  IN updated_by varchar(255)
)
BEGIN
  IF (site_id = -1)
  THEN
    insert into site(
      version,
      updated_at,
      updated_by,
      created_at,
      created_by
    ) values (
               1,
               current_timestamp,
               updated_by,
               current_timestamp,
               updated_by);
    set site_id = LAST_INSERT_ID();
  ELSE
    update site set
                  version = version+1,
                  updated_at = current_timestamp,
                  updated_by = updated_by
    where site_id = site_id
      and version = version;
  END IF;
  select site_id;
END