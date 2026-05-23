-- Auth closure: pending signup table (matches Java pendingSignup).

CREATE TABLE IF NOT EXISTS pendingsignup (
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(100) NOT NULL UNIQUE,
    signuptime  BIGINT,
    language    VARCHAR(10),
    token       VARCHAR(40)
);
