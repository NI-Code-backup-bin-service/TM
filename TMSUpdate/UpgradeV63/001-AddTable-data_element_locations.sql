--multiline
create table data_element_locations
(
    location_id INT AUTO_INCREMENT,
    profile_type_id INT NOT NULL,
    location_name VARCHAR(30) NOT NULL,
    location_display_name VARCHAR(30) NOT NULL,
    CONSTRAINT PK_data_element_locations PRIMARY KEY (location_id),
    CONSTRAINT FK_data_element_locations_profile_type_id FOREIGN KEY (profile_type_id) REFERENCES profile_type(profile_type_id),
    UNIQUE KEY UQ_profile_type_location_name (profile_type_id, location_name),
    UNIQUE INDEX IDX_data_element_locations (location_id),
    INDEX IDX_data_element_locations_profile_type_id (profile_type_id)
);