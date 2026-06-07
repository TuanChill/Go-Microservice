package models

import (
	"os"
	"strings"
	"testing"
)

func TestServicePermissionMigrationsDeclareIsolation(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		schema      string
		serviceUser string
		ownerRole   string
		revokes     []string
	}{
		{
			name:        "auth service",
			path:        "../../../auth-service/migrations/000002_auth_permissions.sql",
			schema:      "auth",
			serviceUser: "auth_service_user",
			ownerRole:   "auth_schema_owner",
			revokes: []string{
				"REVOKE ALL ON SCHEMA user_profiles FROM auth_service_user",
				"REVOKE ALL ON SCHEMA notification FROM auth_service_user",
			},
		},
		{
			name:        "user service",
			path:        "../../../user-service/migrations/000002_user_profiles_permissions.sql",
			schema:      "user_profiles",
			serviceUser: "user_service_user",
			ownerRole:   "user_profiles_schema_owner",
			revokes: []string{
				"REVOKE ALL ON SCHEMA auth FROM user_service_user",
				"REVOKE ALL ON SCHEMA notification FROM user_service_user",
			},
		},
		{
			name:        "notification service",
			path:        "../../../notification-otp-service/migrations/000002_notification_permissions.sql",
			schema:      "notification",
			serviceUser: "notification_service_user",
			ownerRole:   "notification_schema_owner",
			revokes: []string{
				"REVOKE ALL ON SCHEMA auth FROM notification_service_user",
				"REVOKE ALL ON SCHEMA user_profiles FROM notification_service_user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}
			migration := string(content)
			required := []string{
				"CREATE ROLE " + tt.serviceUser + " LOGIN",
				"CREATE ROLE " + tt.ownerRole + " NOLOGIN",
				"ALTER SCHEMA " + tt.schema + " OWNER TO " + tt.ownerRole,
				"GRANT USAGE ON SCHEMA " + tt.schema + " TO " + tt.serviceUser,
				"GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA " + tt.schema + " TO " + tt.serviceUser,
				"GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA " + tt.schema + " TO " + tt.serviceUser,
				"ALTER DEFAULT PRIVILEGES FOR ROLE " + tt.ownerRole + " IN SCHEMA " + tt.schema,
				"ALTER DEFAULT PRIVILEGES FOR ROLE " + tt.serviceUser + " IN SCHEMA " + tt.schema,
			}
			for _, expected := range required {
				if !strings.Contains(migration, expected) {
					t.Fatalf("missing %q", expected)
				}
			}
			for _, revoke := range tt.revokes {
				if !strings.Contains(migration, revoke) {
					t.Fatalf("missing revoke %q", revoke)
				}
			}
			if strings.Contains(migration, "GRANT "+tt.ownerRole+" TO "+tt.serviceUser) {
				t.Fatal("runtime service user must not be granted schema-owner role")
			}
			if strings.Contains(migration, "GRANT ALL ON SCHEMA") {
				t.Fatal("migration must not grant all privileges on schema")
			}
		})
	}
}
