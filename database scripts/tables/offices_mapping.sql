-- create table for institutions 
CREATE TABLE offices_mapping.institutions (
    institution_id SERIAL PRIMARY KEY,
    institution_code TEXT NULL,
    institution_name TEXT NOT NULL,
    institution_description TEXT,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);