-- Migration to create users table
-- Save as: migrations/001_create_users_table.sql

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    -- Premium subscription fields
    is_premium_user BOOLEAN DEFAULT FALSE,
    subscription_tier VARCHAR(20) CHECK (subscription_tier IN ('free', 'pro', 'enterprise')) DEFAULT 'free',
    subscription_start_date TIMESTAMP,
    subscription_end_date TIMESTAMP,

    -- notification settings
    notification_preferences JSONB DEFAULT '{"email": true, "push": true, "sms": false}',
    timezone VARCHAR(50) DEFAULT 'UTC',

    -- metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_subscription_tier ON users(subscription_tier);
CREATE INDEX idx_users_is_premium ON users(is_premium_user);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Update the notifications table to reference users
ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id);

-- Index for notifications-user relationship
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);