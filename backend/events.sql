CREATE TABLE events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id VARCHAR(64) NOT NULL,
    action VARCHAR(16) NOT NULL,
    tracking_id VARCHAR(64),
    user_id VARCHAR(64),
    ip VARCHAR(45),
    created_at DATETIME,
    UNIQUE KEY unique_event (event_id)
);