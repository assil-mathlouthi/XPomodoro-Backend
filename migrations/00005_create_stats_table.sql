-- +goose Up
create table if not exists stats (
    user_id int primary key,
    longest_streak int default 0,
    current_streak int default 0,
    last_updated datetime default current_timestamp,
    created_at datetime default current_timestamp,
    foreign key (user_id) references users(id) on delete cascade
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +goose Down
DROP table if exists stats;
