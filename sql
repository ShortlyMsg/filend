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
    file_name TEXT,
    file_hash TEXT,
    file_size BIGINT,
    isUploaded BOOLEAN DEFAULT FALSE,
    file_model_id VARCHAR(36) NOT NULL REFERENCES file_models(file_model_id) ON DELETE CASCADE
);