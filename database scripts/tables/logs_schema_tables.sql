CREATE TABLE logs.user_tokens (
        token_id serial primary key,
        username TEXT,
        staff_id TEXT, 
        token TEXT,
        insti_code TEXT,
        app_code TEXT,
        created_at timestamp default current_timestamp,
);