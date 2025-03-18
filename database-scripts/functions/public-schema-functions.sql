-- register user
CREATE OR REPLACE FUNCTION public.register_user(p_username text, p_first_name text, p_middle_name text, p_last_name text, p_email text, p_phone_no text, p_staff_id text, p_institution_id integer, p_password_hash text, p_requires_password_reset boolean, p_last_password_reset text, p_birthdate text, p_insti_code text, p_app_code text, p_app_id integer)
 RETURNS text
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_error_message TEXT;
    userId INT;
BEGIN
    BEGIN
        INSERT INTO user_accounts(username, first_name, middle_name, last_name, email, phone_no, staff_id, institution_id, birthdate, application_id) VALUES (p_username, p_first_name, p_middle_name, p_last_name, p_email, p_phone_no, p_staff_id, p_institution_id, p_birthdate::DATE, p_application_id) RETURNING user_id INTO userId;
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed [user_accounts]: ' || v_error_message;
    END;

    -- set user's password
    BEGIN
        INSERT INTO user_passwords(user_id, password_hash, requires_password_reset, last_password_reset, insti_code, app_code) VALUES (userId, p_password_hash::BYTEA, p_requires_password_reset, p_last_password_reset, p_insti_code, p_app_code);
    EXCEPTION
        WHEN OTHERS THEN
            -- Capture the error message
            v_error_message := SQLERRM;
            RETURN 'Inserting Data Failed [user_passwords]: ' || v_error_message;
    END;

    RETURN 'Success';
END;
$function$
;

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
