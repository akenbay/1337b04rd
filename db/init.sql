-- init.sql
-- Database initialization for 1337b04rd (anonymous imageboard)

BEGIN;

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users/Sessions table (anonymous users identified by cookies)
CREATE TABLE IF NOT EXISTS user_sessions (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    avatar_url TEXT NOT NULL,               -- URL from Rick and Morty API
    username TEXT NOT NULL,           -- Character name from API
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + INTERVAL '7 days'
);

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    post_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES user_sessions(session_id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_key TEXT,                         -- S3 object key for the image
    bucket_name TEXT,                       -- S3 bucket name
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,    -- When post was moved to archive
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(comment_id) ON DELETE CASCADE, -- For nested comments
    session_id UUID REFERENCES user_sessions(session_id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at);
CREATE INDEX IF NOT EXISTS idx_posts_archived ON posts(is_archived, archived_at);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at);

-- Function to update timestamp on post update
CREATE OR REPLACE FUNCTION update_post_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_post_timestamp
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_post_timestamp();

-- Function to archive old posts
CREATE OR REPLACE FUNCTION archive_old_posts()
RETURNS VOID AS $$
BEGIN
    -- Archive posts without comments after 10 minutes
    UPDATE posts
    SET is_archived = TRUE, archived_at = NOW()
    WHERE is_archived = FALSE
    AND created_at < NOW() - INTERVAL '10 minutes'
    AND NOT EXISTS (
        SELECT 1 FROM comments 
        WHERE comments.post_id = posts.post_id
    );
    
    -- Archive posts with comments after 15 minutes of inactivity
    UPDATE posts
    SET is_archived = TRUE, archived_at = NOW()
    WHERE is_archived = FALSE
    AND (
        SELECT MAX(created_at) FROM comments 
        WHERE comments.post_id = posts.post_id
    ) < NOW() - INTERVAL '15 minutes';
END;
$$ LANGUAGE plpgsql;

COMMIT;