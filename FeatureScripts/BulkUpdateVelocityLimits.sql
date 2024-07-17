set @NumTxnsDaily = 10;
set @TotalLimitDaily = 20;
set @NumOfTxnsBatch = 30;
set @TotalLimitBatch = 40;
set @TxnAmount = 50;
set @CleanseTime = "00:00";
-- add/amend profile data to all profiles for cleanse time, add a new line for each profile
-- configurable parameters:
-- profile_id - change to the profile of choice
-- datavalue - The desired time to set
insert into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (2735, (select data_element_id from data_element where name = "dailyTxnCleanseTime"), @CleanseTime, 1, NOW(), "Feature Script", NOW(), "Feature Script", 1, 0, 0) on duplicate key update datavalue = @CleanseTime;
insert into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (2736, (select data_element_id from data_element where name = "dailyTxnCleanseTime"), @CleanseTime, 1, NOW(), "Feature Script", NOW(), "Feature Script", 1, 0, 0) on duplicate key update datavalue = @CleanseTime;

-- add/amend entries for each profile id into velocity_limits, add new line for each profile, velocity_limit_id must be unique
-- configurable parameters:
-- velocity_limit_id - Must be unique, can consist of alphanumeric characters
-- site_id - update the profile id here
-- transaction_limit_daily - new daily txn limit to set
-- transaction_limit_batch - new batch txn limit to set
-- single_transaction_limit - new single txn amount limit to set
-- cumulative_daily - new total daily txn amount limit to set
-- cumulative_batch - new total batch txn amount to set
insert into velocity_limits(velocity_limit_id, site_id, tid_id, limit_level, scheme, transaction_limit_daily, transaction_limit_batch, single_transaction_limit, cumulative_daily, cumulative_batch) values ("fs1", (select site_id from site_profiles where profile_id = 2735), -1, 3, NULL, @NumTxnsDaily, @NumOfTxnsBatch, @TxnAmount, @TotalLimitDaily, @TotalLimitBatch) on duplicate key update transaction_limit_daily = @NumTxnsDaily, transaction_limit_batch = @NumOfTxnsBatch, single_transaction_limit = @TxnAmount, cumulative_daily = @TotalLimitDaily, cumulative_batch = @TotalLimitBatch;
insert into velocity_limits(velocity_limit_id, site_id, tid_id, limit_level, scheme, transaction_limit_daily, transaction_limit_batch, single_transaction_limit, cumulative_daily, cumulative_batch) values ("fs2", (select site_id from site_profiles where profile_id = 2736), -1, 3, NULL, @NumTxnsDaily, @NumOfTxnsBatch, @TxnAmount, @TotalLimitDaily, @TotalLimitBatch) on duplicate key update transaction_limit_daily = @NumTxnsDaily, transaction_limit_batch = @NumOfTxnsBatch, single_transaction_limit = @TxnAmount, cumulative_daily = @TotalLimitDaily, cumulative_batch = @TotalLimitBatch;