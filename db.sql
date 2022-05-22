CREATE
DATABASE blog;
USE
blog;
/* All 9 Tables: article visitor upload download ip_location admin comment login search */

/* Article */
CREATE TABLE article
(
    id             INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title          VARCHAR(100) NOT NULL,
    subtitle       VARCHAR(100),
    tag            VARCHAR(30)  NOT NULL,
    content_md     MEDIUMTEXT   NOT NULL,
    content_html   MEDIUMTEXT,
    content_text   MEDIUMTEXT,
    author         VARCHAR(30)  NOT NULL,
    author_id      INT          NOT NULL,
    read_count     INT          NOT NULL,
    create_date    DATETIME     NOT NULL,
    last_edit_date DATETIME
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8mb4;

/* Visitor */
CREATE TABLE visitor
(
    id         INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    ip         VARCHAR(20) NOT NULL,
    location   VARCHAR(50),
    visit_date DATETIME    NOT NULL,
    platform   VARCHAR(20),
    browser    VARCHAR(20),
    version    VARCHAR(20)
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8;

/* Upload */
CREATE TABLE upload
(
    id          INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    filename    VARCHAR(100) NOT NULL,
    ip          VARCHAR(20)  NOT NULL,
    location    VARCHAR(50),
    upload_date DATETIME     NOT NULL,
    platform    VARCHAR(20),
    browser     VARCHAR(20),
    version     VARCHAR(20)
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8;

/* Download */
CREATE TABLE download
(
    id            INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    filename      VARCHAR(100) NOT NULL,
    ip            VARCHAR(20)  NOT NULL,
    location      VARCHAR(50),
    download_date DATETIME     NOT NULL,
    platform      VARCHAR(20),
    browser       VARCHAR(20),
    version       VARCHAR(20)
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8;

/* Ip_location */
CREATE TABLE ip_location
(
    id       INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    ip       VARCHAR(20) NOT NULL,
    location VARCHAR(50) NOT NULL
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8;

/* Admin */
CREATE TABLE admin
(
    id            INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username      VARCHAR(30) NOT NULL UNIQUE,
    password      VARCHAR(30) NOT NULL,
    nickname      VARCHAR(30) NOT NULL,
    sex           CHAR(1),    NOT NULL,
    authority     TINYINT     NOT NULL,
    avatar        VARCHAR(50),
    register_date DATETIME    NOT NULL
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8;

/* Comment */
CREATE TABLE comment
(
    id           INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    article_id   INT         NOT NULL,
    reply_name   VARCHAR(30) NOT NULL,
    content      MEDIUMTEXT  NOT NULL,
    comment_date DATETIME    NOT NULL,
    platform     VARCHAR(20),
    browser      VARCHAR(20),
    version      VARCHAR(20),
    ip           VARCHAR(20) NOT NULL,
    location     VARCHAR(50),
    reviewer_id  INT         NOT NULL
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8mb4;

/* Login */
CREATE TABLE login
(
    id         INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    ip         VARCHAR(20) NOT NULL,
    location   VARCHAR(50),
    username   VARCHAR(30) NOT NULL,
    password   VARCHAR(30) NOT NULL,
    login_flag TINYINT     NOT NULL,
    login_date DATETIME    NOT NULL,
    platform   VARCHAR(20),
    browser    VARCHAR(20),
    version    VARCHAR(20)
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8;

/* Search */
CREATE TABLE search
(
    id          INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    ip          VARCHAR(20) NOT NULL,
    location    VARCHAR(50),
    keyword     VARCHAR(30) NOT NULL,
    search_date DATETIME    NOT NULL,
    platform    VARCHAR(20),
    browser     VARCHAR(20),
    version     VARCHAR(20)
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8mb4;

/* Dictionary */
CREATE TABLE dictionary
(
    id          INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    key_        VARCHAR(30)  NOT NULL,
    value_      VARCHAR(200) NOT NULL,
    key_type    VARCHAR(10),
    value_type  VARCHAR(10),
    create_date DATETIME     NOT NULL
) AUTO_INCREMENT = 1,DEFAULT CHARSET = utf8mb4;


INSERT INTO admin(username, password, nickname, sex, authority, register_date, avatar)
VALUES ("admin", "123456", "Admin", "女", 1, CURRENT_TIMESTAMP, "head.jfif");
INSERT INTO article(title, subtitle, tag, content_md, create_date, author, author_id, read_count)
VALUES ("这是文章标题", "这是副标题", "这是标签", "# this is content", CURRENT_TIMESTAMP, "admin", "0", 0);
INSERT INTO ip_location(ip, location)
VALUES("127.0.0.1", "本机地址");
INSERT INTO comment(article_id, reply_name, content, comment_date, platform, browser, version, ip, location, reviewer_id, reviewer)
VALUES (0, "Admin", "这是评论", CURRENT_TIMESTAMP, "windows", "chrome", "94.0.4606.54", "127.0.0.0", "本机地址", 0, "admin");
INSERT INTO dictionary(key_, value_, key_type, value_type, create_date)
VALUES ("Shape", "形状", "en", "zh", CURRENT_TIMESTAMP);
INSERT INTO dictionary(key_, value_, key_type, value_type, create_date)
VALUES ("Six", "六", "en", "zh", CURRENT_TIMESTAMP);
INSERT INTO dictionary(key_, value_, key_type, value_type, create_date)
VALUES ("CNN", "卷积神经网络", "en", "zh", CURRENT_TIMESTAMP);

/*
USE
information_schema;
SELECT concat(round(sum(DATA_LENGTH / 1024 / 1024), 2), 'MB')
FROM TABLES;

SELECT concat(round(sum(DATA_LENGTH / 1024 / 1024), 2), 'MB') AS DATA
FROM TABLES
WHERE table_schema = 'blog';

SELECT concat(round(sum(DATA_LENGTH / 1024 / 1024), 2), 'M')
FROM TABLES
WHERE table_schema = 'blog'
  AND table_name = 'visitor';

show
variables like 'log_bin';
show
master status
show binlog events;
show full
processlist
show variables like '%char%';

ALTER TABLE article MODIFY COLUMN content MEDIUMTEXT;
ALTER TABLE article ADD COLUMN last_edit_date DATETIME;
ALTER TABLE article ADD COLUMN author_id MEDIUMINT NOT NULL;
ALTER TABLE article ADD COLUMN content_html MEDIUMTEXT;
ALTER TABLE article ADD COLUMN content_text MEDIUMTEXT;
ALTER TABLE article change content content_md MEDIUMTEXT;
ALTER TABLE article DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE article MODIFY content MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE article MODIFY content_html MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE article MODIFY content_text MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE article MODIFY title VARCHAR ( 100 ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE comment DEFAULT CHARACTER SET utf8mb4;
ALTER TABLE comment change replyName reply_name VARCHAR ( 50 );
ALTER TABLE comment MODIFY content MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE admin ADD COLUMN authority TINYINT;
ALTER TABLE admin ADD UNIQUE (username);
ALTER TABLE admin change username username VARCHAR ( 30 ) BINARY UNIQUE;

ALTER TABLE article change date create_date DATETIME;
ALTER TABLE article MODIFY author_id INT NOT NULL;
ALTER TABLE article MODIFY read_count INT;
ALTER TABLE article MODIFY title VARCHAR(100) NOT NULL;
ALTER TABLE article MODIFY tag VARCHAR(30) NOT NULL;
ALTER TABLE article MODIFY author VARCHAR(30) NOT NULL;
ALTER TABLE article MODIFY content_md MEDIUMTEXT NOT NULL;

ALTER TABLE visitor MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE visitor MODIFY location VARCHAR(50);
ALTER TABLE visitor change date visit_date DATETIME;
ALTER TABLE visitor MODIFY visit_date DATETIME NOT NULL;

ALTER TABLE upload MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE upload MODIFY location VARCHAR(50);
ALTER TABLE upload change date upload_date DATETIME;

ALTER TABLE download MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE download MODIFY location VARCHAR(50);
ALTER TABLE download change date download_date DATETIME;

ALTER TABLE ip_location MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE ip_location MODIFY location VARCHAR(50);

ALTER TABLE admin MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE admin MODIFY nickname VARCHAR(30);

ALTER TABLE comment MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE comment MODIFY article_id INT;
ALTER TABLE comment change date comment_date DATETIME;
ALTER TABLE comment change address location VARCHAR(30);
ALTER TABLE comment change osname platform VARCHAR(30);
ALTER TABLE comment MODIFY reply_name VARCHAR (30);
ALTER TABLE comment ADD COLUMN version VARCHAR(20);
ALTER TABLE comment ADD COLUMN reviewer_id INT;
ALTER TABLE comment ADD COLUMN reviewer VARCHAR(30);
ALTER TABLE comment DROP COLUMN reviewer

ALTER TABLE login MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE login change date login_date DATETIME;

ALTER TABLE search DEFAULT CHARACTER SET utf8mb4;
ALTER TABLE search MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE search change date search_date DATETIME;

ALTER TABLE dictionary MODIFY id INT NOT NULL AUTO_INCREMENT;
ALTER TABLE dictionary change date create_date DATETIME;

show create table article;
explain select * from article
explain select * from article where id = '21';
show variables like 'slow_query%';
show variables like 'long_query_time';
set global slow_query_log='ON';
set global long_query_time=2;
*/