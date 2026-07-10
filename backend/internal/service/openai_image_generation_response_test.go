package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestNormalizeOpenAIImageGenerationResponsePayload_OutputItemDone(t *testing.T) {
	payload := []byte(`{"type":"response.output_item.done","item":{"type":"image_generation_call","result":"aW1hZ2U="}}`)

	got, changed := normalizeOpenAIImageGenerationResponsePayload(payload)

	require.True(t, changed)
	require.Equal(t, "completed", gjson.GetBytes(got, "item.status").String())
	require.Contains(t, gjson.GetBytes(got, "item.id").String(), "ig_sub2api_")
	require.Equal(t, "aW1hZ2U=", gjson.GetBytes(got, "item.result").String())
}

func TestNormalizeOpenAIImageGenerationResponsePayload_PreservesCompleteItem(t *testing.T) {
	payload := []byte(`{"type":"response.output_item.done","item":{"id":"ig_upstream","type":"image_generation_call","status":"completed","result":"aW1hZ2U="}}`)

	got, changed := normalizeOpenAIImageGenerationResponsePayload(payload)

	require.False(t, changed)
	require.Equal(t, payload, got)
}

func TestNormalizeOpenAIImageGenerationResponsePayload_TerminalAndJSONOutput(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		path    string
	}{
		{
			name:    "terminal event",
			payload: []byte(`{"type":"response.completed","response":{"output":[{"type":"image_generation_call","result":"dGVybWluYWw="}]}}`),
			path:    "response.output.0",
		},
		{
			name:    "non streaming response",
			payload: []byte(`{"output":[{"type":"image_generation_call","result":"anNvbg=="}]}`),
			path:    "output.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := normalizeOpenAIImageGenerationResponsePayload(tt.payload)

			require.True(t, changed)
			require.Equal(t, "completed", gjson.GetBytes(got, tt.path+".status").String())
			require.NotEmpty(t, gjson.GetBytes(got, tt.path+".id").String())
		})
	}
}
