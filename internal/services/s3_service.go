package services

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

// S3Service handles AWS S3 operations
type S3Service struct {
	client  *s3.Client
	bucket  string
	region  string
	baseURL string
	logger  *logrus.Logger
}

// NewS3Service creates a new S3 service instance
func NewS3Service(region, bucket, accessKeyID, secretAccessKey, baseURL string, logger *logrus.Logger) (*S3Service, error) {
	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:  client,
		bucket:  bucket,
		region:  region,
		baseURL: baseURL,
		logger:  logger,
	}, nil
}

// UploadProfileImage uploads a profile image to S3 and returns the URL
func (s *S3Service) UploadProfileImage(ctx context.Context, userID int, imageData io.Reader, contentType string) (string, error) {
	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("profiles/%d/profile_%d_%d%s", userID, userID, timestamp, getFileExtension(contentType))

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(filename),
		Body:        imageData,
		ContentType: aws.String(contentType),
		ACL:         "public-read", // Make the image publicly accessible
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to upload profile image to S3")
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	// Generate public URL
	imageURL := fmt.Sprintf("%s/%s", s.baseURL, filename)

	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"filename": filename,
		"url":      imageURL,
	}).Info("Profile image uploaded successfully")

	return imageURL, nil
}

// DeleteProfileImage deletes a profile image from S3
func (s *S3Service) DeleteProfileImage(ctx context.Context, imageURL string) error {
	// Extract key from URL
	key := strings.TrimPrefix(imageURL, s.baseURL+"/")

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.WithError(err).WithField("key", key).Error("Failed to delete profile image from S3")
		return fmt.Errorf("failed to delete image: %w", err)
	}

	s.logger.WithField("key", key).Info("Profile image deleted successfully")
	return nil
}

// getFileExtension returns the appropriate file extension based on content type
func getFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg" // Default to jpg
	}
}

// ValidateImageType checks if the content type is a valid image type
func (s *S3Service) ValidateImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// GetImageURL generates the full URL for an image key
func (s *S3Service) GetImageURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}
