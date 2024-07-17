SET @versionId = (SELECT package_id FROM  package WHERE version = 307);
SET @versionDataElement = (SELECT data_element_id FROM data_element WHERE name = "RequiredSoftwareVersion");
call bulkUpdate(@versionId, @versionDataElement);