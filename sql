CREATE TABLE file_models (
    id TEXT PRIMARY KEY,
    otp TEXT NOT NULL,
    user_security_code TEXT NOT NULL,
    file_names TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);