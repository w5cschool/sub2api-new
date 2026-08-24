package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestDocsHandlerLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDocsHandler(t.TempDir())
	router := gin.New()
	router.GET("/docs", h.ListPublished)
	router.GET("/docs/:slug", h.GetPublished)
	router.GET("/admin/docs", h.ListAdmin)
	router.POST("/admin/docs", h.Create)
	router.PUT("/admin/docs/:slug", h.Update)
	router.DELETE("/admin/docs/:slug", h.Delete)

	created := performDocJSONRequest(t, router, http.MethodPost, "/admin/docs", docMutationRequest{
		Slug: "getting-started", Title: "Getting started", Status: docStatusPublished,
		SortOrder: 2, Content: "# Hello\n\nWelcome.",
	})
	require.Equal(t, http.StatusCreated, created.Code)

	list := performRequest(router, http.MethodGet, "/docs", nil, "")
	require.Equal(t, http.StatusOK, list.Code)
	var listEnvelope struct {
		Data []DocSummary `json:"data"`
	}
	require.NoError(t, json.Unmarshal(list.Body.Bytes(), &listEnvelope))
	require.Len(t, listEnvelope.Data, 1)
	require.Equal(t, "getting-started", listEnvelope.Data[0].Slug)

	get := performRequest(router, http.MethodGet, "/docs/getting-started", nil, "")
	require.Equal(t, http.StatusOK, get.Code)
	var getEnvelope struct {
		Data Doc `json:"data"`
	}
	require.NoError(t, json.Unmarshal(get.Body.Bytes(), &getEnvelope))
	require.Equal(t, "# Hello\n\nWelcome.", getEnvelope.Data.Content)

	updated := performDocJSONRequest(t, router, http.MethodPut, "/admin/docs/getting-started", docMutationRequest{
		Title: "Draft guide", Status: docStatusDraft, SortOrder: 1, Content: "Draft",
	})
	require.Equal(t, http.StatusOK, updated.Code)
	require.Equal(t, http.StatusNotFound, performRequest(router, http.MethodGet, "/docs/getting-started", nil, "").Code)

	deleted := performRequest(router, http.MethodDelete, "/admin/docs/getting-started", nil, "")
	require.Equal(t, http.StatusOK, deleted.Code)
	require.Equal(t, http.StatusNotFound, performRequest(router, http.MethodGet, "/docs/getting-started", nil, "").Code)
}

func TestDocsHandlerEmptyAdminListReturnsJSONArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDocsHandler(t.TempDir())
	router := gin.New()
	router.GET("/admin/docs", h.ListAdmin)

	list := performRequest(router, http.MethodGet, "/admin/docs", nil, "")
	require.Equal(t, http.StatusOK, list.Code)

	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(list.Body.Bytes(), &envelope))
	require.JSONEq(t, `[]`, string(envelope.Data))
}

func TestDocsHandlerUploadImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewDocsHandler(t.TempDir())
	router := gin.New()
	router.POST("/admin/docs", h.Create)
	router.POST("/admin/docs/:slug/images", h.UploadImage)
	router.GET("/docs/:slug/images/*filename", h.ServeImage)

	created := performDocJSONRequest(t, router, http.MethodPost, "/admin/docs", docMutationRequest{
		Slug: "images", Title: "Images", Status: docStatusPublished, Content: "# Images",
	})
	require.Equal(t, http.StatusCreated, created.Code)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "pixel.png")
	require.NoError(t, err)
	// Valid 1x1 transparent PNG.
	png := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0xf0, 0x1f, 0x00, 0x05, 0x00, 0x01, 0xff, 0x89, 0x99, 0x3d, 0x1d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}
	_, err = part.Write(png)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	upload := performRequest(router, http.MethodPost, "/admin/docs/images/images", &body, writer.FormDataContentType())
	require.Equal(t, http.StatusOK, upload.Code, upload.Body.String())
	var envelope struct {
		Data struct {
			Filename string `json:"filename"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(upload.Body.Bytes(), &envelope))
	require.NotEmpty(t, envelope.Data.Filename)

	image := performRequest(router, http.MethodGet, "/docs/images/images/"+envelope.Data.Filename, nil, "")
	require.Equal(t, http.StatusOK, image.Code)
	require.Equal(t, png, image.Body.Bytes())
}

func performDocJSONRequest(t *testing.T, router http.Handler, method, path string, payload docMutationRequest) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	return performRequest(router, method, path, bytes.NewReader(data), "application/json")
}

func performRequest(router http.Handler, method, path string, body io.Reader, contentType string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}
