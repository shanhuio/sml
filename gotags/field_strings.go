package gotags

import (
	"fmt"
	"go/ast"
	"strings"
)

// getTypes returns a comma separated list of types in fields. If includeNames
// is true each type is preceded by a comma separated list of parameter names.
func getTypes(fields *ast.FieldList, includeNames bool) string {
	if fields == nil {
		return ""
	}

	types := make([]string, len(fields.List))
	for i, param := range fields.List {
		if len(param.Names) > 0 {
			// there are named parameters, there may be multiple names for a
			// single type
			t := getType(param.Type, true)

			if includeNames {
				// join all the names, followed by their type
				names := make([]string, len(param.Names))
				for j, n := range param.Names {
					names[j] = n.Name
				}
				t = fmt.Sprintf("%s %s", strings.Join(names, ", "), t)
			} else {
				if len(param.Names) > 1 {
					// repeat t len(param.Names) times
					t = strings.Repeat(fmt.Sprintf("%s, ", t), len(param.Names))

					// remove trailing comma and space
					t = t[:len(t)-2]
				}
			}

			types[i] = t
		} else {
			// no named parameters
			types[i] = getType(param.Type, true)
		}
	}

	return strings.Join(types, ", ")
}

// getType returns a string representation of the type of node. If star is true
// and the type is a pointer, a * will be prepended to the string.
func getType(node ast.Node, star bool) (paramType string) {
	switch t := node.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if star {
			return "*" + getType(t.X, star)
		}
		return getType(t.X, star)
	case *ast.SelectorExpr:
		return getType(t.X, star) + "." + getType(t.Sel, star)
	case *ast.ArrayType:
		if l, ok := t.Len.(*ast.BasicLit); ok {
			return fmt.Sprintf("[%s]%s", l.Value, getType(t.Elt, star))
		}
		return "[]" + getType(t.Elt, star)
	case *ast.FuncType:
		fparams := getTypes(t.Params, true)
		fresult := getTypes(t.Results, false)
		if len(fresult) > 0 {
			return fmt.Sprintf("func(%s) %s", fparams, fresult)
		}
		return fmt.Sprintf("func(%s)", fparams)
	case *ast.MapType:
		return fmt.Sprintf(
			"map[%s]%s", getType(t.Key, true), getType(t.Value, true),
		)
	case *ast.ChanType:
		return fmt.Sprintf("chan %s", getType(t.Value, true))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return fmt.Sprintf("...%s", getType(t.Elt, true))
	default:
		return ""
	}
}

// getAccess returns the string "public" if name is considered an exported
// name, otherwise the string "private" is returned.
func getAccess(name string) (access string) {
	if idx := strings.LastIndex(name, "."); idx > -1 && idx < len(name) {
		name = name[idx+1:]
	}

	if ast.IsExported(name) {
		return "public"
	}
	return "private"
}
