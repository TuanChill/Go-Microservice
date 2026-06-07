DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'user_service_user') THEN
        CREATE ROLE user_service_user LOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'user_profiles_schema_owner') THEN
        CREATE ROLE user_profiles_schema_owner NOLOGIN;
    END IF;
END $$;

ALTER SCHEMA user_profiles OWNER TO user_profiles_schema_owner;
GRANT USAGE ON SCHEMA user_profiles TO user_service_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA user_profiles TO user_service_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA user_profiles TO user_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE user_profiles_schema_owner IN SCHEMA user_profiles GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO user_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE user_profiles_schema_owner IN SCHEMA user_profiles GRANT USAGE, SELECT ON SEQUENCES TO user_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE user_service_user IN SCHEMA user_profiles GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO user_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE user_service_user IN SCHEMA user_profiles GRANT USAGE, SELECT ON SEQUENCES TO user_service_user;

REVOKE ALL ON SCHEMA auth FROM user_service_user;
REVOKE ALL ON SCHEMA notification FROM user_service_user;
