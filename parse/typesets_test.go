package parse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tehbilly/genny/parse"
)

func TestArgsToTypeset(t *testing.T) {

	args := "Person=dude:man,woman Animal=pet:dog,cat Place=london,paris"
	ts, err := parse.TypeSet(args)

	if assert.NoError(t, err) {
		if assert.Equal(t, 8, len(ts)) {

			assert.Equal(t, "dude", ts[0]["Person"].Alias)
			assert.Equal(t, "man", ts[0]["Person"].Type)
			assert.Equal(t, "pet", ts[0]["Animal"].Alias)
			assert.Equal(t, "dog", ts[0]["Animal"].Type)
			assert.Equal(t, "london", ts[0]["Place"].Alias)
			assert.Equal(t, "london", ts[0]["Place"].Type)

			assert.Equal(t, "man", ts[1]["Person"].Type)
			assert.Equal(t, "dog", ts[1]["Animal"].Type)
			assert.Equal(t, "paris", ts[1]["Place"].Type)

			assert.Equal(t, "man", ts[2]["Person"].Type)
			assert.Equal(t, "cat", ts[2]["Animal"].Type)
			assert.Equal(t, "london", ts[2]["Place"].Type)

			assert.Equal(t, "man", ts[3]["Person"].Type)
			assert.Equal(t, "cat", ts[3]["Animal"].Type)
			assert.Equal(t, "paris", ts[3]["Place"].Type)

			assert.Equal(t, "woman", ts[4]["Person"].Type)
			assert.Equal(t, "dog", ts[4]["Animal"].Type)
			assert.Equal(t, "london", ts[4]["Place"].Type)

			assert.Equal(t, "woman", ts[5]["Person"].Type)
			assert.Equal(t, "dog", ts[5]["Animal"].Type)
			assert.Equal(t, "paris", ts[5]["Place"].Type)

			assert.Equal(t, "woman", ts[6]["Person"].Type)
			assert.Equal(t, "cat", ts[6]["Animal"].Type)
			assert.Equal(t, "london", ts[6]["Place"].Type)

			assert.Equal(t, "woman", ts[7]["Person"].Type)
			assert.Equal(t, "cat", ts[7]["Animal"].Type)
			assert.Equal(t, "paris", ts[7]["Place"].Type)

		}
	}

	ts, err = parse.TypeSet("Person=man Animal=dog Place=london")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ts))
	}
	ts, err = parse.TypeSet("Person=1,2,3,4,5 Animal=1,2,3,4,5 Place=1,2,3,4,5")
	if assert.NoError(t, err) {
		assert.Equal(t, 125, len(ts))
	}
	ts, err = parse.TypeSet("Person=1 Animal=1,2,3,4,5 Place=1,2")
	if assert.NoError(t, err) {
		assert.Equal(t, 10, len(ts))
	}

	ts, err = parse.TypeSet("Person=interface{} Animal=interface{} Place=interface{}")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ts))
		assert.Equal(t, ts[0]["Animal"].Type, "interface{}")
		assert.Equal(t, ts[0]["Person"].Type, "interface{}")
		assert.Equal(t, ts[0]["Place"].Type, "interface{}")
	}

}
