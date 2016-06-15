package main

import (
	"fmt"
	"log"
	"strings"

)

func (s *supervisor) authorize(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `verdict`, Super: &msg.Supervisor{}}

	switch svPermissionActionScopeMap[q.Super.PermAction] {
	case `global`:
		result.Super.Verdict, result.Super.VerdictAdmin = s.authorize_global(q)
	case `repository`:
	case `team`:
	default:
		goto unauthorized
	}
	q.Reply <- result
	return

unauthorized:
	result.Super.Verdict = 403
	q.Reply <- result
}

func (s *supervisor) authorize_global(q *msg.Request) (uint16, bool) {
	var (
		userUUID, permUUID string
		ok                 bool
	)
	// unknown user
	if userUUID, ok = s.id_user_rev.get(q.User); !ok {
		return 403, false
	}
	// user has omnipotence
	if s.global_permissions.assess(userUUID,
		`00000000-0000-0000-0000-000000000000`) {
		return 200, true
	}
	for _, perm := range svGlobalRequiredPermission[q.Super.PermAction] {
		// incomplete permission schema -> denied
		if permUUID, ok = s.id_permission.get(perm); !ok {
			return 403, false
		}
		if s.global_permissions.assess(userUUID, permUUID) {
			// has system permission (ordered + checked first)
			if strings.HasPrefix(perm, `system_`) {
				return 200, true
			}
			// has limited permission
			return 200, false
		}
	}
	return 403, false
}

func IsAuthorized(user, action, repository, monitoring, node string) (bool, bool) {
	returnChannel := make(chan msg.Result)
	// honour request for sandbox environment
	if SomaCfg.OpenInstance {
		return true, true
	}
	handler := handlerMap[`supervisor`].(supervisor)
	handler.input <- msg.Request{
		Type:   `supervisor`,
		Action: `authorize`,
		User:   user,
		Reply:  returnChannel,
		Super: &msg.Supervisor{
			Action:         `authorize`,
			PermAction:     action,
			PermRepository: repository,
			PermMonitoring: monitoring,
			PermNode:       node,
		},
	}
	result := <-returnChannel
	if result.Super.Verdict == 200 {
		if result.Super.VerdictAdmin {
			// authorized, admin access
			return true, true
		}
		// authorized, non-admin access
		return true, false
	}
	// not authorized
	log.Printf(LogStrErr, `supervisor`, `authorize`, result.Super.Verdict, fmt.Sprintf("Forbidden: %s, %s", user, action))
	return false, false
}

