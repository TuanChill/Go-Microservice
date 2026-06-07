DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'auth_service_user') THEN
        CREATE ROLE auth_service_user LOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'auth_schema_owner') THEN
        CREATE ROLE auth_schema_owner NOLOGIN;
    END IF;
END $$;

ALTER SCHEMA auth OWNER TO auth_schema_owner;
GRANT USAGE ON SCHEMA auth TO auth_service_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA auth TO auth_service_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA auth TO auth_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE auth_schema_owner IN SCHEMA auth GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO auth_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE auth_schema_owner IN SCHEMA auth GRANT USAGE, SELECT ON SEQUENCES TO auth_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE auth_service_user IN SCHEMA auth GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO auth_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE auth_service_user IN SCHEMA auth GRANT USAGE, SELECT ON SEQUENCES TO auth_service_user;

REVOKE ALL ON SCHEMA user_profiles FROM auth_service_user;
REVOKE ALL ON SCHEMA notification FROM auth_service_user;
