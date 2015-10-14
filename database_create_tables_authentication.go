package main

func createTablesAuthentication(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 25)

	queryMap["createTableUserAuthentication"] = `
create table if not exists auth.user_authentication (
  user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
  algorithm                   varchar(16)     NOT NULL,
  rounds                      numeric(10,0)   NOT NULL,
  salt                        text            NOT NULL,
  crypt                       text            NOT NULL,
  reset_pending               boolean         NOT NULL DEFAULT 'no',
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableUserAuthentication"
	idx++

	queryMap["createTableUserTokenAuthentication"] = `
create table if not exists auth.user_token_authentication (
  user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
  token                       varchar(256)    UNIQUE NOT NULL,
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableUserTokenAuthentication"
	idx++

	queryMap["createTableUserKeys"] = `
create table if not exists auth.user_keys (
  user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
  user_key_fingerprint        varchar(128)    NOT NULL,
  user_key_public             text            NOT NULL,
  user_key_active             boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableUserKeys"
	idx++

	queryMap["createIndexUniqueActiveUserKey"] = `
create index _unique_active_user_key
  on auth.user_keys ( user_id, user_key_active )
  where user_key_active`
	queries[idx] = "createIndexUniqueActiveUserKey"
	idx++

	queryMap["createTableUserClientCertificates"] = `
create table if not exists auth.user_client_certificates (
  user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
  user_cert_fingerprint       varchar(128)    NOT NULL,
  user_cert_active            boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableUserClientCertificates"
	idx++

	queryMap["createIndexUniqueActiveUserCert"] = `
create index _unique_active_user_cert
  on auth.user_client_certificates ( user_id, user_cert_active )
  where user_cert_active`
	queries[idx] = "createIndexUniqueActiveUserCert"
	idx++

	queryMap["createTableAdmins"] = `
create table if not exists auth.admins (
  admin_id                    uuid            PRIMARY KEY,
  admin_uid                   varchar(256)    UNIQUE NOT NULL,
  admin_user_uid              varchar(256)    NOT NULL REFERENCES inventory.users ( user_uid ) ON DELETE CASCADE,
  admin_is_active             boolean         NOT NULL DEFAULT 'yes',
  CHECK( left( admin_uid, 6 ) = 'admin_' ),
  CHECK( position( admin_user_uid in admin_uid ) != 0 )
);`
	queries[idx] = "createTableAdmins"
	idx++

	queryMap["createTableAdminAuthentication"] = `
create table if not exists auth.admin_authentication (
  admin_id                    uuid            NOT NULL REFERENCES auth.admins ( admin_id ) ON DELETE CASCADE,
  algorithm                   varchar(16)     NOT NULL,
  rounds                      numeric(10,0)   NOT NULL,
  salt                        text            NOT NULL,
  crypt                       text            NOT NULL,
  reset_pending               boolean         NOT NULL DEFAULT 'no',
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableAdminAuthentication"
	idx++

	queryMap["createTableAdminTokenAuthentication"] = `
create table if not exists auth.admin_token_authentication (
  admin_id                    uuid            NOT NULL REFERENCES auth.admins ( admin_id ) ON DELETE CASCADE,
  token                       varchar(256)    UNIQUE NOT NULL,
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableAdminTokenAuthentication"
	idx++

	queryMap["createTableAdminKeys"] = `
create table if not exists auth.admin_keys (
  admin_id                    uuid            NOT NULL REFERENCES auth.admins ( admin_id ) ON DELETE CASCADE,
  admin_key_fingerprint       varchar(128)    NOT NULL,
  admin_key_public            text            NOT NULL,
  admin_key_active            boolean         NOT NULL DEFAULT 'yes' 
);`
	queries[idx] = "createTableAdminKeys"
	idx++

	queryMap["createIndexUniqueActiveAdminKey"] = `
create index _unique_active_admin_key
  on auth.admin_keys ( admin_id, admin_key_active )
  where admin_key_active`
	queries[idx] = "createIndexUniqueActiveAdminKey"
	idx++

	queryMap["createTableAdminClientCertificates"] = `
create table if not exists auth.admin_client_certificates (
  admin_id                    uuid            NOT NULL REFERENCES auth.admins ( admin_id ) ON DELETE CASCADE,
  admin_cert_fingerprint      varchar(128)    NOT NULL,
  admin_cert_active           boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableAdminClientCertificates"
	idx++

	queryMap["createIndexUniqueActiveAdminCert"] = `
create index _unique_active_admin_cert
  on auth.admin_client_certificates ( admin_id, admin_cert_active )
  where admin_cert_active`
	queries[idx] = "createIndexUniqueActiveAdminCert"
	idx++

	queryMap["createTableTools"] = `
create table if not exists auth.tools (
  tool_id                     uuid            PRIMARY KEY,
  tool_name                   varchar(256)    UNIQUE NOT NULL,
  tool_owner_id               uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE RESTRICT,
  created                     timestamptz(3)  NOT NULL DEFAULT NOW(),
  CHECK( EXTRACT( TIMEZONE FROM created ) = '0' )
);`
	queries[idx] = "createTableTools"
	idx++

	queryMap["createTableToolAuthentication"] = `
create table if not exists auth.tool_authentication (
  tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE,
  algorithm                   varchar(16)     NOT NULL,
  rounds                      numeric(10,0)   NOT NULL,
  salt                        text            NOT NULL,
  crypt                       text            NOT NULL,
  reset_pending               boolean         NOT NULL DEFAULT 'no',
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableToolAuthentication"
	idx++

	queryMap["createTableToolTokenAuthentication"] = `
create table if not exists auth.tool_token_authentication (
  tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE,
  token                       varchar(256)    UNIQUE NOT NULL,
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableToolTokenAuthentication"
	idx++

	queryMap["createTableToolKeys"] = `
create table if not exists auth.tool_keys (
  tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE,
  tool_key_fingerprint        varchar(128)    NOT NULL,
  tool_key_public             text            NOT NULL,
  tool_key_active             boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableToolKeys"
	idx++

	queryMap["createIndexUniqueActiveToolKey"] = `
create index _unique_active_tool_key
  on auth.tool_keys ( tool_id, tool_key_active )
  where tool_key_active`
	queries[idx] = "createIndexUniqueActiveToolKey"
	idx++

	queryMap["createTableToolClientCertificates"] = `
create table if not exists auth.tool_client_certificates (
  tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE,
  tool_cert_fingerprint       varchar(128)    NOT NULL,
  tool_cert_active            boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableToolClientCertificates"
	idx++

	queryMap["createIndexUniqueActiveToolCert"] = `
create index _unique_active_tool_cert
  on auth.tool_client_certificates ( tool_id, tool_cert_active )
  where tool_cert_active`
	queries[idx] = "createIndexUniqueActiveToolCert"
	idx++

	queryMap["createTablePasswordReset"] = `
create table if not exists auth.password_reset (
  user_id                     uuid            NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
  admin_id                    uuid            NULL REFERENCES auth.admins ( admin_id ) ON DELETE CASCADE,
  tool_id                     uuid            NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE,
  token                       varchar(256)    UNIQUE NOT NULL,
  valid_from                  timestamptz(3)  NOT NULL,
  valid_until                 timestamptz(3)  NOT NULL,
  token_used                  boolean         NOT NULL DEFAULT 'no',
  token_invalidated           boolean         NOT NULL DEFAULT 'no',
  CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
  CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' ),
  CHECK(    ( user_id IS NOT NULL AND admin_id IS     NULL AND tool_id IS     NULL )
         OR ( user_id IS     NULL AND admin_id IS NOT NULL AND tool_id IS     NULL )
         OR ( user_id IS     NULL AND admin_id IS     NULL AND tool_id IS NOT NULL ) )
);`
	queries[idx] = "createTablePasswordReset"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
