--multiline
CREATE PROCEDURE get_aquirer_count(IN search_term varchar(255), IN acquirers text)
BEGIN
    SELECT COUNT(*)
    FROM profile p
    WHERE (UPPER(p.name) like search_term OR p.profile_id like search_term) AND FIND_IN_SET(p.name, acquirers) AND p.profile_type_id = (SELECT profile_type_id from profile_type WHERE name='acquirer');
END;