package minimock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CamelToSnake(t *testing.T) {
	assert.Equal(t, "i_love_golang_and_json_so_much", CamelToSnake("ILoveGolangAndJSONSoMuch"))
	assert.Equal(t, "i_love_json", CamelToSnake("ILoveJSON"))
	assert.Equal(t, "json", CamelToSnake("json"))
	assert.Equal(t, "json", CamelToSnake("JSON"))
	assert.Equal(t, "привет_мир", CamelToSnake("ПриветМир"))
}

func Benchmark_CamelToSnake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CamelToSnake("TestTable")
	}
}
