--multiline;
create table if not exists apk (
    apk_id int(11) NOT NULL AUTO_INCREMENT, 
    name varchar(45) unique, 
    PRIMARY KEY(apk_id)
);