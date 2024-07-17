--multiline
INSERT IGNORE INTO data_element_locations
(
    profile_type_id,
    location_name,
    location_display_name
)
VALUES
    /*=================TIDs=================*/
(5, 'fraud', 'Fraud'),
    /*=================Sites=================*/
(4, 'tid_configuration', 'TID Configuration'),
(4, 'site_configuration', 'Site Configuration'),
(4, 'data_groups', 'Data Groups'),
(4, 'default_configuration', 'Default Configuration'),
(4, 'change_history', 'Change History'),
(4, 'users', 'Users'),
(4, 'fraud', 'Fraud'),
    /*=================Chains=================*/
(3, 'chain_configuration', 'Chain Configuration'),
(3, 'data_groups', 'Data Groups'),
(3, 'default_configuration', 'Default Configuration'),
(3, 'change_history', 'Change History'),
    /*=================Acquirers=================*/
(2, 'acquirer_configuration', 'Acquirer Configuration'),
(2, 'data_groups', 'Data Groups'),
(2, 'change_history', 'Change History');