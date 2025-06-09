DROP TABLE Users;
DROP TABLE BlankItems;
DROP TABLE LyricsItems;
DROP TABLE MessageItems;
DROP TABLE SpeakerItems;
DROP TABLE TimerItems;
DROP TABLE Meetings;
DROP TABLE Organization;
DROP TABLE OrgUsers;

CREATE TABLE BlankItems (
    id CHAR(36) NOT NULL,
    meeting_id CHAR(36) NOT NULL,
    meeting_role VARCHAR(50) NOT NULL,
    item_type VARCHAR(20) NOT NULL DEFAULT 'blank',
    item_order INT NOT NULL,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
INSERT INTO BlankItems (id, meeting_id, meeting_role, item_type, item_order) VALUES ('c4ce7194-0f38-4b7b-89d1-09be87b902fd', '6cd5b59a-413a-4815-b3a9-e99a5dc91b50','Pre-meeting', 'blank', 0);

CREATE TABLE LyricsItems (
    id CHAR(36) NOT NULL,
    meeting_id CHAR(36) NOT NULL,
    meeting_role VARCHAR(50) NOT NULL,
    item_type VARCHAR(20) NOT NULL DEFAULT 'blank',
    item_order INT NOT NULL,
    hymn_id CHAR(36) NULL,
    show_translation TINYINT(0) NOT NULL DEFAULT 0,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
INSERT INTO LyricsItems (id, meeting_id, meeting_role, item_type, item_order, hymn_id, show_translation) VALUES ('5535277e-4192-4872-9320-c0f7a52569b0', '6cd5b59a-413a-4815-b3a9-e99a5dc91b50','Opening Hymn', 'lyrics', 4, 'bb125745-55eb-448c-b255-dac7ef6444cc', 1);

CREATE TABLE MessageItems (
    id CHAR(36) NOT NULL,
    meeting_id CHAR(36) NOT NULL,
    meeting_role VARCHAR(50) NOT NULL,
    item_type VARCHAR(20) NOT NULL DEFAULT 'blank',
    item_order INT NOT NULL,
    primary_text CHAR(200) NOT NULL,
    secondary_text CHAR(200) NULL,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TABLE SpeakerItems (
    id CHAR(36) NOT NULL,
    meeting_id CHAR(36) NOT NULL,
    meeting_role VARCHAR(50) NOT NULL,
    item_type VARCHAR(20) NOT NULL DEFAULT 'blank',
    item_order INT NOT NULL,
    speaker_name CHAR(200) NOT NULL,
    title CHAR(200) NULL,
    expected_duration INT NULL,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TABLE TimerItems (
    id CHAR(36) NOT NULL,
    meeting_id CHAR(36) NOT NULL,
    meeting_role VARCHAR(50) NOT NULL,
    item_type VARCHAR(20) NOT NULL DEFAULT 'blank',
    item_order INT NOT NULL,
    show_meeting_details TINYINT(1) NOT NULL DEFAULT 0,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
INSERT INTO TimerItems (id, meeting_id, meeting_role, item_type, item_order) VALUES ('b68d7dfa-d318-4fde-9381-f4992393c981', '6cd5b59a-413a-4815-b3a9-e99a5dc91b50', 'Meeting Countdown', 'timer', 1);

CREATE TABLE Users (
    id CHAR(36) NOT NULL,
    social_id VARCHAR(60) NULL,
    email VARCHAR(60) NOT NULL,
    first_name VARCHAR(50) NULL,
    full_name VARCHAR(50) NULL,
    last_name VARCHAR(50) NULL,
    photo_url VARCHAR(2000) NULL,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_email (email),
    UNIQUE KEY unique_social (social_id)
);
INSERT INTO `Users` (`id`, `email`, `first_name`, `full_name`, `last_name`) VALUES ('3cd5fe4e-9ecb-4ec2-b7c7-0d19288c08e0', 'pendenga@gmail.com', 'Grant', 'Grant Anderson', 'Anderson');
INSERT INTO `Users` (`id`, `email`, `first_name`, `full_name`, `last_name`) VALUES ('a5659535-43a8-486d-9b68-1da5d3fdee06', 'rskabelund@gmail.com', 'Randy', 'Randy Skabelund', 'Skabelund');

CREATE TABLE Meetings (
    id CHAR(36) NOT NULL,
    org_id CHAR(36) NOT NULL,
    conference VARCHAR(200) NULL,
    meeting VARCHAR(200) NOT NULL,
    meeting_date DATETIME NOT NULL,
    duration INT,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
INSERT INTO Meetings (id, org_id, conference, meeting, meeting_date, duration) VALUES ('6cd5b59a-413a-4815-b3a9-e99a5dc91b50', 'd65ad59c-216c-11f0-a191-ac1f6bbcd39a', 'Stake Conference', 'General Session', '2025-06-01 10:00:00.000000', '120');
INSERT INTO Meetings (id, org_id, meeting, meeting_date, duration) VALUES ('958a87d5-19b8-4e97-8016-dc9ca23072c5', 'e7d7a025-5bcd-43c8-ba35-e80d91ead4b2', 'Sacrament Meeting', '2025-04-27 09:00:00.000000', '120');
INSERT INTO Meetings (id, org_id, meeting, meeting_date, duration) VALUES ('f7c65b79-1d5a-45fc-a935-f3ed1bef75f9', '1b951e53-89d4-403c-b7e4-23984ac8aa15', 'Sacrament Meeting', '2025-04-27 09:00:00.000000', '120');

CREATE TABLE Organization (
    id CHAR(36) NOT NULL,
    name CHAR(200) NOT NULL,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_org_name (name)
);
INSERT INTO `Organization` (`id`, `name`) VALUES ('d65ad59c-216c-11f0-a191-ac1f6bbcd39a', 'Mesa Flatiron Stake');
INSERT INTO `Organization` (`id`, `name`) VALUES ('e7d7a025-5bcd-43c8-ba35-e80d91ead4b2', 'Boulder Mountain Ward');
INSERT INTO `Organization` (`id`, `name`) VALUES ('6b11edd1-1d25-4268-a6a2-6e564ebde510', 'Ironwood Ward');
INSERT INTO `Organization` (`id`, `name`) VALUES ('62a29b38-5255-43bc-b857-2480bf7b0de5', 'Ocotillo Ward');
INSERT INTO `Organization` (`id`, `name`) VALUES ('37a98239-bf69-430f-bc28-7af63edd52c7', 'Twin Knolls Ward');
INSERT INTO `Organization` (`id`, `name`) VALUES ('1b951e53-89d4-403c-b7e4-23984ac8aa15', 'Signal Butte 1st Ward');
INSERT INTO `Organization` (`id`, `name`) VALUES ('bf6b6624-71a5-49dd-8c89-59938606577b', 'Adobe Branch');

CREATE TABLE OrgUsers (
    org_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    deleted_dt DATETIME NULL,
    inserted_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY (org_id, user_id, deleted_dt)
);
INSERT INTO `OrgUsers` (`org_id`, `user_id`) VALUES ('d65ad59c-216c-11f0-a191-ac1f6bbcd39a', '3cd5fe4e-9ecb-4ec2-b7c7-0d19288c08e0');
INSERT INTO `OrgUsers` (`org_id`, `user_id`) VALUES ('e7d7a025-5bcd-43c8-ba35-e80d91ead4b2', '3cd5fe4e-9ecb-4ec2-b7c7-0d19288c08e0');
INSERT INTO `OrgUsers` (`org_id`, `user_id`) VALUES ('d65ad59c-216c-11f0-a191-ac1f6bbcd39a', 'a5659535-43a8-486d-9b68-1da5d3fdee06');
INSERT INTO `OrgUsers` (`org_id`, `user_id`) VALUES ('e7d7a025-5bcd-43c8-ba35-e80d91ead4b2', 'a5659535-43a8-486d-9b68-1da5d3fdee06');

/*
SELECT * FROM Users;
SELECT * FROM BlankItems;
SELECT * FROM LyricsItems;
SELECT * FROM MessageItems;
SELECT * FROM SpeakerItems;
SELECT * FROM TimerItems;
SELECT * FROM Meetings;
SELECT * FROM Organization;
SELECT * FROM OrgUsers;
*/