--multiline;
CREATE PROCEDURE `get_change_approval_history`(IN afterDate VARCHAR(255), IN name VARCHAR(255), IN user VARCHAR(255), IN beforeDate VARCHAR(255), IN field VARCHAR(255))
BEGIN
  set @name = upper(concat('%', ifnull(name,''), '%'));
  set @user = upper(concat('%', ifnull(user,''), '%'));
  set @field = upper(concat('%', ifnull(field,''), '%'));
  SELECT
    a.approval_id as profile_data_id,
    p.name,
    de.name as 'field',
    a.current_value as original_value,
    a.new_value as updated_value,
    a.created_by as updated_by,
    a.created_at as updated_at,
    a.approved,
    a.approved_by as reviewed_by,
    a.approved_at as reviewed_at,
    a.tid_id
  FROM approvals a
         LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
         LEFT JOIN profile p ON p.profile_id = a.profile_id
  WHERE a.approved != 0
    AND (afterDate IS NULL OR a.created_at >= afterDate)
    AND p.name like @name
    AND a.approved_by like @user
    AND (beforeDate IS NULL OR a.created_at <= beforeDate)
    AND de.name like @field
  ORDER BY a.created_at DESC;
END