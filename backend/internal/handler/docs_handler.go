package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	docStatusDraft     = "draft"
	docStatusPublished = "published"
	maxDocContentSize  = 2 << 20
	maxDocImageSize    = 8 << 20
)

type DocSummary struct {
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Doc struct {
	DocSummary
	Content string `json:"content"`
}

type docManifest struct {
	Documents []DocSummary `json:"documents"`
}

type docMutationRequest struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	SortOrder int    `json:"sort_order"`
	Content   string `json:"content"`
}

type DocsHandler struct {
	docsDir      string
	manifestPath string
	mu           sync.RWMutex
}

func NewDocsHandler(dataDir string) *DocsHandler {
	docsDir := filepath.Join(dataDir, "docs")
	_ = os.MkdirAll(docsDir, 0o755)
	return &DocsHandler{
		docsDir:      docsDir,
		manifestPath: filepath.Join(docsDir, "manifest.json"),
	}
}

func (h *DocsHandler) ListPublished(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	manifest, err := h.readManifestLocked()
	if err != nil {
		response.InternalError(c, "Failed to load documents")
		return
	}
	docs := make([]DocSummary, 0, len(manifest.Documents))
	for _, doc := range manifest.Documents {
		if doc.Status == docStatusPublished {
			docs = append(docs, doc)
		}
	}
	sortDocs(docs)
	response.Success(c, docs)
}

func (h *DocsHandler) GetPublished(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	doc, err := h.getDocLocked(c.Param("slug"), true)
	if errors.Is(err, os.ErrNotExist) {
		response.NotFound(c, "Document not found")
		return
	}
	if err != nil {
		response.InternalError(c, "Failed to load document")
		return
	}
	response.Success(c, doc)
}

func (h *DocsHandler) ServeImage(c *gin.Context) {
	slug := c.Param("slug")
	filename := strings.TrimPrefix(c.Param("filename"), "/")
	if !validDocSlug(slug) {
		c.Status(http.StatusNotFound)
		return
	}

	h.mu.RLock()
	manifest, err := h.readManifestLocked()
	available := err == nil && findDoc(manifest.Documents, slug) != nil
	h.mu.RUnlock()
	if !available {
		c.Status(http.StatusNotFound)
		return
	}

	imagesDir := filepath.Join(h.docsDir, "assets", slug)
	target, ok := resolvePageImagePath(h.docsDir, imagesDir, filename)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "public, max-age=86400")
	c.File(target)
}

func (h *DocsHandler) ListAdmin(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	manifest, err := h.readManifestLocked()
	if err != nil {
		response.InternalError(c, "Failed to load documents")
		return
	}
	// Keep the empty collection non-nil so JSON clients receive [] instead of null.
	docs := append([]DocSummary{}, manifest.Documents...)
	sortDocs(docs)
	response.Success(c, docs)
}

func (h *DocsHandler) GetAdmin(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	doc, err := h.getDocLocked(c.Param("slug"), false)
	if errors.Is(err, os.ErrNotExist) {
		response.NotFound(c, "Document not found")
		return
	}
	if err != nil {
		response.InternalError(c, "Failed to load document")
		return
	}
	response.Success(c, doc)
}

