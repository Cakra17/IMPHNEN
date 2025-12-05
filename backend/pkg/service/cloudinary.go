package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld          *cloudinary.Cloudinary
	parentFolder string
}

func NewCloudinaryService(name, apiKey, secret, parentFolder string) (CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(name, apiKey, secret)
	if err != nil {
		return CloudinaryService{}, fmt.Errorf("Failed connect to cloudinary service: %v", err)
	}
	return CloudinaryService{cld: cld, parentFolder: parentFolder}, nil
}

func (s *CloudinaryService) UploadMedia(ctx context.Context, subFolder string, image multipart.File) (secureURL, publidID string, err error) {
	if seeker, ok := image.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			log.Printf("[ERROR] Failed to seek file: %v", err)
			return "", "", fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}
	resp, err := s.cld.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder:         fmt.Sprintf("%s/%s", s.parentFolder, subFolder),
		UniqueFilename: api.Bool(true),
		UseFilename:    api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		log.Printf("[ERROR] Failed to upload photo to claudinary: %v", err)
		return "", "", fmt.Errorf("Failed to save photo: %s", err.Error())
	}

	return resp.SecureURL, resp.PublicID, nil
}

func (s *CloudinaryService) DeleteMedia(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("Failed to delete image: %s", err.Error())
	}

	return nil
}
