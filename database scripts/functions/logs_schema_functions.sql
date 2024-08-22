-- to reset user bad log attemptsz
CREATE OR REPLACE FUNCTION logs.delete_from_login_logs( 
	p_user_id INT, 
	p_is_success BOOLEAN 
) 
RETURNS TEXT AS $$ 
DECLARE v_error_message TEXT; 
BEGIN 
	BEGIN 
		DELETE FROM logs.login_logs WHERE user_id = p_user_id AND is_success = p_is_success; 
	RETURN 'Success'; 
	EXCEPTION WHEN OTHERS THEN -- Capture the error message 
	v_error_message := SQLERRM; RETURN 'Deleting Data Failed: ' || v_error_message; 
	END; 
END; $$ 
LANGUAGE plpgsql;

-- to log user attempt
CREATE OR REPLACE FUNCTION logs.add_login_logs(
        p_user_id INT,
        p_ip TEXT,
        p_login_attempt_at TEXT,
        p_is_success BOOLEAN
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        INSERT INTO logs.login_logs(user_id, ip, login_attempt_at, is_success) VALUES (p_user_id, p_ip, p_login_attempt_at, p_is_success);
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;
END;
$$ LANGUAGE plpgsql;

-- insert to audit logs
CREATE OR REPLACE FUNCTION logs.add_audit_logs(
        p_user_id INT,
        p_staff_id TEXT,
        p_username TEXT,
        p_user_status TEXT,
        p_institution_name TEXT,
        p_user_action TEXT,
        p_old_value TEXT,
        p_new_value TEXT,
        p_log_date TEXT,
        p_log_time TEXT 
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        INSERT INTO logs.audit_logs(user_id, staff_id, username, user_status, institution_name, user_action, old_value, new_value, log_date, log_time) VALUES (p_user_id, p_staff_id, p_username, p_user_status, p_institution_name, p_user_action, p_old_value, p_new_value, p_log_date, p_log_time);
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;    
END;
$$ LANGUAGE plpgsql;

-- insert to locked users
CREATE OR REPLACE FUNCTION logs.add_locked_users(
        p_user_id INT,
        p_user_ip TEXT,
        p_created_at TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
BEGIN
    BEGIN
        INSERT INTO logs.locked_users(user_id, ip, created_at) VALUES (p_user_id, p_user_ip, p_created_at);
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;    
END;
$$ LANGUAGE plpgsql;

-- insert to user's token
CREATE OR REPLACE FUNCTION logs.(
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
        UPDATE logs.user_tokens SET token = p_token, updated_at = current_timestamp WHERE (username = p_username OR staff_id = p_staff_id) AND insti_code = p_insti_code AND app_code = p_app_code;
        RETURN 'Success';
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed: ' || v_error_message;
    END;    
END;
$$ LANGUAGE plpgsql;
