PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE USR(
ID INTEGER PRIMARY KEY  AUTOINCREMENT,
NAME TEXT NOT NULL,
PASSWORD TEXT NOT NULL
);
INSERT INTO USR VALUES(1,'admin','adminPeople');
INSERT INTO USR VALUES(2,'evpeople','evpeople');
INSERT INTO USR VALUES(3,'verso','verso');
INSERT INTO USR VALUES(4,'username','password');
COMMIT;