func (h *DocsHandler) Create(c *gin.Context) {
	var req docMutationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid document payload")
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	req.Title = strings.TrimSpace(req.Title)
	if message := validateDocMutation(req); message != "" {
		response.BadRequest(c, message)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	manifest, err := h.readManifestLocked()
	if err != nil {
		response.InternalError(c, "Failed to load documents")
		return
	}
	if findDoc(manifest.Documents, req.Slug) != nil {
		response.Error(c, http.StatusConflict, "Document slug already exists")
		return
	}

	now := time.Now().UTC()
	meta := DocSummary{
		Slug: req.Slug, Title: req.Title, Status: normalizedDocStatus(req.Status),
		SortOrder: req.SortOrder, CreatedAt: now, UpdatedAt: now,
	}
	if err := h.writeContentLocked(meta.Slug, req.Content); err != nil {
		response.InternalError(c, "Failed to save document")
		return
	}
	manifest.Documents = append(manifest.Documents, meta)
	if err := h.writeManifestLocked(manifest); err != nil {
		_ = os.Remove(h.contentPath(meta.Slug))
		response.InternalError(c, "Failed to save document metadata")
		return
	}
	response.Created(c, Doc{DocSummary: meta, Content: req.Content})
}

func (h *DocsHandler) Update(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if !validDocSlug(slug) {
		response.BadRequest(c, "Invalid document slug")
		return
	}
	var req docMutationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid document payload")
		return
	}
	req.Slug = slug
	req.Title = strings.TrimSpace(req.Title)
	if message := validateDocMutation(req); message != "" {
		response.BadRequest(c, message)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	manifest, err := h.readManifestLocked()
	if err != nil {
		response.InternalError(c, "Failed to load documents")
		return
	}
	meta := findDoc(manifest.Documents, slug)
	if meta == nil {
		response.NotFound(c, "Document not found")
		return
	}
	if err := h.writeContentLocked(slug, req.Content); err != nil {
		response.InternalError(c, "Failed to save document")
		return
	}
	meta.Title = req.Title
	meta.Status = normalizedDocStatus(req.Status)
	meta.SortOrder = req.SortOrder
	meta.UpdatedAt = time.Now().UTC()
	if err := h.writeManifestLocked(manifest); err != nil {
		response.InternalError(c, "Failed to save document metadata")
		return
	}
	response.Success(c, Doc{DocSummary: *meta, Content: req.Content})
}

func (h *DocsHandler) Delete(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if !validDocSlug(slug) {
		response.BadRequest(c, "Invalid document slug")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	manifest, err := h.readManifestLocked()
	if err != nil {
		response.InternalError(c, "Failed to load documents")
		return
	}
	index := -1
	for i := range manifest.Documents {
		if manifest.Documents[i].Slug == slug {
			index = i
			break
		}
	}
	if index < 0 {
		response.NotFound(c, "Document not found")
		return
	}
	manifest.Documents = append(manifest.Documents[:index], manifest.Documents[index+1:]...)
	if err := h.writeManifestLocked(manifest); err != nil {
		response.InternalError(c, "Failed to delete document metadata")
		return
	}
	_ = os.Remove(h.contentPath(slug))
	_ = os.RemoveAll(filepath.Join(h.docsDir, "assets", slug))
	response.Success(c, gin.H{"message": "Document deleted"})
}

func (h *DocsHandler) UploadImage(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if !validDocSlug(slug) {
		response.BadRequest(c, "Invalid document slug")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxDocImageSize+(1<<20))
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Image file is required")
		return
	}
	if fileHeader.Size <= 0 || fileHeader.Size > maxDocImageSize {
		response.BadRequest(c, "Image must be smaller than 8 MB")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	manifest, err := h.readManifestLocked()
	if err != nil || findDoc(manifest.Documents, slug) == nil {
		response.NotFound(c, "Document not found")
		return
	}
	name, err := h.saveImageLocked(slug, fileHeader)
	if err != nil {
		if errors.Is(err, errUnsupportedDocImage) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to upload image")
		return
	}
	publicURL := docImagePublicURL(c, slug, name)
	response.Success(c, gin.H{
		"filename": name,
		"url":      publicURL,
		"markdown": fmt.Sprintf("![%s](%s)", docImageAlt(fileHeader.Filename), publicURL),
	})
}

func docImagePublicURL(c *gin.Context, slug, name string) string {
	prefix, found := strings.CutSuffix(c.Request.URL.Path, "/admin/docs/"+slug+"/images")
	if !found {
		prefix = "/api/v1"
	}
	return strings.TrimRight(prefix, "/") + "/docs/" + url.PathEscape(slug) + "/images/" + url.PathEscape(name)
}

func docImageAlt(filename string) string {
	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	name = strings.ReplaceAll(name, `\`, `\\`)
	name = strings.ReplaceAll(name, "[", `\[`)
	return strings.ReplaceAll(name, "]", `\]`)
}

var errUnsupportedDocImage = errors.New("only PNG, JPEG, WebP and GIF images are supported")

func (h *DocsHandler) saveImageLocked(slug string, header *multipart.FileHeader) (string, error) {
	source, err := header.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = source.Close() }()

	first := make([]byte, 512)
	n, err := io.ReadFull(source, first)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}
	first = first[:n]
	mimeType := http.DetectContentType(first)
	extensionByMIME := map[string]string{
		"image/png": ".png", "image/jpeg": ".jpg", "image/webp": ".webp", "image/gif": ".gif",
	}
	ext, ok := extensionByMIME[mimeType]
	if !ok {
		return "", errUnsupportedDocImage
	}
	if _, err := source.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	dir := filepath.Join(h.docsDir, "assets", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := uuid.NewString() + ext
	target, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	_, copyErr := io.Copy(target, io.LimitReader(source, maxDocImageSize+1))
	closeErr := target.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(filepath.Join(dir, name))
		if copyErr != nil {
			return "", copyErr
		}
		return "", closeErr
	}
	return name, nil
}

func (h *DocsHandler) getDocLocked(slug string, publishedOnly bool) (*Doc, error) {
	if !validDocSlug(slug) {
		return nil, os.ErrNotExist
	}
	manifest, err := h.readManifestLocked()
	if err != nil {
		return nil, err
	}
	meta := findDoc(manifest.Documents, slug)
	if meta == nil || (publishedOnly && meta.Status != docStatusPublished) {
		return nil, os.ErrNotExist
	}
	content, err := os.ReadFile(h.contentPath(slug))
	if err != nil {
		return nil, err
	}
	if len(content) > maxDocContentSize {
		return nil, fmt.Errorf("document content exceeds size limit")
	}
	return &Doc{DocSummary: *meta, Content: string(content)}, nil
}

func (h *DocsHandler) readManifestLocked() (*docManifest, error) {
	data, err := os.ReadFile(h.manifestPath)
	if errors.Is(err, os.ErrNotExist) {
		return &docManifest{Documents: []DocSummary{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var manifest docManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if manifest.Documents == nil {
		manifest.Documents = []DocSummary{}
	}
	return &manifest, nil
}

func (h *DocsHandler) writeManifestLocked(manifest *docManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return atomicWriteFile(h.manifestPath, append(data, '\n'), 0o644)
}

func (h *DocsHandler) writeContentLocked(slug, content string) error {
	return atomicWriteFile(h.contentPath(slug), []byte(content), 0o644)
}

func (h *DocsHandler) contentPath(slug string) string {
	return filepath.Join(h.docsDir, slug+".md")
}

func atomicWriteFile(path string, content []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".docs-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer func() { _ = os.Remove(tempName) }()
	if err := temp.Chmod(perm); err != nil {
		_ = temp.Close()
		return err
	}
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempName, path)
}

func validateDocMutation(req docMutationRequest) string {
	if !validDocSlug(req.Slug) {
		return "Slug must contain only letters, numbers, hyphens and underscores"
	}
	if req.Title == "" || len([]rune(req.Title)) > 200 {
		return "Title is required and must be at most 200 characters"
	}
	if len(req.Content) > maxDocContentSize {
		return "Document content must be smaller than 2 MB"
	}
	if req.Status != "" && req.Status != docStatusDraft && req.Status != docStatusPublished {
		return "Status must be draft or published"
	}
	return ""
}

func validDocSlug(slug string) bool {
	return len(slug) <= 64 && validSlugPattern.MatchString(slug)
}

func normalizedDocStatus(status string) string {
	if status == docStatusPublished {
		return docStatusPublished
	}
	return docStatusDraft
}

func findDoc(docs []DocSummary, slug string) *DocSummary {
	for i := range docs {
		if docs[i].Slug == slug {
			return &docs[i]
		}
	}
	return nil
}

func sortDocs(docs []DocSummary) {
	sort.SliceStable(docs, func(i, j int) bool {
		if docs[i].SortOrder != docs[j].SortOrder {
			return docs[i].SortOrder < docs[j].SortOrder
		}
		return docs[i].CreatedAt.Before(docs[j].CreatedAt)
	})
}

func RegisterDocsRoutes(v1 *gin.RouterGroup, dataDir string, adminAuth gin.HandlerFunc, settingService *service.SettingService) {
	h := NewDocsHandler(dataDir)

	docs := v1.Group("/docs")
	{
		docs.GET("", h.ListPublished)
		docs.GET("/:slug", h.GetPublished)
	}

	v1.GET("/docs/:slug/images/*filename", h.ServeImage)

	adminDocs := v1.Group("/admin/docs")
	adminDocs.Use(adminAuth)
	adminDocs.Use(middleware2.AdminComplianceGuard(settingService))
	{
		adminDocs.GET("", h.ListAdmin)
		adminDocs.GET("/:slug", h.GetAdmin)
		adminDocs.POST("", h.Create)
		adminDocs.PUT("/:slug", h.Update)
		adminDocs.DELETE("/:slug", h.Delete)
		adminDocs.POST("/:slug/images", h.UploadImage)
	}
}
