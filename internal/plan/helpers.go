package plan

import "sort"

// Das Manifest kommt als generische JSON-kompatible Werte an (siehe
// internal/manifest.Load) - dieselben Helfer wie in internal/hardware und
// internal/catalog, hier fuer den zusaetzlichen Fall verschachtelter Werte
// hinter interface{}.

func asMap(v interface{}) map[string]interface{} {
	m, _ := v.(map[string]interface{})
	return m
}

func asList(v interface{}) []interface{} {
	l, _ := v.([]interface{})
	return l
}

func asString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func asBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func asFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	default:
		return 0
	}
}

func asInt(v interface{}) int {
	return int(asFloat(v))
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
