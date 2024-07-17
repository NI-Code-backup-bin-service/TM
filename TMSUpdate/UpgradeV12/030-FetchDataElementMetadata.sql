--multiline;
CREATE PROCEDURE `fetch_data_element_metadata`(IN dataElementId INT)
BEGIN
  SELECT
    data_element_id,
    name,
    datatype,
    max_length,
    validation_expression,
    validation_message,
    front_end_validate,
    `unique`,
    `options`,
    displayname_en
  FROM data_element
  WHERE data_element_id = dataElementId;
END