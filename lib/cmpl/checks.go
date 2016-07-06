package cmpl

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// I'm sorry as well.
func CheckAdd(c *cli.Context) {
	topArgs := []string{`in`, `on`, `with`, `interval`, `inheritance`, `childrenonly`, `extern`, `threshold`, `constraint`}
	thrArgs := []string{`predicate`, `level`, `value`}
	ctrArgs := []string{`service`, `oncall`, `attribute`, `system`, `native`, `custom`}
	onArgs := []string{`repository`, `bucket`, `group`, `cluster`, `node`}

	if c.NArg() == 0 {
		return
	}

	if c.NArg() == 1 {
		for _, t := range topArgs {
			fmt.Println(t)
		}
	}

	skipNext := 0
	subON := false
	subTHRESHOLD := false
	subCONSTRAINT := false

	hasIN := false
	hasON := false
	hasWITH := false
	hasINTERVAL := false
	hasINHERITANCE := false
	hasCHILDRENONLY := false
	hasEXTERN := false

	hasTHR_predicate := false
	hasTHR_level := false
	hasTHR_value := false

	hasCTR_service := false
	hasCTR_oncall := false
	hasCTR_attribute := false
	hasCTR_system := false
	hasCTR_native := false
	hasCTR_custom := false
	hasCTR_selected_service := false
	hasCTR_selected_oncall := false

	for _, t := range c.Args().Tail() {
		if skipNext > 0 {
			skipNext--
			continue
		}
		if subON {
			skipNext = 1
			subON = false
		}
		if subTHRESHOLD {
			if hasTHR_predicate && hasTHR_level && hasTHR_value {
				subTHRESHOLD = false
				hasTHR_predicate = false
				hasTHR_level = false
				hasTHR_value = false
			} else {
				switch t {
				case `predicate`:
					skipNext = 1
					hasTHR_predicate = true
					continue
				case `level`:
					skipNext = 1
					hasTHR_level = true
					continue
				case `value`:
					skipNext = 1
					hasTHR_value = true
					continue
				}
			}
		}
		if subCONSTRAINT {
			if hasCTR_selected_service {
				skipNext = 1
				hasCTR_selected_service = false
				continue
			}
			if hasCTR_selected_oncall {
				skipNext = 1
				hasCTR_selected_oncall = false
				continue
			}
			if hasCTR_service || hasCTR_oncall || hasCTR_attribute || hasCTR_system || hasCTR_native || hasCTR_custom {
				subCONSTRAINT = false
				hasCTR_service = false
				hasCTR_oncall = false
				hasCTR_attribute = false
				hasCTR_system = false
				hasCTR_native = false
				hasCTR_custom = false
				hasCTR_selected_service = false
				hasCTR_selected_oncall = false
			} else {
				switch t {
				case `service`:
					hasCTR_selected_service = true
					hasCTR_service = true
					continue
				case `oncall`:
					hasCTR_selected_oncall = true
					hasCTR_oncall = true
					continue
				case `attribute`:
					skipNext = 2
					hasCTR_attribute = true
					continue
				case `system`:
					skipNext = 2
					hasCTR_system = true
					continue
				case `native`:
					skipNext = 2
					hasCTR_native = true
					continue
				case `custom`:
					skipNext = 2
					hasCTR_custom = true
					continue
				}
			}

		}
		switch t {
		case `in`:
			skipNext = 1
			hasIN = true
			continue
		case `on`:
			hasON = true
			subON = true
			continue
		case `with`:
			skipNext = 1
			hasWITH = true
			continue
		case `interval`:
			skipNext = 1
			hasINTERVAL = true
			continue
		case `inheritance`:
			skipNext = 1
			hasINHERITANCE = true
			continue
		case `childrenonly`:
			skipNext = 1
			hasCHILDRENONLY = true
			continue
		case `extern`:
			skipNext = 1
			hasEXTERN = true
			continue
		case `threshold`:
			subTHRESHOLD = true
			continue
		case `constraint`:
			subCONSTRAINT = true
			continue
		}
	}
	// skipNext not yet consumed
	if skipNext > 0 {
		return
	}
	// in subchain: ON
	if subON {
		for _, t := range onArgs {
			fmt.Println(t)
		}
		return
	}
	// in subchain: CONSTRAINT
	if subCONSTRAINT {
		if hasCTR_selected_service || hasCTR_selected_oncall {
			fmt.Println(`name`)
			return
		}
		if !(hasCTR_service || hasCTR_oncall || hasCTR_attribute || hasCTR_system || hasCTR_native || hasCTR_custom) {
			for _, t := range ctrArgs {
				fmt.Println(t)
			}
			return
		}
	}
	// in subchain: THRESHOLD
	if subTHRESHOLD {
		if !(hasTHR_predicate && hasTHR_level && hasTHR_value) {
			for _, t := range thrArgs {
				switch t {
				case `predicate`:
					if !hasTHR_predicate {
						fmt.Println(t)
					}
				case `level`:
					if !hasTHR_level {
						fmt.Println(t)
					}
				case `value`:
					if !hasTHR_value {
						fmt.Println(t)
					}
				}
			}
			return
		}
	}
	// not in any subchain
	for _, t := range topArgs {
		switch t {
		case `in`:
			if !hasIN {
				fmt.Println(t)
			}
		case `on`:
			if !hasON {
				fmt.Println(t)
			}
		case `with`:
			if !hasWITH {
				fmt.Println(t)
			}
		case `interval`:
			if !hasINTERVAL {
				fmt.Println(t)
			}
		case `inheritance`:
			if !hasINHERITANCE {
				fmt.Println(t)
			}
		case `childrenonly`:
			if !hasCHILDRENONLY {
				fmt.Println(t)
			}
		case `extern`:
			if !hasEXTERN {
				fmt.Println(t)
			}
		default:
			fmt.Println(t)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
