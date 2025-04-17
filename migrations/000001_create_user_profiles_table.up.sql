-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for fixed-choice fields
CREATE TYPE profile_created_by AS ENUM (
    'Self', 'Brother', 'Sister', 'Parents', 'Friend', 'Relative'
);

CREATE TYPE community_type AS ENUM (
    'A muslim', 'Hanafi', 'Salafi', 'Sunni', 'Thableegh', 'Shia', 'Jamat Islami'
);

CREATE TYPE nationality_type AS ENUM (
    'India', 'UAE', 'UK', 'USA'
);

CREATE TYPE marital_status_type AS ENUM (
    'Never married', 'Widower', 'Divorced', 'Nikah Divorce'
);

-- Enum for Kerala districts
CREATE TYPE home_district_type AS ENUM (
    'Thiruvananthapuram', 'Kollam', 'Pathanamthitta', 'Alappuzha',
    'Kottayam', 'Idukki', 'Ernakulam', 'Thrissur', 'Palakkad',
    'Malappuram', 'Kozhikode', 'Wayanad', 'Kannur', 'Kasaragod'
);

-- Create user_profiles table without foreign key constraint
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    is_groom BOOLEAN NOT NULL,
    profile_created_by profile_created_by NOT NULL,
    name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    community community_type NOT NULL,
    nationality nationality_type NOT NULL,
    height DECIMAL(5,2) NOT NULL CHECK (height > 0),
    weight DECIMAL(5,2) NOT NULL CHECK (weight > 0),
    marital_status marital_status_type NOT NULL,
    is_physically_challenged BOOLEAN NOT NULL DEFAULT FALSE,
    home_district home_district_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Ensure one profile per user (without foreign key)
    CONSTRAINT unique_user_profile UNIQUE (user_id)
);

-- Create indexes for filtering and searching
CREATE INDEX idx_user_profiles_is_groom ON user_profiles(is_groom);
CREATE INDEX idx_user_profiles_community ON user_profiles(community);
CREATE INDEX idx_user_profiles_nationality ON user_profiles(nationality);
CREATE INDEX idx_user_profiles_marital_status ON user_profiles(marital_status);
CREATE INDEX idx_user_profiles_home_district ON user_profiles(home_district);