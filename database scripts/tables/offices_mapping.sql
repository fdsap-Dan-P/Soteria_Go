-- create table for institutions 
CREATE TABLE offices_mapping.institutions (
    institution_id SERIAL PRIMARY KEY,
    institution_code TEXT NULL,
    institution_name TEXT NOT NULL,
    institution_description TEXT,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

insert into offices_mapping.institutions
(institution_code, institution_name, institution_description)
values
('2', 'CARD, Inc.',''),
('1409', 'CARD Bank, Inc.',''),
('2039', 'CARD MBA, Inc.',''),
('2137', 'CMDI, Inc.',''),
('2162', 'CARD SME Bank, Inc. A Thrift Bank',''),
('2482', 'CAMIA, Inc.',''),
('2492', 'CARD BDSF, Inc.',''),
('2500', 'CMIT, Inc.',''),
('2528', 'BotiCARD, Inc.',''),
('2546', 'CARD MRI Rizal Bank, Inc.',''),
('2686', 'CLFC',''),
('2693', 'RISE',''),
('2699', 'MLNI',''),
('2705', 'CARD EMPC',''),
('5399', 'CARD MRI Hijos Tours, Inc.',''),
('5492', 'CARD MRI Publishing House, Inc.',''),
('5556', 'CMPMI',''),
('6773', 'CARD MRI Astro Laboratories, Inc.',''),
('7536', 'FDSAP',''),
('9039', 'CARD Pioneer Microinsurance, Inc.',''),
('9272', 'CARD Masikhay Consultancy Services, Inc.',''),
('9659', 'CARD Indogrosir, Inc.',''),
('9660', 'CARD Ottokonek, Inc.',''),
('9857', 'PADAYON Microfinance Inc.',''),
('9863', 'BENTE Productions Inc.',''),
('9869', 'Bakawan Data Analytics, Inc.',''),
('9875', 'CARD Clinics and Allied Services Inc. (CCASI)','');