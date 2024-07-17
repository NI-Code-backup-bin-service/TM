--multiline;
CREATE PROCEDURE `get_packages`()
BEGIN
  SELECT *
  FROM package p
  order by p.package_id desc;
END