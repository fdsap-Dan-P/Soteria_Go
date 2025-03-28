-- system paramters 
create table parameters.system_config(
    config_id serial primary key,
    config_code text not null,
    config_name text not null,
    config_description text null,
    config_type text,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

create table parameters.insti_app_config(
    config_id int references parameters.system_config(config_id) on delete cascade,
    config_value text not null,
    config_insti_code text,
    config_app_code text,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

insert into parameters.insti_app_config
(config_id, config_value, config_insti_code, config_app_code)
values 
(1, '1440', '2162', 'OR0003-1732775142'),
(2, '90', '2162', 'OR0003-1732775142'),
(3, 'true', '2162', 'OR0003-1732775142'),
(4, 'true', '2162', 'OR0003-1732775142'),
(5, 'true', '2162', 'OR0003-1732775142'),
(6, 'true', '2162', 'OR0003-1732775142'),
(7, '8', '2162', 'OR0003-1732775142'),
(8, '5', '2162', 'OR0003-1732775142')


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
    ('116', 'Validation Failed', 'Terminated OTP', false), -- replace
    ('117', 'Validation Failed', 'Invalid OTP', false),
    ('118', 'Validation Failed', 'Expired OTP', false),
    ('119', 'Validation Failed', 'Password Not Match', false),
    ('120', 'Validation Failed', 'Invalid Date Range', false),
    ('121', 'Validation Failed', 'No Data Available', false),
    ('122', 'Validation Failed', 'Unlocked: Password Reset Required', false),
    ('123', 'Validation Failed', 'Inactive User Account', false), -- dynamic
    ('124', 'Validation Failed', 'Device Currently in Use by Another User', false),
    ('125', 'Validation Failed', 'Role In Use', false),
    ('126', 'Validation Failed', 'Temporary Password Sent Not Used', false),
    ('200', 'Successful', '', true),
    ('201', 'Successful', 'Successfully Logged In', true),
    ('202', 'Successful', 'Successfully Logged Out', true),
    ('203', 'Successful', 'Successfully Inserted', true),
    ('204', 'Successful', 'Successfully Updated', true),
    ('205', 'Successful', 'Successfully Downloaded', true),
    ('206', 'Successful', 'Successfully Sent OTP', true ),
    ('207', 'Successful', 'Successfully Validated', true ),
    ('208', 'Successful', 'Successfully Sent Email', true ),
	('209', 'Successful', 'Successfully Sent SMS', true ),
	('210', 'Successful', 'Successfully Deleted', true ),
    ('211', 'Successful', 'Successfully Activate User Account', true ),
    ('212', 'Successful', 'Successfully Deactivate User Account', true ),
    ('213', 'Successful', 'Successfully Unblock User Account', true ),
    ('214', 'Successful', 'Successfully Unlock User Account', true ),
    ('215', 'Successful', 'Successfully Validated Headers', true ),
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