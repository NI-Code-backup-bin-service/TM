--multiline
CREATE TABLE data_element_locations_data_element
(
    location_id INT NOT NULL,
    data_element_id INT NOT NULL,
    CONSTRAINT PK_data_element_locations_data_element PRIMARY KEY (location_id, data_element_id),
    CONSTRAINT FK_data_element_locations_profile_type_location_id FOREIGN KEY (location_id) REFERENCES data_element_locations(location_id),
    CONSTRAINT FK_data_element_locations_profile_type_data_element_id FOREIGN KEY (data_element_id) REFERENCES data_element(data_element_id),
    UNIQUE INDEX IDX_data_element_locations_data_element (location_id, data_element_id),
    INDEX IDX_data_element_locations_profile_type_location_id (location_id),
    INDEX IDX_data_element_locations_profile_type_data_element_id (data_element_id)
);