package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestExplicitImageControlsRequireAPIKey(t *testing.T) {
	for _, body := range []string{
		`{"tools":[{"type":"image_generation","model":"gpt-image-2","size":"2160x3040","quality":"high"}]}`,
		`{"tools":[{"type":"image_generation","size":"1024x1536"}]}`,
		`{"tools":[{"type":"image_generation","quality":"high"}]}`,
	} {
		require.True(t, HasExplicitOpenAIResponsesImageControls([]byte(body)))
	}
	for _, body := range []string{
		`{"tools":[{"type":"image_generation","size":"auto","quality":"auto"}]}`,
		`{"tools":[{"type":"image_generation","model":"gpt-image-2"}]}`,
		`{"tools":[{"type":"function","name":"image_generation","size":"2160x3040"}]}`,
	} {
		require.False(t, HasExplicitOpenAIResponsesImageControls([]byte(body)))
	}

	oauth := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth}
	setup := &Account{Platform: PlatformOpenAI, Type: AccountTypeSetupToken}
	apiKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey}
	for _, account := range []*Account{oauth, setup} {
		require.False(t, accountSupportsOpenAICapabilities(account, "", OpenAIImagesCapabilityExact))
		require.False(t, account.SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityResponsesImageExact))
		require.True(t, account.SupportsOpenAIImageCapability(OpenAIImagesCapabilityNative))
	}
	require.True(t, accountSupportsOpenAICapabilities(apiKey, "", OpenAIImagesCapabilityExact))
	require.True(t, apiKey.SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityResponsesImageExact))
}

func TestAPIKeyImagesRejectsDecodedSizeMismatchBeforeSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
	image := encodeOpenAIImageTestPNG(t, 16, 32)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"data":[{"b64_json":%q,"size":"64x128"}]}`, image))),
	}
	_, _, _, err := (&OpenAIGatewayService{}).handleOpenAIImagesNonStreamingResponse(
		context.Background(), response, c,
		&Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey},
		&OpenAIImagesRequest{Size: "64x128"},
	)
	require.ErrorContains(t, err, "16x32 differs from requested 64x128")
	require.Empty(t, recorder.Body.String())
}

func TestAPIKeyImagesRejectsReportedQualityDowngrade(t *testing.T) {
	err := validateOpenAIImagesResponseControls(
		&OpenAIImagesRequest{Quality: "high"},
		[]byte(`{"quality":"auto","data":[{"b64_json":"aGVsbG8="}]}`),
	)
	require.ErrorContains(t, err, `quality "auto" differs from requested "high"`)
}

func TestForwardImagesRejectsOAuthExplicitControlsBeforeUpstream(t *testing.T) {
	for _, request := range []*OpenAIImagesRequest{
		{Size: "2160x3040", Quality: "high"},
		{Size: "auto", Quality: "high"},
		{Size: "2160x3040", Quality: "auto"},
	} {
		require.Equal(t, OpenAIImagesCapabilityExact, request.RequiredCapabilityForModel("gpt-image-2"))
		_, err := (&OpenAIGatewayService{}).ForwardImages(context.Background(), nil,
			&Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth}, nil, request, "")
		require.ErrorContains(t, err, "requires an API-key image account")
	}
	require.Equal(t, OpenAIImagesCapabilityNative,
		(&OpenAIImagesRequest{RequiredCapability: OpenAIImagesCapabilityNative, Size: "auto", Quality: "auto"}).RequiredCapabilityForModel("gpt-image-2"))
}
