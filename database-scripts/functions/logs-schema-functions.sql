-- insert to user's token
CREATE OR REPLACE FUNCTION logs.create_user_token(
        p_username TEXT,
        p_staff_id TEXT,
        p_token TEXT,
        p_insti_code TEXT,
        p_app_code TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        INSERT INTO logs.user_tokens(username, staff_id, token, insti_code, app_code) VALUES (p_username, p_staff_id, p_token, p_insti_code, p_app_code);
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;    
END;
$$ LANGUAGE plpgsql;

-- update to user's token
CREATE OR REPLACE FUNCTION logs.update_user_token(
        p_username TEXT,
        p_staff_id TEXT,
        p_token TEXT,
        p_insti_code TEXT,
        p_app_code TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        UPDATE logs.user_tokens SET token = p_token, created_at = current_timestamp WHERE (username = p_username OR staff_id = p_staff_id) AND insti_code = p_insti_code AND app_code = p_app_code;
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;    
END;
$$ LANGUAGE plpgsql;
