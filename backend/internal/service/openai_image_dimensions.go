package service

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

const maxOpenAIImageDimensionProbeBytes int64 = 1 << 20

func detectOpenAIImageResultSize(encoded string) string {
	payload := strings.TrimSpace(encoded)
	if strings.HasPrefix(strings.ToLower(payload), "data:") {
		comma := strings.IndexByte(payload, ',')
		if comma < 0 || comma+1 >= len(payload) {
			return ""
		}
		payload = strings.TrimSpace(payload[comma+1:])
	}
	if payload == "" {
		return ""
	}

	for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding} {
		decoded := base64.NewDecoder(encoding, strings.NewReader(payload))
		buffered := bufio.NewReader(io.LimitReader(decoded, maxOpenAIImageDimensionProbeBytes))
		prefix, _ := buffered.Peek(30)
		if width, height, ok := detectOpenAIWebPDimensions(prefix); ok {
			return fmt.Sprintf("%dx%d", width, height)
		}
		cfg, _, err := image.DecodeConfig(buffered)
		if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
			continue
		}
		return fmt.Sprintf("%dx%d", cfg.Width, cfg.Height)
	}
	return ""
}

func detectOpenAIWebPDimensions(header []byte) (int, int, bool) {
	if len(header) < 16 || string(header[:4]) != "RIFF" || string(header[8:12]) != "WEBP" {
		return 0, 0, false
	}

	switch string(header[12:16]) {
	case "VP8X":
		if len(header) < 30 {
			return 0, 0, false
		}
		width := 1 + int(header[24]) + int(header[25])<<8 + int(header[26])<<16
		height := 1 + int(header[27]) + int(header[28])<<8 + int(header[29])<<16
		return width, height, width > 0 && height > 0
	case "VP8 ":
		if len(header) < 30 || string(header[23:26]) != "\x9d\x01\x2a" {
			return 0, 0, false
		}
		width := int(binary.LittleEndian.Uint16(header[26:28]) & 0x3fff)
		height := int(binary.LittleEndian.Uint16(header[28:30]) & 0x3fff)
		return width, height, width > 0 && height > 0
	case "VP8L":
		if len(header) < 25 || header[20] != 0x2f {
			return 0, 0, false
		}
		width := 1 + int(header[21]) + int(header[22]&0x3f)<<8
		height := 1 + int(header[22]>>6) + int(header[23])<<2 + int(header[24]&0x0f)<<10
		return width, height, width > 0 && height > 0
	default:
		return 0, 0, false
	}
}

func reconcileOpenAIResponsesImageResultSizes(results []openAIResponsesImageResult, firstMeta *openAIResponsesImageResult) {
	for i := range results {
		// ChatGPT OAuth can normalize requested controls to "auto". The final
		// image bytes are authoritative for response metadata and tier billing.
		if actualSize := detectOpenAIImageResultSize(results[i].Result); actualSize != "" {
			results[i].Size = actualSize
		}
	}
	if firstMeta == nil || len(results) == 0 {
		return
	}
	if size := strings.TrimSpace(results[0].Size); size != "" {
		firstMeta.Size = size
	}
}

// Check the actual image bytes when available; lifecycle metadata alone can
// claim the requested size even when the image is smaller.
func validateOpenAIImagesResponseControls(request *OpenAIImagesRequest, body []byte) error {
	return validateOpenAIImagesResponseControlsWithEvidence(request, body, false)
}

// The Codex direct endpoint is undocumented and has previously normalized
// explicit controls to auto. Require positive evidence for every requested
// control before returning a successful OAuth image response.
func validateVerifiedOpenAIImagesResponseControls(request *OpenAIImagesRequest, body []byte) error {
	return validateOpenAIImagesResponseControlsWithEvidence(request, body, true)
}

func validateOpenAIImagesResponseControlsWithEvidence(request *OpenAIImagesRequest, body []byte, requireEvidence bool) error {
	if request == nil || !request.RequiresExactControls() || !gjson.ValidBytes(body) {
		return nil
	}
	requestedWidth, requestedHeight, exactPixels := parseImageBillingDimensions(request.Size)
	root := gjson.ParseBytes(body)
	items := root.Get("data").Array()
	if requireEvidence && len(items) == 0 {
		return fmt.Errorf("upstream returned no image to verify explicit controls")
	}
	for _, item := range items {
		if explicitOpenAIImageControl(request.Quality) {
			quality := item.Get("quality").String()
			if quality == "" {
				quality = root.Get("quality").String()
			}
			if requireEvidence && quality == "" {
				return fmt.Errorf("upstream did not report image quality; cannot verify requested %q", request.Quality)
			}
			if quality != "" && !strings.EqualFold(quality, request.Quality) {
				return fmt.Errorf("upstream image quality %q differs from requested %q", quality, request.Quality)
			}
		}
		if !exactPixels {
			if requireEvidence && explicitOpenAIImageControl(request.Size) {
				return fmt.Errorf("cannot verify requested image size %q", request.Size)
			}
			continue
		}
		actual := detectOpenAIImageResultSize(item.Get("b64_json").String())
		if requireEvidence && actual == "" {
			return fmt.Errorf("upstream image bytes do not reveal dimensions; cannot verify requested %q", request.Size)
		}
		if actual == "" {
			actual = strings.TrimSpace(item.Get("size").String())
		}
		if actual == "" {
			actual = strings.TrimSpace(root.Get("size").String())
		}
		if actualWidth, actualHeight, ok := parseImageBillingDimensions(actual); ok &&
			(actualWidth != requestedWidth || actualHeight != requestedHeight) {
			return fmt.Errorf("upstream image size %s differs from requested %s", actual, request.Size)
		}
		if strings.EqualFold(actual, "auto") {
			return fmt.Errorf("upstream image size %q differs from requested %q", actual, request.Size)
		}
	}
	return nil
}