var svGlobalRequiredPermission = map[string][]string{
	`attributes_create`:        []string{`system_all`},
	`attributes_delete`:        []string{`system_all`},
	`attributes_list`:          []string{`system_all`, `global_schema`},
	`attributes_show`:          []string{`system_all`, `global_schema`},
	`category_create`:          []string{`system_all`},
	`category_delete`:          []string{`system_all`},
	`category_list`:            []string{`system_all`, `global_schema`},
	`category_show`:            []string{`system_all`, `global_schema`},
	`datacenters_create`:       []string{`system_all`},
	`datacenters_delete`:       []string{`system_all`},
	`datacenters_list`:         []string{`system_all`, `global_schema`},
	`datacenters_rename`:       []string{`system_all`},
	`datacenters_show`:         []string{`system_all`, `global_schema`},
	`datacenters_sync`:         []string{`system_all`},
	`environments_create`:      []string{`system_all`},
	`environments_delete`:      []string{`system_all`},
	`environments_list`:        []string{`system_all`, `global_schema`},
	`environments_rename`:      []string{`system_all`},
	`environments_show`:        []string{`system_all`, `global_schema`},
	`grant_global_right`:       []string{`system_all`},
	`grant_limited_right`:      []string{`system_all`},
	`grant_search`:             []string{`system_all`},
	`grant_system_right`:       []string{`system_all`},
	`levels_create`:            []string{`system_all`},
	`levels_delete`:            []string{`system_all`},
	`levels_list`:              []string{`system_all`, `global_schema`},
	`levels_search`:            []string{`system_all`, `global_schema`},
	`levels_show`:              []string{`system_all`, `global_schema`},
	`metrics_create`:           []string{`system_all`},
	`metrics_delete`:           []string{`system_all`},
	`metrics_list`:             []string{`system_all`, `global_schema`},
	`metrics_show`:             []string{`system_all`, `global_schema`},
	`modes_create`:             []string{`system_all`},
	`modes_delete`:             []string{`system_all`},
	`modes_list`:               []string{`system_all`, `global_schema`},
	`modes_show`:               []string{`system_all`, `global_schema`},
	`monitoring_create`:        []string{`system_all`},
	`monitoring_delete`:        []string{`system_all`},
	`node_create`:              []string{`system_all`},
	`node_delete`:              []string{`system_all`},
	`node_sync`:                []string{`system_all`},
	`node_update`:              []string{`system_all`},
	`oncall_create`:            []string{`system_all`},
	`oncall_delete`:            []string{`system_all`},
	`oncall_list`:              []string{`system_all`, `global_schema`},
	`oncall_search`:            []string{`system_all`, `global_schema`},
	`oncall_show`:              []string{`system_all`, `global_schema`},
	`oncall_update`:            []string{`system_all`},
	`permission_create`:        []string{`system_all`},
	`permission_delete`:        []string{`system_all`},
	`permission_list`:          []string{`system_all`, `global_schema`},
	`permission_search`:        []string{`system_all`, `global_schema`},
	`permission_show`:          []string{`system_all`, `global_schema`},
	`predicates_create`:        []string{`system_all`},
	`predicates_delete`:        []string{`system_all`},
	`predicates_list`:          []string{`system_all`, `global_schema`},
	`predicates_show`:          []string{`system_all`, `global_schema`},
	`property_native_create`:   []string{`system_all`},
	`property_native_delete`:   []string{`system_all`},
	`property_native_list`:     []string{`system_all`, `global_schema`},
	`property_native_show`:     []string{`system_all`, `global_schema`},
	`property_system_create`:   []string{`system_all`},
	`property_system_delete`:   []string{`system_all`},
	`property_system_list`:     []string{`system_all`, `global_schema`},
	`property_system_search`:   []string{`system_all`, `global_schema`},
	`property_system_show`:     []string{`system_all`, `global_schema`},
	`property_template_create`: []string{`system_all`},
	`property_template_delete`: []string{`system_all`},
	`property_template_list`:   []string{`system_all`, `global_schema`},
	`property_template_search`: []string{`system_all`, `global_schema`},
	`property_template_show`:   []string{`system_all`, `global_schema`},
	`providers_create`:         []string{`system_all`},
	`providers_delete`:         []string{`system_all`},
	`providers_list`:           []string{`system_all`, `global_schema`},
	`providers_show`:           []string{`system_all`, `global_schema`},
	`repository_create`:        []string{`system_all`, `global_schema`},
	`repository_delete`:        []string{`system_all`},
	`revoke_global_right`:      []string{`system_all`},
	`revoke_limited_right`:     []string{`system_all`},
	`revoke_system_right`:      []string{`system_all`},
	`servers_create`:           []string{`system_all`},
	`servers_delete`:           []string{`system_all`},
	`servers_list`:             []string{`system_all`, `global_schema`},
	`servers_search`:           []string{`system_all`, `global_schema`},
	`servers_show`:             []string{`system_all`, `global_schema`},
	`servers_sync`:             []string{`system_all`},
	`servers_update`:           []string{`system_all`},
	`states_create`:            []string{`system_all`},
	`states_delete`:            []string{`system_all`},
	`states_list`:              []string{`system_all`, `global_schema`},
	`states_rename`:            []string{`system_all`, `global_schema`},
	`states_show`:              []string{`system_all`, `global_schema`},
	`status_create`:            []string{`system_all`},
	`status_delete`:            []string{`system_all`},
	`status_list`:              []string{`system_all`, `global_schema`},
	`status_show`:              []string{`system_all`, `global_schema`},
	`team_create`:              []string{`system_all`},
	`team_delete`:              []string{`system_all`},
	`team_list`:                []string{`system_all`, `global_schema`},
	`team_search`:              []string{`system_all`, `global_schema`},
	`team_show`:                []string{`system_all`, `global_schema`},
	`team_sync`:                []string{`system_all`},
	`team_update`:              []string{`system_all`},
	`types_create`:             []string{`system_all`},
	`types_delete`:             []string{`system_all`},
	`types_list`:               []string{`system_all`, `global_schema`},
	`types_rename`:             []string{`system_all`, `global_schema`},
	`types_show`:               []string{`system_all`, `global_schema`},
	`units_create`:             []string{`system_all`},
	`units_delete`:             []string{`system_all`},
	`units_list`:               []string{`system_all`, `global_schema`},
	`units_show`:               []string{`system_all`, `global_schema`},
	`users_create`:             []string{`system_all`},
	`users_delete`:             []string{`system_all`},
	`users_list`:               []string{`system_all`, `global_schema`},
	`users_search`:             []string{`system_all`, `global_schema`},
	`users_show`:               []string{`system_all`, `global_schema`},
	`users_sync`:               []string{`system_all`},
	`users_update`:             []string{`system_all`},
	`validity_create`:          []string{`system_all`},
	`validity_delete`:          []string{`system_all`},
	`validity_list`:            []string{`system_all`, `global_schema`},
	`validity_show`:            []string{`system_all`, `global_schema`},
	`view_create`:              []string{`system_all`},
	`view_delete`:              []string{`system_all`},
	`view_list`:                []string{`system_all`, `global_schema`},
	`view_rename`:              []string{`system_all`},
	`view_show`:                []string{`system_all`, `global_schema`},
}

