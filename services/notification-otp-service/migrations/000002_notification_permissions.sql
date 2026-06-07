DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'notification_service_user') THEN
        CREATE ROLE notification_service_user LOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'notification_schema_owner') THEN
        CREATE ROLE notification_schema_owner NOLOGIN;
    END IF;
END $$;

ALTER SCHEMA notification OWNER TO notification_schema_owner;
GRANT USAGE ON SCHEMA notification TO notification_service_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA notification TO notification_service_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA notification TO notification_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE notification_schema_owner IN SCHEMA notification GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO notification_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE notification_schema_owner IN SCHEMA notification GRANT USAGE, SELECT ON SEQUENCES TO notification_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE notification_service_user IN SCHEMA notification GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO notification_service_user;
ALTER DEFAULT PRIVILEGES FOR ROLE notification_service_user IN SCHEMA notification GRANT USAGE, SELECT ON SEQUENCES TO notification_service_user;

REVOKE ALL ON SCHEMA auth FROM notification_service_user;
REVOKE ALL ON SCHEMA user_profiles FROM notification_service_user;
