package parse

import (
	"errors"
	"strings"
)

const (
	typeSep     = " "
	keyValueSep = "="
	valuesSep   = ","
	aliasSep    = ":"
	builtins    = "BUILTINS"
	numbers     = "NUMBERS"
)

type TypeRef struct {
	Alias string
	Type  string
}

func ParseTypeRef(s string) (*TypeRef, error) {
	parts := strings.Split(s, aliasSep)
	switch len(parts) {
	case 1:
		return &TypeRef{Alias: parts[0], Type: parts[0]}, nil
	case 2:
		return &TypeRef{Alias: parts[0], Type: parts[1]}, nil
	default:
		return nil, errors.New("unable to parse type ref: " + s)
	}
}

// TypeSet turns a type string into a []map[string]string
// that can be given to parse.Generics for it to do its magic.
//
// Acceptable args are:
//
//     Person=man
//     Person=dude:man
//     Person=man Animal=dog
//     Person=man Animal=dog Animal2=cat
//     Person=man,woman Animal=dog,cat
//     Person=man,woman,child Animal=dog,cat Place=london,paris
//     Place=London:city.London
func TypeSet(arg string) ([]map[string]TypeRef, error) {

	types := make(map[string][]string)
	var keys []string
	for _, pair := range strings.Split(arg, typeSep) {
		segs := strings.Split(pair, keyValueSep)
		if len(segs) != 2 {
			return nil, &errBadTypeArgs{Arg: arg, Message: "Generic=Specific expected"}
		}
		key := segs[0]
		keys = append(keys, key)
		types[key] = make([]string, 0)
		for _, t := range strings.Split(segs[1], valuesSep) {
			if t == builtins {
				types[key] = append(types[key], Builtins...)
			} else if t == numbers {
				types[key] = append(types[key], Numbers...)
			} else {
				types[key] = append(types[key], t)
			}
		}
	}

	cursors := make(map[string]int)
	for _, key := range keys {
		cursors[key] = 0
	}

	outChan := make(chan map[string]TypeRef)
	go func() {
		buildTypeSet(keys, 0, cursors, types, outChan)
		close(outChan)
	}()

	var typeSets []map[string]TypeRef
	for typeSet := range outChan {
		typeSets = append(typeSets, typeSet)
	}

	return typeSets, nil

}

func buildTypeSet(keys []string, keyI int, cursors map[string]int, types map[string][]string, out chan<- map[string]TypeRef) {
	key := keys[keyI]
	for cursors[key] < len(types[key]) {
		if keyI < len(keys)-1 {
			buildTypeSet(keys, keyI+1, copycursors(cursors), types, out)
		} else {
			// build the typeset for this combination
			ts := make(map[string]TypeRef)
			for k, vals := range types {
				tl := vals[cursors[k]]
				tr, err := ParseTypeRef(tl)
				if err == nil {
					// TODO: Perhaps log if err != nil
					ts[k] = *tr
				}
			}
			out <- ts
		}
		cursors[key]++
	}
}

func copycursors(source map[string]int) map[string]int {
	c := make(map[string]int)
	for k, v := range source {
		c[k] = v
	}
	return c
}
