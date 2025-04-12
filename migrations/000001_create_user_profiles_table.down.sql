-- Drop indexes (PostgreSQL automatically drops indexes with the table,
-- but explicitly dropping them can be helpful for clarity)
DROP INDEX IF EXISTS idx_user_profiles_home_district;
DROP INDEX IF EXISTS idx_user_profiles_marital_status;
DROP INDEX IF EXISTS idx_user_profiles_nationality;
DROP INDEX IF EXISTS idx_user_profiles_community;
DROP INDEX IF EXISTS idx_user_profiles_is_groom;

-- Drop the user_profiles table
DROP TABLE IF EXISTS user_profiles;

-- Drop the enum types created by this migration
DROP TYPE IF EXISTS home_district_type;
DROP TYPE IF EXISTS marital_status_type;
DROP TYPE IF EXISTS nationality_type;
DROP TYPE IF EXISTS community_type;
DROP TYPE IF EXISTS profile_created_by;

-- Note: The "uuid-ossp" extension is left installed.
