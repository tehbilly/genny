// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/tehbilly/genny

package bugreports

func ContainsString(slice []string, element string) bool {
	return false
}

// ContainsAllString targets github issue 36
func ContainsAllString(slice []string, other []string) bool {
	for _, e := range other {
		if !ContainsString(slice, e) {
			return false
		}
	}
	return true
}
