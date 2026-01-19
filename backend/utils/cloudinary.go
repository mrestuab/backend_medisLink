package utils

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	CloudName = "duyidxlsk"
	ApiKey    = "824833296986543"
	ApiSecret = "x5eqY_Z3K-dyYX8e5AZLS6nRsww"
)

func UploadToCloudinary(file interface{}, folderName string) (string, error) {
	cld, err := cloudinary.NewFromParams(CloudName, ApiKey, ApiSecret)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folderName,
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
