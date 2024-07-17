--multiline
CREATE PROCEDURE get_aquirer_page(IN searchTerm varchar(255), IN acquirers text, IN pageLimit INT, IN offSetValue INT)
BEGIN
    SELECT p.profile_id, p.name FROM profile p
    WHERE (UPPER(p.name) like searchTerm OR p.profile_id like searchTerm) AND FIND_IN_SET(p.name, acquirers) AND p.profile_type_id = (SELECT profile_type_id from profile_type WHERE name='acquirer')
    LIMIT pageLimit
    OFFSET offSetValue;
END;