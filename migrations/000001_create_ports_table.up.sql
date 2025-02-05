CREATE TABLE IF NOT EXISTS ports (
    id          VARCHAR(50) PRIMARY KEY,
    name        VARCHAR(100),
    city        VARCHAR(100),
    country     VARCHAR(100),
    alias       TEXT[],
    regions     TEXT[],
    coordinates FLOAT[],
    province    VARCHAR(100),
    timezone    VARCHAR(50),
    unlocs      TEXT[],
    code        VARCHAR(10)
);