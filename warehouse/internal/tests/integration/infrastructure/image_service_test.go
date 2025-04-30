//go:build integration

package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"testing"
	imageApp "warehouse/internal/application/product"
	productImageInfra "warehouse/internal/infrastructure/image/product"
	"warehouse/internal/tests/testutils"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	DefaultBucketName = "test-bucket"
	DefaultImageName  = "default-image.png"
	DefaultImageData  = "default-image-data"
)

type ImageServiceTestSuite struct {
	suite.Suite

	ctx       context.Context
	testMinio *testutils.TestMinio
	config    *productImageInfra.Config
}

func (s *ImageServiceTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.config = &productImageInfra.Config{
		BucketName:       DefaultBucketName,
		DefaultImageName: DefaultImageName,
	}

	minioC, err := testutils.NewTestMinio(s.ctx)
	require.NoError(s.T(), err)
	s.testMinio = minioC

	err = s.testMinio.Client.MakeBucket(s.ctx, s.config.BucketName, minio.MakeBucketOptions{})
	require.NoError(s.T(), err)

	_, err = s.testMinio.Client.PutObject(s.ctx, s.config.BucketName, s.config.DefaultImageName,
		bytes.NewReader([]byte(DefaultImageData)), -1, minio.PutObjectOptions{})
	require.NoError(s.T(), err)
}

func (s *ImageServiceTestSuite) TearDownSuite() {
	if s.testMinio != nil {
		err := s.testMinio.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *ImageServiceTestSuite) getService() imageApp.ImageService {
	return productImageInfra.NewImageService(s.config, s.testMinio.Client)
}

func (s *ImageServiceTestSuite) TestCreateImage() {
	defaultHash, err := testutils.GetFileHash(bytes.NewReader([]byte(DefaultImageData)))
	require.NoError(s.T(), err)

	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "Success",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			service := s.getService()

			path, err := service.Create(s.ctx)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				reader, err := s.testMinio.Client.GetObject(s.ctx, s.config.BucketName, path, minio.GetObjectOptions{})
				require.NoError(s.T(), err)
				defer reader.Close()

				downloadedHash, hashErr := testutils.GetFileHash(reader)
				require.NoError(s.T(), hashErr)
				require.Equal(s.T(), defaultHash, downloadedHash)
			}
		})
	}
}

func (s *ImageServiceTestSuite) TestUpdateImage() {
	tests := []struct {
		name          string
		fileName      string
		fileData      []byte
		fileType      string
		expectedError error
	}{
		{
			name:          "Success",
			fileName:      "upload-image-success.txt",
			fileData:      []byte("upload test data"),
			fileType:      "text/plain",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			service := s.getService()

			err := service.Update(s.ctx, tc.fileName, bytes.NewReader(tc.fileData), tc.fileType)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				reader, err := s.testMinio.Client.GetObject(s.ctx, s.config.BucketName, tc.fileName, minio.GetObjectOptions{})
				require.NoError(s.T(), err)
				defer reader.Close()

				originalHash, err := testutils.GetFileHash(bytes.NewReader(tc.fileData))
				require.NoError(s.T(), err)

				downloadedHash, hashErr := testutils.GetFileHash(reader)
				require.NoError(s.T(), hashErr)
				require.Equal(s.T(), originalHash, downloadedHash)
			}
		})
	}
}

func (s *ImageServiceTestSuite) TestGetImage() {
	tests := []struct {
		name          string
		fileName      string
		contentType   string
		setup         func(fileName string, contentType string) string
		expectedError error
	}{
		{
			name:        "Success",
			fileName:    "get-image-success.txt",
			contentType: "image/png",
			setup: func(fileName string, contentType string) string {
				fileData := []byte("download test data")
				options := minio.PutObjectOptions{
					ContentType: contentType,
				}

				_, err := s.testMinio.Client.PutObject(s.ctx, s.config.BucketName, fileName,
					bytes.NewReader(fileData), -1, options)
				require.NoError(s.T(), err)

				originalHash, err := testutils.GetFileHash(bytes.NewReader(fileData))
				require.NoError(s.T(), err)

				return originalHash
			},
			expectedError: nil,
		},
		{
			name:        "Failure: Image not found",
			fileName:    "get-image-failure-not-found.txt",
			contentType: "image/png",
			setup: func(_, _ string) string {
				return ""
			},
			expectedError: errors.New("image service error"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			service := s.getService()
			originalHash := tc.setup(tc.fileName, tc.contentType)

			reader, contentType, err := service.Get(s.ctx, tc.fileName)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.ErrorContains(s.T(), err, tc.expectedError.Error())
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), reader)
				defer reader.Close()

				require.Equal(s.T(), tc.contentType, contentType)

				downloadedHash, hashErr := testutils.GetFileHash(reader)
				require.NoError(s.T(), hashErr)
				require.Equal(s.T(), originalHash, downloadedHash)
			}
		})
	}
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTestSuite))
}
