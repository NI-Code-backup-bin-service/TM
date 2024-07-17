--multiline;
CREATE PROCEDURE `check_tid_exists_for_mid`(IN mid text,IN tid text)
BEGIN
select count(*) from tid_site_profiles tsp
		join site_profiles sp ON sp.site_id =  tsp.site_id
		join profile_data pd ON sp.profile_id = pd.profile_id
        where pd.datavalue = mid AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'merchantNo')
		AND tsp.tid_id = tid;
END