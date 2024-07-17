--multiline;
CREATE PROCEDURE discard_change(IN profile_data_id int)
BEGIN
	delete pd from profile_data pd WHERE pd.profile_data_id = profile_data_id;
END;