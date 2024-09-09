CREATE TABLE file_models (
    id VARCHAR(36) PRIMARY KEY,
    otp VARCHAR(6),
    user_security_code VARCHAR(4),
    file_names TEXT[],
    file_hashes TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);