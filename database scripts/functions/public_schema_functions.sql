-- register user
CREATE OR REPLACE FUNCTION public.register_user(
        p_username TEXT,
        p_first_name TEXT,
        p_last_name TEXT,
        p_email TEXT,
        p_phone_no TEXT,
        p_staff_id TEXT,
        p_status_id INT,
        p_last_login TEXT,
        p_institution_id INT
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        INSERT INTO user_accounts(username, first_name, last_name, email, phone_no, staff_id, status_id, last_login, institution_id) VALUES (p_username, p_first_name, p_last_name, p_email, p_phone_no, p_staff_id, p_status_id, p_last_login, p_institution_id);
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;
END;
$$ LANGUAGE plpgsql;

-- insert user's new password
CREATE OR REPLACE FUNCTION public.add_user_passwords(
        p_user_id INT,
        p_password_hash TEXT,
        p_requires_password_reset BOOLEAN,
        p_last_password_reset TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        INSERT INTO user_passwords(user_id, password_hash, requires_password_reset, last_password_reset) VALUES (p_user_id, p_password_hash::BYTEA, p_requires_password_reset, p_last_password_reset);
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;
END;
$$ LANGUAGE plpgsql;
