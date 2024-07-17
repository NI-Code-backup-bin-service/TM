--multiline;
CREATE PROCEDURE `delete_site`(IN siteId INT)
BEGIN
  START TRANSACTION;
  SET @profileID = (SELECT sp.profile_id
                    FROM site_profiles sp
                           LEFT JOIN profile p ON p.profile_id = sp.profile_id
                           LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
                    WHERE sp.site_id = siteId
                    ORDER BY pt.priority
                    LIMIT 1);
  DELETE FROM profile_data WHERE profile_id = @profileID AND @profileID IS NOT NULL;
  DELETE FROM profile_data_group WHERE profile_id = @profileID AND @profileID IS NOT NULL;
  DELETE FROM site_profiles WHERE site_id = siteId AND siteId IS NOT NULL;
  DELETE FROM tid_site WHERE site_id = siteId AND siteId IS NOT NULL;
  DELETE FROM site_level_users where site_id = siteId AND siteId IS NOT NULL;
  DELETE FROM site WHERE site_id = siteId AND siteId IS NOT NULL;
  DELETE u FROM tid_user_override u LEFT OUTER join tid_site ts ON ts.tid_id = u.tid_id WHERE ts.tid_id IS NULL;
  DELETE t FROM tid t LEFT OUTER join tid_site ts ON ts.tid_id = t.tid_id WHERE ts.tid_id IS NULL;
  COMMIT;
END