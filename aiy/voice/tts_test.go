package voice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSay test tts speak function
func TestSay(t *testing.T) {
	err := Say(Speacker{})
	assert.Equal(t, nil, err)
}
