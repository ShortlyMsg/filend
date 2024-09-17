CREATE TABLE file_models (
    file_model_id VARCHAR(36) PRIMARY KEY,
    otp VARCHAR(6),
    user_security_code VARCHAR(4),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE file_details (
    file_details_id VARCHAR(36) PRIMARY KEY,
    file_names TEXT[],
    file_hashes TEXT[],
    file_model_id NOT NULL REFERENCES file_models(file_model_id) ON DELETE CASCADE
);