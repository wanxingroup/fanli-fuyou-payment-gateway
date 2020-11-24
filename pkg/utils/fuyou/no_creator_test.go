package fuyou

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenFyOrderNo(t *testing.T) {

	tests := []struct {
		inputSeed int64
		want      string
	}{
		{
			inputSeed: 1,
			want:      "1419" + time.Now().Format("20060102150405") + "08100001",
		},
		{
			inputSeed: 0,
			want:      "1419" + time.Now().Format("20060102150405") + "88700002",
		},
		{
			inputSeed: 1,
			want:      "1419" + time.Now().Format("20060102150405") + "08100003",
		},
	}

	for _, test := range tests {

		if test.inputSeed > 0 {
			rand.Seed(test.inputSeed)
		}

		assert.Equal(t, test.want, GeneralFuYouOrderId(), test)
	}
}
