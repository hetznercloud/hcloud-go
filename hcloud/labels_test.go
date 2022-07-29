package hcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLabels(t *testing.T) {
	t.Run("correct labels", func(t *testing.T) {
		labelMap := map[string]interface{}{
			"label1":                    "correct.de",
			"label2.de/hallo":           "1correct2.de",
			"label3-test.de/hallo.welt": "233344444443",
			"d/d":                       "d",
			"empty/label":               "",
		}

		ok, err := ValidateResourceLabels(labelMap)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, ok)
	})
	t.Run("incorrect string label values", func(t *testing.T) {
		incorrectLabels := []string{
			"incorrect .com",
			"-incorrect.com",
			"incorrect.com-",
			"incorr,ect.com-",
			"incorrect-111111111111111111111111111111111111111111111111111111111111.com",
		}

		for _, label := range incorrectLabels {
			labelMap := map[string]interface{}{
				"test1": "valid",
				"test2": label,
			}
			ok, err := ValidateResourceLabels(labelMap)
			assert.Error(t, err)
			assert.False(t, ok)
		}
	})
	t.Run("incorrect string label keys", func(t *testing.T) {
		incorrectLabels := []string{
			"incorrect.de/",
			"incor rect.de/",
			"incorrect.de/+",
			"-incorrect.de",
			"incorrect.de-",
			"incorrect.de/tes t",
			"incorrect.de/test-",
			"incorrect.de/test-dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
			"incorrect-11111111111111111111111111111111111111111111111111111111111111111111111111111111" +
				"11111111111111111111111111111111111111111111111111111111111111111111111111111111" +
				"11111111111111111111111111111111111111111111111111111111111111111111111111111111" +
				"11111111111111111111111111111111111111111111111111111111111111111111111111111111" +
				".de/test",
		}

		for _, label := range incorrectLabels {
			labelMap := map[string]interface{}{
				label: "cor-rect.de",
			}
			ok, err := ValidateResourceLabels(labelMap)
			assert.Error(t, err)
			assert.False(t, ok)
		}
	})
}
