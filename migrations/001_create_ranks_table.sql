
-- +goose up
CREATE TABLE if not exists ranks (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `min_xp` INT NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose down
DROP TABLE if exists ranks;