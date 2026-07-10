package apicompat

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// NormalizeImageGenerationResponsePayload fills fields required by strict
// Responses clients when an upstream returns a completed image item with a
// result but omits status or id.
func NormalizeImageGenerationResponsePayload(payload []byte) ([]byte, bool) {
	if len(payload) == 0 || !gjson.ValidBytes(payload) {
		return payload, false
	}

	updated := payload
	changed := false
	normalizeAt := func(path string, index int) {
		item := gjson.GetBytes(updated, path)
		if !item.IsObject() || strings.TrimSpace(item.Get("type").String()) != "image_generation_call" {
			return
		}
		result := strings.TrimSpace(item.Get("result").String())
		if result == "" {
			return
		}
		if strings.TrimSpace(item.Get("id").String()) == "" {
			if next, err := sjson.SetBytes(updated, path+".id", fallbackImageGenerationCallID(item, index)); err == nil {
				updated = next
				changed = true
			}
		}
		if strings.TrimSpace(item.Get("status").String()) == "" {
			if next, err := sjson.SetBytes(updated, path+".status", "completed"); err == nil {
				updated = next
				changed = true
			}
		}
	}

	if strings.TrimSpace(gjson.GetBytes(updated, "type").String()) == "response.output_item.done" {
		normalizeAt("item", 0)
	}
	for _, prefix := range []string{"response.output", "output"} {
		count := int(gjson.GetBytes(updated, prefix+".#").Int())
		for i := 0; i < count; i++ {
			normalizeAt(fmt.Sprintf("%s.%d", prefix, i), i)
		}
	}

	return updated, changed
}

func fallbackImageGenerationCallID(item gjson.Result, index int) string {
	seed := strings.TrimSpace(item.Get("result").String()) + "\n" +
		strings.TrimSpace(item.Get("revised_prompt").String()) + "\n" +
		fmt.Sprintf("%d", index)
	sum := sha256.Sum256([]byte(seed))
	return fmt.Sprintf("ig_sub2api_%x", sum[:8])
}
