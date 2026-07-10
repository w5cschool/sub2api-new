package service

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
)

// normalizeOpenAIImageGenerationResponsePayload fills fields that Codex needs
// to deserialize and render image_generation_call output items. Some upstream
// response variants omit status and, less commonly, id even when result is
// present.
func normalizeOpenAIImageGenerationResponsePayload(payload []byte) ([]byte, bool) {
	return apicompat.NormalizeImageGenerationResponsePayload(payload)
}
