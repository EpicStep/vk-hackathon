BEGIN;

CREATE TABLE images (
    id VARCHAR(36) PRIMARY KEY,
    image LONGBLOB,
    hash BIGINT UNSIGNED,
    Height INT UNSIGNED,
    Width INT UNSIGNED
);

COMMIT;