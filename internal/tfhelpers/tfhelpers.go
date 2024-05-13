package tfhelpers

import (
	"fmt"
	"os"

	"github.com/markphelps/optional"
)

func TFProp(field string, value interface{}, skip_comp interface{}) string {
	if value == skip_comp {
		return ""
	}
	switch val := value.(type) {
	case bool:
		return fmt.Sprintf("\t%s = %t\n", field, value)
	case int:
		return fmt.Sprintf("\t%s = %d\n", field, value)
	case string:
		return fmt.Sprintf("\t%s = \"%s\"\n", field, value)
	case float64:
		return fmt.Sprintf("\t%s = %f\n", field, value)
	case optional.Bool:
		if val.Present() {
			return fmt.Sprintf("\t%s = %t\n", field, val.MustGet())
		}
	case optional.Int:
		if val.Present() {
			return fmt.Sprintf("\t%s = %d\n", field, val.MustGet())
		}
	case optional.Int64:
		if val.Present() {
			return fmt.Sprintf("\t%s = %d\n", field, val.MustGet())
		}
	case optional.String:
		if val.Present() {
			return fmt.Sprintf("\t%s = \"%s\"\n", field, val.MustGet())
		}
	default:
		return fmt.Sprintf("ERRORERROR %T", val)
	}

	return ""
}

func TFUnquote(field string, value string) string {
	return fmt.Sprintf("\t%s = %s\n", field, value)
}

func TFFile(field string, value string, content_path string) string {
	f, e := os.Create(content_path)
	if e != nil {
		return ""
	}
	f.WriteString(value)
	return fmt.Sprintf("\t%s = file(\"%s\")\n", field, content_path)
}
