DELETE pdg.* FROM `profile_data_group` pdg INNER JOIN `data_group` dg ON pdg.data_group_id = dg.data_group_id WHERE (dg.name = "loyalty");
DELETE FROM `data_group` WHERE `name` = 'loyalty';