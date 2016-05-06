package main

const MaxInt = int(^uint(0) >> 1)

var UpgradeVersions = map[string]map[int]func(int) int{
	"inventory": map[int]func(int) int{
		201605060001: upgrade_inventory_201605060001,
	},
	"auth": map[int]func(int) int{
		201605060001: upgrade_auth_201605060001,
	},
	"soma": map[int]func(int) int{
		201605060001: upgrade_soma_201605060001,
		201605060002: upgrade_soma_201605060002,
	},
}

func UpgradeSchema(target int) error {
	// no specific target specified => upgrade all the way
	if target == 0 {
		target = MaxInt
	}

loop:
	for schema, _ := range UpgradeVersions {
		// fetch current version from database
		version := getCurrentSchemaVersion(schema)

		if version >= target {
			// schema is already as updated as we need
			continue loop
		}

		for f, ok := UpgradeVersions[schema][version]; ok; f, ok = UpgradeVersions[schema][version] {
			version = f(version)
			if version == 0 {
				// something broke
				// TODO: set error
				break loop
			} else if version >= target {
				// job done, continue with next schema
				continue loop
			}
		}
	}
	return nil
}

func upgrade_inventory_201605060001(curr int) int {
	if curr != 201605060001 {
		return 0
	}
	return 201605060002
}

func upgrade_auth_201605060001(curr int) int {
	if curr != 201605060001 {
		return 0
	}
	return 201605060002
}

func upgrade_soma_201605060001(curr int) int {
	if curr != 201605060001 {
		return 0
	}
	return 201605060002
}

func upgrade_soma_201605060002(curr int) int {
	if curr != 201605060002 {
		return 0
	}
	return 201605060003
}

func getCurrentSchemaVersion(schema string) int {
	// TODO: needs hook to mock current version to report for no-execute
	//       case
	return 201605060001
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
