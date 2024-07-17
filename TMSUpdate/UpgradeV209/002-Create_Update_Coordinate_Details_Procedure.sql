--multiline
CREATE PROCEDURE `update_coordinate_details`(IN tid int, IN update_coordinates varchar(50), IN update_accuracy varchar(50), IN update_coordinate_time DATETIME)
BEGIN
        UPDATE tid t SET t.coordinates = update_coordinates,
                         t.accuracy = update_accuracy,
                         t.last_coordinate_time = update_coordinate_time
        WHERE t.tid_id = tid;
END