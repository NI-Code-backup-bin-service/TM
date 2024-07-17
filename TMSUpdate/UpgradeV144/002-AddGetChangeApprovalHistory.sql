--multiline
CREATE PROCEDURE `get_change_approval_history`(
  IN afterDate VARCHAR(255),
  IN name VARCHAR(255),
  IN user VARCHAR(255),
  IN beforeDate VARCHAR(255),
  IN field VARCHAR(255),
  IN acquirers TEXT
)
BEGIN
  SET @name  = upper(concat('%', ifnull(name,  ''), '%'));
  SET @user  = upper(concat('%', ifnull(user,  ''), '%'));
  SET @field = upper(concat('%', ifnull(field, ''), '%'));

SELECT
  p.name,
  CONCAT(dg.name, "/", de.name) AS 'field',
  a.current_value AS original_value,
  a.new_value AS updated_value,
  a.created_by AS updated_by,
  a.created_at AS updated_at,
  a.approved,
  a.approved_by AS reviewed_by,
  a.approved_at AS reviewed_at,
  a.tid_id,
  a.merchant_id,
  a.is_password,
  a.is_encrypted

FROM
  approvals a
    LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
    LEFT JOIN data_group dg ON dg.data_group_id = de.data_group_id
    LEFT JOIN profile p ON p.profile_id = a.profile_id


    LEFT JOIN chain_profiles cp ON cp.chain_profile_id = p.profile_id
    LEFT JOIN profile p2 ON p2.profile_id = cp.acquirer_id


    LEFT JOIN chain_profiles cp2 ON cp2.acquirer_id = p.profile_id
    LEFT JOIN profile p3 ON p3.profile_id = cp2.acquirer_id


    LEFT JOIN (
      profile p4
      JOIN site_profiles sp ON sp.profile_id = p4.profile_id
      JOIN site_profiles sp2 ON sp2.site_id = sp.site_id
      JOIN profile p5 ON p5.profile_id = sp2.profile_id
    ) ON p5.profile_id = a.profile_id
    AND p5.profile_type_id = 4


    LEFT JOIN tid td ON td.tid_id = p.name
    LEFT JOIN tid_site ts ON ts.tid_id = p.name
    LEFT JOIN site t ON t.site_id = ts.site_id
    LEFT JOIN (
      site_profiles tp4
      JOIN profile p6 ON p6.profile_id = tp4.profile_id
      JOIN profile_type pt4 ON pt4.profile_type_id = p6.profile_type_id
      AND pt4.priority = 4
    ) ON tp4.site_id = t.site_id
    AND td.tid_id != 0

WHERE
  a.approved != 0
  AND (
    afterDate IS NULL
    OR a.created_at >= afterDate
  )
  AND p.name LIKE @name
  AND a.approved_by LIKE @user
  AND (
    beforeDate IS NULL
    OR a.created_at <= beforeDate
  )
  AND de.name LIKE @field
  AND (
    FIND_IN_SET(p2.name, acquirers)
    OR FIND_IN_SET(p3.name, acquirers)
    OR (
      FIND_IN_SET(p4.name, acquirers)
      OR FIND_IN_SET(a.acquirer, acquirers)
    )
    OR (
      FIND_IN_SET(p6.name, acquirers)
      OR FIND_IN_SET(a.acquirer, acquirers)
    )
  )

GROUP BY
  p.name,
  dg.name,
  de.name,
  original_value,
  updated_value,
  updated_by,
  updated_at,
  a.approved,
  reviewed_by,
  reviewed_at,
  a.tid_id,
  a.merchant_id,
  a.is_password,
  a.is_encrypted

ORDER BY
  a.approved_at DESC;
END;