// https://gist.github.com/regeda/969a067ff4ed6ffa8ed6
package underscore_test

import (
	"testing"

	"github.com/PavelVershinin/xpg/underscore"
	"github.com/stretchr/testify/assert"
)

func TestUnderscore(t *testing.T) {
	assert.Equal(t, "i_love_golang_and_json_so_much", underscore.Underscore("ILoveGolangAndJSONSoMuch"))
	assert.Equal(t, "i_love_json", underscore.Underscore("ILoveJSON"))
	assert.Equal(t, "json", underscore.Underscore("json"))
	assert.Equal(t, "json", underscore.Underscore("JSON"))
	assert.Equal(t, "привет_мир", underscore.Underscore("ПриветМир"))
}

// BenchmarkUnderscore-4           10000000               209 ns/op
func BenchmarkUnderscore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		underscore.Underscore("TestTable")
	}
}
