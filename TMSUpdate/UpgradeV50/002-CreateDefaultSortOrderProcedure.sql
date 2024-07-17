--multiline
CREATE PROCEDURE SORTGROUPS()
BEGIN
	DECLARE totalDataGroups INT DEFAULT 0;
	DECLARE currentIndex INT DEFAULT 0;
	DECLARE currentDataGroupId INT DEFAULT 0;

	SELECT COUNT(1) FROM data_group INTO totalDataGroups;
	SET currentIndex = 0;

	WHILE currentIndex < totalDataGroups DO 
	  SELECT data_group_id FROM data_group LIMIT currentIndex, 1 INTO currentDataGroupId;
	  
	  UPDATE data_element de
	  SET sort_order_in_group = (SELECT y.position
								 FROM (SELECT x.data_element_id, x.position, x.name
									   FROM (SELECT de2.data_element_id, de2.name, @rownum := @rownum + 1 AS position
											 FROM data_element de2
											 JOIN (SELECT @rownum := 0) r
											 WHERE data_group_id = currentDataGroupId
											 ORDER BY de2.name) x
									   ) y
								 WHERE y.data_element_id = de.data_element_id
								)
	  WHERE data_group_id = currentDataGroupId;
	  
	  SET currentIndex = currentIndex + 1;
	END WHILE;
END