alter table profile_data add column approved int not null default 0;
update profile_data set approved = 1;