package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/Cakra17/imphnen/internal/models"
)

type KolosalService struct {
	API_KEY string
}

func NewKolosalService(apiKey string) KolosalService {
	return KolosalService{API_KEY: apiKey}
}

func (s *KolosalService) OCRForm(image multipart.File, filename string) (models.OCR, error) {
	if s.API_KEY == "" {
		return models.OCR{}, fmt.Errorf("KOLOSAL_API_KEY not set")
	}

	if seeker, ok := image.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return models.OCR{}, fmt.Errorf("failed to seek file: %w", err)
		}
	}

	buffer := make([]byte, 512)
	n, err := image.Read(buffer)
	if err != nil && err != io.EOF {
		return models.OCR{}, fmt.Errorf("failed to read file for type detection: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])

	if seeker, ok := image.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return models.OCR{}, fmt.Errorf("failed to seek file: %w", err)
		}
	}

	if !strings.HasPrefix(contentType, "image/") {
		return models.OCR{}, fmt.Errorf("file is not an image, detected type: %s", contentType)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return models.OCR{}, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, image); err != nil {
		return models.OCR{}, fmt.Errorf("failed to copy image data: %w", err)
	}

	if err := writer.WriteField("invoice", "true"); err != nil {
		return models.OCR{}, fmt.Errorf("failed to write invoice field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return models.OCR{}, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.kolosal.ai/ocr/form", &body)
	if err != nil {
		return models.OCR{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.API_KEY))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.OCR{}, fmt.Errorf("OCR API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.OCR{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.OCR{}, fmt.Errorf("OCR API returned status %d: %s",
			resp.StatusCode, string(respBody))
	}

	var payload models.OCR
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return models.OCR{}, fmt.Errorf("failed to decode OCR response: %w, body: %s",
			err, string(respBody))
	}

	return payload, nil
}
