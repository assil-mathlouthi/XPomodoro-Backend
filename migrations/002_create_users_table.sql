
-- +goose up
create table if not exists users (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `email` VARCHAR(255) UNIQUE,
    `country` VARCHAR(50) ,
    `rank_id` INT NOT NULL references ranks(id) ON DELETE CASCADE,
    `xp` INT NOT NULL DEFAULT 0,
    `password_hash` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose down
DROP TABLE if exists users ;