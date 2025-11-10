-- +goose Up
create table if not exists pomodoros (
    id int auto_increment primary key,
    user_id int not null,
    type ENUM('pomodoro','short break','long break') DEFAULT 'pomodoro', 
    completed BOOLEAN DEFAULT FALSE,
    session_duration int default 25,
    start_time datetime not null,
    end_time datetime,
    created_at datetime default current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
drop table if exists pomodoros;
