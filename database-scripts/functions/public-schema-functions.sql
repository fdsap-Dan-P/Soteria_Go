-- register user
CREATE OR REPLACE FUNCTION public.register_user(
        p_username TEXT,
        p_first_name TEXT,
        p_middle_name TEXT,
        p_last_name TEXT,
        p_email TEXT,
        p_phone_no TEXT,
        p_staff_id TEXT,
        p_institution_id INT,

        p_password_hash TEXT,
        p_requires_password_reset BOOLEAN,
        p_last_password_reset TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_error_message TEXT;
    userId INT;
BEGIN
    BEGIN
        INSERT INTO user_accounts(username, first_name, middle_name, last_name, email, phone_no, staff_id, institution_id) VALUES (p_username, p_first_name, p_middle_name, p_last_name, p_email, p_phone_no, p_staff_id, p_institution_id) RETURNING user_id INTO userId;
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed [user_accounts]: ' || v_error_message;
    END;

    -- set user's password
    BEGIN
        INSERT INTO user_passwords(user_id, password_hash, requires_password_reset, last_password_reset) VALUES (userId, p_password_hash::BYTEA, p_requires_password_reset, p_last_password_reset);
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed [user_passwords]: ' || v_error_message;
    END;

    RETURN 'Success';
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


-- get user's password reuse
CREATE OR REPLACE FUNCTION public.password_reuse(
    user_id numeric, 
    reuse_limit numeric, 
    insti_code text, 
    app_code text
)
RETURNS TABLE(password_hash bytea)
LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    EXECUTE
        'SELECT password_hash FROM (
            SELECT * FROM user_passwords
            WHERE user_id = $1
              AND last_password_reset IS NOT NULL
              AND last_password_reset <> ''''  -- Handles empty string case for TEXT columns
              AND insti_code = $2
              AND app_code = $3
            ORDER BY created_at DESC
            LIMIT $4
        ) subq'
    USING user_id, insti_code, app_code, reuse_limit;
END;
$function$;
