# The new column `not_overridable` may seem an unconventional name where `overridable` might seem more logical. It is
# named as such because if it were named overrideable then if we ever missed populating this field on an object in the
# server code then the field would not be overrideable and that is not the desired default behavior.
CALL AddColumn('profile_data', 'not_overridable', 'tinyint(1)', 'NOT NULL DEFAULT 0 COMMENT \'Determines if the data element can be overridden in the frontend\'');
UPDATE profile_data SET not_overridable = 1 WHERE profile_id = 1 AND data_element_id IN (select data_element_id from data_element where name in ('workstationNumber', 'secondaryTid'));