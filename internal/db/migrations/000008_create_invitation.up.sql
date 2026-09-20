CREATE TABLE IF NOT EXISTS USER_INVITATIONS (
    token bytea PRIMARY KEY, 
    user_id bigint NOT NULL 
); 