-- +goose Up
-- 1. Completed pomodoro session
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (1, 'pomodoro', TRUE, 25, '2024-01-15 09:00:00', '2024-01-15 09:25:00');

-- 2. Completed short break
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (1, 'short break', TRUE, 5, '2024-01-15 09:25:00', '2024-01-15 09:30:00');

-- 3. In-progress pomodoro session (no end_time, not completed)
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time) 
VALUES (1, 'pomodoro', FALSE, 25, '2024-01-15 10:00:00');

-- 4. Completed long break after 4 pomodoros
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (1, 'long break', TRUE, 15, '2024-01-15 11:30:00', '2024-01-15 11:45:00');

-- 5. Custom duration pomodoro session
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (1, 'pomodoro', TRUE, 30, '2024-01-15 14:00:00', '2024-01-15 14:30:00'); 

-- User 2 Sample Data
-- 6. Completed pomodoro session for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'pomodoro', TRUE, 25, '2024-01-16 08:30:00', '2024-01-16 08:55:00');

-- 7. Completed short break for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'short break', TRUE, 5, '2024-01-16 08:55:00', '2024-01-16 09:00:00');

-- 8. Second completed pomodoro for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'pomodoro', TRUE, 25, '2024-01-16 09:00:00', '2024-01-16 09:25:00');

-- 9. Third completed pomodoro for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'pomodoro', TRUE, 25, '2024-01-16 09:30:00', '2024-01-16 09:55:00');

-- 10. Fourth completed pomodoro for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'pomodoro', TRUE, 25, '2024-01-16 10:00:00', '2024-01-16 10:25:00');

-- 11. Long break after 4 pomodoros for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'long break', TRUE, 15, '2024-01-16 10:25:00', '2024-01-16 10:40:00');

-- 12. In-progress pomodoro session for user 2 (no end_time, not completed)
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time) 
VALUES (2, 'pomodoro', FALSE, 25, '2024-01-16 11:00:00');

-- 13. Custom duration pomodoro session for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'pomodoro', TRUE, 45, '2024-01-16 14:00:00', '2024-01-16 14:45:00');

-- 14. Another day - completed pomodoro for user 2
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'pomodoro', TRUE, 25, '2024-01-17 09:15:00', '2024-01-17 09:40:00');

-- 15. Short break for user 2 on second day
INSERT INTO pomodoros (user_id, type, completed, session_duration, start_time, end_time) 
VALUES (2, 'short break', TRUE, 5, '2024-01-17 09:40:00', '2024-01-17 09:45:00');

-- +goose Down
delete from pomodoros;
