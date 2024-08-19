-- system paramters 
create table paramaters.system_config(
    config_id serial primary key,
    config_code text not null,
    config_name text not null,
    config_description text null,
    config_type text,
    config_value text not null,
    config_insti_code text not null,
    config_app_code text not null,
    created_at timestamp default current_timestamp,
    updated_at_at timestamp default current_timestamp,
);

insert into parameters.system_config
(config_code, config_name, config_description, config_type, config_value, config_insti_code, config_app_code)
VALUES 
('jwt', 'Token', 'Time', 'Minute', '1440', 'rbi', 'data-plaform'),
('jwt', 'Token', 'Time', 'Minute', '1440', 'rbi', 'cagabay'),

('lockout_attempt', 'Password Lockout Attempt', 'Minimum login attempts', 'Numeric', '3', 'rbi', 'data-plaform'),
('lockout_attempt', 'Password Lockout Attempt', 'Minimum login attempts', 'Numeric', '3', 'rbi', 'cagabay'),

('lockout_period', 'Password Lockout Period', 'Lockout period in minutes', 'Numeric', '5', 'rbi', 'data-plaform'),
('lockout_period', 'Password Lockout Period', 'Lockout period in minutes', 'Numeric', '5', 'rbi', 'cagabay')

('pass_exp', 'Password Expiration', 'Password Active Duration Starting Password Creation', 'Day', '90', 'rbi', 'data-plaform'),
('pass_exp', 'Password Expiration', 'Password Active Duration Starting Password Creation', 'Day', '90', 'rbi', 'cagabay'),

('pass_sym', 'Symbol Validation', 'Ensure the password contains at least one symbol character', 'bool', 'true', 'rbi', 'data-plaform'),
('pass_sym', 'Symbol Validation', 'Ensure the password contains at least one symbol character', 'bool', 'true', 'rbi', 'cagabay'),

('pass_upc', 'Uppercase Validation', 'Ensure the password contains at least one uppercase character', 'bool', 'true', 'rbi', 'data-plaform'),
('pass_upc', 'Uppercase Validation', 'Ensure the password contains at least one uppercase character', 'bool', 'true', 'rbi', 'cagabay'),

('pass_lwr', 'Lowercase Validation', 'Ensure the password contains at least one lowercase character', 'bool', 'true', 'rbi', 'data-plaform'),
('pass_lwr', 'Lowercase Validation', 'Ensure the password contains at least one lowercase character', 'bool', 'true', 'rbi', 'cagabay'),

('pass_num', 'Number Validation', 'Ensure the password contains at least one number character', 'bool', 'true', 'rbi', 'data-plaform'),
('pass_num', 'Number Validation', 'Ensure the password contains at least one number character', 'bool', 'true', 'rbi', 'cagabay'),

('pass_min', 'Minimum Character', 'Setting password minimum character requirement', 'Numeric', '8', 'rbi', 'data-plaform'),
('pass_min', 'Minimum Character', 'Setting password minimum character requirement', 'Numeric', '8', 'rbi', 'cagabay'),

('via_email', 'Via Email', 'Reset Password via redirect link', 'bool', 'true', 'rbi', 'data-plaform'),
('via_email', 'Via Email', 'Reset Password via redirect link', 'bool', 'true', 'rbi', 'cagabay'),

('via_sms', 'Via SMS', 'Reset Password via SMS', 'bool', 'true', 'rbi', 'data-plaform'),
('via_sms', 'Via SMS', 'Reset Password via SMS', 'bool', 'true', 'rbi', 'cagabay'),

('pass_reuse', 'Password Reusability', 'Minimum Password Reusability', 'Numeric', '3', 'rbi', 'data-plaform'),
('pass_reuse', 'Password Reusability', 'Minimum Password Reusability', 'Numeric', '3', 'rbi', 'cagabay');

