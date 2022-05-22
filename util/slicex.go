package util

import "reflect"

func RemoveDuplicates(slc []string) []string {
	m := make(map[string]bool)
	for _, s := range slc {
		m[s] = true
	}
	var res []string
	for k := range m {
		res = append(res, k)
	}
	return res
}

func ElementInSlice(e interface{}, slc interface{}) bool {
	switch reflect.TypeOf(slc).Kind() {
	case reflect.Slice:
		ss := reflect.ValueOf(slc)
		for i := 0; i < ss.Len(); i++ {
			if reflect.DeepEqual(e, ss.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
