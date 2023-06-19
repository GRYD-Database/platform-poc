CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS storage (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    wallet TEXT NOT NULL,
    txHash TEXT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

---- create above / drop below ----
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE IF EXISTS storage;