CREATE TABLE parameters.return_message (
    retcode TEXT,
    category TEXT,
    error_message TEXT,
    is_success BOOLEAN,
    created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

INSERT INTO parameters.return_message(
    retcode, category, error_message, is_success)
    VALUES
    ('100', 'Validation Failed', '', false),
    ('101', 'Validation Failed', 'First Login: Password Reset Required', false),
    ('102', 'Validation Failed', 'Password Expired', false),
    ('103', 'Validation Failed', 'Invalid Password', false),
    ('104', 'Validation Failed', 'Invalid Token', false),
    ('105', 'Validation Failed', 'User Already Logged In', false),
    ('106', 'Validation Failed', 'Password Cannot Be Reused', false),
    ('107', 'Validation Failed', 'Invalid Date', false),
    ('108', 'Validation Failed', 'Invalid Phone number', false),
    ('109', 'Validation Failed', 'Invalid Email Address', false),
    ('110', 'Validation Failed', 'Expired Token', false),
    ('111', 'Validation Failed', 'Token Missing', false),
    ('112', 'Validation Failed', 'Invalid Staff Id', false),
    ('113', 'Validation Failed', 'Password Reuse Prohibited', false), -- should be dynamic
    ('114', 'Validation Failed', 'User Lockout Period', false),
    ('115', 'Validation Failed', 'Session Id Missing', false),
    ('116', 'Validation Failed', 'Email Reuse Prohibited', false), -- replace
    ('117', 'Validation Failed', 'Invalid OTP', false),
    ('118', 'Validation Failed', 'Expired OTP', false),
    ('119', 'Validation Failed', 'Password Not Match', false),
    ('120', 'Validation Failed', 'Invalid Date Range', false),
    ('121', 'Validation Failed', 'No Data Available', false),
    ('122', 'Validation Failed', 'Unlocked: Password Reset Required', false),
    ('123', 'Validation Failed', 'Inactive User Account', false), -- dynamic
    ('124', 'Validation Failed', 'Device Currently in Use by Another User', false),
    ('125', 'Validation Failed', 'Role In Use', false),
    ('126', 'Validation Failed', 'Unlocked: Temporary Password Sent', false),
    ('200', 'Successful', '', NULL),
    ('201', 'Successfully Logged In', '', NULL),
    ('202', 'Successfully Logged Out', '', NULL),
    ('203', 'Successfully Inserted', '', NULL),
    ('204', 'Successfully Updated', '', NULL),
    ('205', 'Successfully Downloaded', '', NULL),
    ('206', 'Successfully Sent OTP', '', NULL ),
    ('207', 'Successfully Validated', '', NULL ),
    ('208', 'Successfully Sent Email', '', NULL ),
	('209', 'Successfully Sent SMS', '', NULL ),
	('210', 'Successfully Deleted', '', NULL ),
    ('211', 'Successfully Activate User Account', '', NULL ),
    ('212', 'Successfully Deactivate User Account', '', NULL ),
    ('213', 'Successfully Unblock User Account', '', NULL ),
    ('214', 'Successfully Unlock User Account', '', NULL ),
    ('300', 'Internal Server Error', '', false),
    ('301', 'Internal Server Error', 'Parsing Data Failed', false),
    ('302', 'Internal Server Error', 'Fetching Data Failed', false),
    ('303', 'Internal Server Error', 'Inserting Data Failed', false),
    ('304', 'Internal Server Error', 'Updating Data Failed', false),
    ('305', 'Internal Server Error', 'Token Generation Failed', false),
    ('306', 'Internal Server Error', 'File Download Failed', false),
    ('307', 'Internal Server Error', 'File Creation Failed', false),
    ('308', 'Internal Server Error', 'URL Retrieval Failed', false),
    ('309', 'Internal Server Error', 'Unable to Open File', false),
    ('310', 'Internal Server Error', 'Unmarshalling Failed', false),
    ('311', 'Internal Server Error', 'Marshalling Failed', false),
    ('300', 'Internal Server Error', 'Timezone Loading Failed', false),
    ('312', 'Internal Server Error', 'Timezone Loading Failed', false),
    ('313', 'Internal Server Error', 'HTML Parsing Failed', false),
    ('314', 'Internal Server Error', 'Deleting Data Failed', false),
    ('315', 'Internal Server Error', 'Send Emails Failed', false),
    ('316', 'Internal Server Error', 'Logging in to Third Party Failed', false ),
    ('317', 'Internal Server Error', 'Reading Third Party Response Failed', false ),
    ('400', 'Bad Request', '', false),
    ('401', 'Bad Request', 'Required Input Missing', false),
    ('402', 'Bad Request', 'Unauthorized Access', false),
    ('403', 'Bad Request', 'Unique Field Already Exists', false),
    ('404', 'Bad Request', 'Not Found', false),
    ('405', 'Bad Request', 'Request Failed', false);