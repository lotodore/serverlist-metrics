// steam server filters and lists for web requests
package webrequest

// Filters according to https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol#Filter

// Use specific filter data type to prevent messing up filters.
type Filter string

// Without parameters:

const (
	// Servers running dedicated
	FilterDedicated Filter = "\\dedicated\\1"
	// Servers that are not password protected
	FilterNoPassword Filter = "\\password\\0"
	// Servers that are not empty
	FilterNotEmpty Filter = "\\empty\\1"
	// Servers that are not full
	FilterNotFull Filter = "\\full\\1"
	// Servers that are empty
	FilterEmpty Filter = "\\noplayers\\1"
)

// With parameters:

// Servers running the specified map
func GetFilterForMap(mapName string) Filter {
	return Filter("\\map\\" + mapName)
}

// Servers that are running game [appid]
func GetFilterForAppId(appId string) Filter {
	return Filter("\\appid\\" + appId)
}

// Additional filters could be added here.

// Build a parameter string from a slice of filters.
func CreateFilterString(filters []Filter) string {
	// Even though "Filter" is actually "string", Golang does not allow to cast this slice.
	// So we cannot use strings.Join here.
	var result string
	for _, f := range filters {
		result += string(f)
	}
	return result
}
