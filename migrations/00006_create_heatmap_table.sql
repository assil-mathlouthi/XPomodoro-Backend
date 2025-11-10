-- +goose Up
CREATE TABLE if NOT exists heatmap (
    user_id INT NOT NULL,
    date DATE NOT NULL,
    count INT DEFAULT 0,
    PRIMARY KEY (user_id, date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +goose Down
DROP table if exists heatmap;
