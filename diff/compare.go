package diff

import (
	"math"
	"reflect"
)

func Compare(old, new interface{}) Diff {
	if reflect.TypeOf(old) != reflect.TypeOf(new) {
		return diffValue{new, Changed}
	}

	switch v := old.(type) {
	case []interface{}:
		return compareSlices(v, new.([]interface{}))
	case map[string]interface{}:
		return compareMaps(v, new.(map[string]interface{}))
	default:
		if old != new {
			return diffValue{new, Changed}
		}
	}

	return Unchanged
}

func compareSlices(old, new []interface{}) Diff {
	if reflect.DeepEqual(old, new) {
		return Unchanged
	}

	max := int(math.Max(float64(len(old)), float64(len(new))))
	d := make(diffSlice, max)

	for i := 0; i < max; i++ {
		if i >= len(old) {
			d[i] = diffValue{new[i], Added}
		} else if i >= len(new) {
			d[i] = Removed
		} else {
			d[i] = Compare(old[i], new[i])
		}
	}

	return d
}

func compareMaps(old, new map[string]interface{}) Diff {
	if reflect.DeepEqual(old, new) {
		return Unchanged
	}

	d := make(diffMap)

	// New and updated keys
	for key, value := range new {
		if _, ok := old[key]; !ok {
			d[key] = diffValue{value, Added}
		} else {
			d[key] = Compare(old[key], value)
		}
	}

	// Removed keys
	for key, _ := range old {
		if _, ok := new[key]; !ok {
			d[key] = Removed
		}
	}

	return d
}