var svPermissionActionScopeMap = map[string]string{
	`attributes_create`:              `global`,
	`attributes_delete`:              `global`,
	`attributes_list`:                `global`,
	`attributes_show`:                `global`,
	`category_create`:                `global`,
	`category_delete`:                `global`,
	`category_list`:                  `global`,
	`category_show`:                  `global`,
	`datacenters_create`:             `global`,
	`datacenters_delete`:             `global`,
	`datacenters_list`:               `global`,
	`datacenters_rename`:             `global`,
	`datacenters_show`:               `global`,
	`datacenters_sync`:               `global`,
	`environments_create`:            `global`,
	`environments_delete`:            `global`,
	`environments_list`:              `global`,
	`environments_rename`:            `global`,
	`environments_show`:              `global`,
	`grant_global_right`:             `global`,
	`grant_limited_right`:            `global`,
	`grant_search`:                   `global`,
	`grant_system_right`:             `global`,
	`levels_create`:                  `global`,
	`levels_delete`:                  `global`,
	`levels_list`:                    `global`,
	`levels_search`:                  `global`,
	`levels_show`:                    `global`,
	`metrics_create`:                 `global`,
	`metrics_delete`:                 `global`,
	`metrics_list`:                   `global`,
	`metrics_show`:                   `global`,
	`modes_create`:                   `global`,
	`modes_delete`:                   `global`,
	`modes_list`:                     `global`,
	`modes_show`:                     `global`,
	`monitoring_create`:              `global`,
	`monitoring_delete`:              `global`,
	`node_create`:                    `global`,
	`node_delete`:                    `global`,
	`node_sync`:                      `global`,
	`node_update`:                    `global`,
	`oncall_create`:                  `global`,
	`oncall_delete`:                  `global`,
	`oncall_list`:                    `global`,
	`oncall_search`:                  `global`,
	`oncall_show`:                    `global`,
	`oncall_update`:                  `global`,
	`permission_create`:              `global`,
	`permission_delete`:              `global`,
	`permission_list`:                `global`,
	`permission_search`:              `global`,
	`permission_show`:                `global`,
	`predicates_create`:              `global`,
	`predicates_delete`:              `global`,
	`predicates_list`:                `global`,
	`predicates_show`:                `global`,
	`property_native_create`:         `global`,
	`property_native_delete`:         `global`,
	`property_native_list`:           `global`,
	`property_native_show`:           `global`,
	`property_service_global_create`: `global`,
	`property_service_global_delete`: `global`,
	`property_service_global_list`:   `global`,
	`property_service_global_search`: `global`,
	`property_service_global_show`:   `global`,
	`property_system_create`:         `global`,
	`property_system_delete`:         `global`,
	`property_system_list`:           `global`,
	`property_system_search`:         `global`,
	`property_system_show`:           `global`,
	`providers_create`:               `global`,
	`providers_delete`:               `global`,
	`providers_list`:                 `global`,
	`providers_show`:                 `global`,
	`repository_create`:              `global`,
	`repository_delete`:              `global`,
	`revoke_global_right`:            `global`,
	`revoke_limited_right`:           `global`,
	`revoke_system_right`:            `global`,
	`servers_create`:                 `global`,
	`servers_delete`:                 `global`,
	`servers_list`:                   `global`,
	`servers_search`:                 `global`,
	`servers_show`:                   `global`,
	`servers_sync`:                   `global`,
	`servers_update`:                 `global`,
	`states_create`:                  `global`,
	`states_delete`:                  `global`,
	`states_list`:                    `global`,
	`states_rename`:                  `global`,
	`states_show`:                    `global`,
	`status_create`:                  `global`,
	`status_delete`:                  `global`,
	`status_list`:                    `global`,
	`status_show`:                    `global`,
	`team_create`:                    `global`,
	`team_delete`:                    `global`,
	`team_list`:                      `global`,
	`team_search`:                    `global`,
	`team_show`:                      `global`,
	`team_sync`:                      `global`,
	`team_update`:                    `global`,
	`types_create`:                   `global`,
	`types_delete`:                   `global`,
	`types_list`:                     `global`,
	`types_rename`:                   `global`,
	`types_show`:                     `global`,
	`units_create`:                   `global`,
	`units_delete`:                   `global`,
	`units_list`:                     `global`,
	`units_show`:                     `global`,
	`users_create`:                   `global`,
	`users_delete`:                   `global`,
	`users_list`:                     `global`,
	`users_search`:                   `global`,
	`users_show`:                     `global`,
	`users_sync`:                     `global`,
	`users_update`:                   `global`,
	`validity_create`:                `global`,
	`validity_delete`:                `global`,
	`validity_list`:                  `global`,
	`validity_show`:                  `global`,
	`view_create`:                    `global`,
	`view_delete`:                    `global`,
	`view_list`:                      `global`,
	`view_rename`:                    `global`,
	`view_show`:                      `global`,
	`buckets_create`:                 `repository`,
	`buckets_property_add`:           `repository`,
	`buckets_list`:                   `repository`,
	`buckets_search`:                 `repository`,
	`buckets_show`:                   `repository`,
	`checks_create`:                  `repository`,
	`checks_list`:                    `repository`,
	`checks_search`:                  `repository`,
	`checks_show`:                    `repository`,
	`clusters_create`:                `repository`,
	`clusters_member_add`:            `repository`,
	`clusters_property_add`:          `repository`,
	`clusters_list`:                  `repository`,
	`clusters_members_list`:          `repository`,
	`clusters_search`:                `repository`,
	`clusters_show`:                  `repository`,
	`groups_create`:                  `repository`,
	`groups_member_add`:              `repository`,
	`groups_property_add`:            `repository`,
	`groups_list`:                    `repository`,
	`groups_members_list`:            `repository`,
	`groups_search`:                  `repository`,
	`groups_show`:                    `repository`,
	`node_property_add`:              `repository`,
	`property_custom_create`:         `repository`,
	`property_custom_delete`:         `repository`,
	`property_custom_list`:           `repository`,
	`property_custom_search`:         `repository`,
	`property_custom_show`:           `repository`,
	`repository_list`:                `repository`,
	`repository_search`:              `repository`,
	`repository_show`:                `repository`,
	`node_list`:                      `team`,
	`node_search`:                    `team`,
	`node_show`:                      `team`,
	`node_show_config`:               `team`,
	`property_service_team_list`:     `team`,
	`property_service_team_search`:   `team`,
	`property_service_team_show`:     `team`,
	`property_service_team_create`:   `team`,
	`property_service_team_delete`:   `team`,
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
