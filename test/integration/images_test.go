package integration

import (
	"bytes"
	"context"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	. "github.com/linode/linodego"
)

// testImageBytes is a minimal Gzipped image.
// This is necessary because the API will reject invalid images.
var testImageBytes = []byte{
	0x1f, 0x8b, 0x08, 0x08, 0xbd, 0x5c, 0x91, 0x60,
	0x00, 0x03, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x69, 0x6d, 0x67, 0x00, 0x03, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestImage_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestImage_GetMissing")
	defer teardown()

	i, err := client.GetImage(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing image, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing image, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing image, got %v", e.Code)
	}
}

func TestImage_GetFound(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestImage_GetFound")
	defer teardown()

	i, err := client.GetImage(context.Background(), "linode/ubuntu22.04")
	if i.Created == nil || i.EOL == nil || i.Updated == nil {
		t.Errorf("Error parsing time, %v, %v, %v", i.Created, i.EOL, i.Updated)
	}
	if err != nil {
		t.Errorf("Error getting image, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "linode/ubuntu22.04" {
		t.Errorf("Expected a specific image, but got a different one %v", i)
	}
}

func TestImages_List_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestImages_List")
	defer teardown()

	i, err := client.ListImages(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of images, but got none %v", i)
	}
}

func TestImage_Upload(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestImage_Upload")
	defer teardown()

	image, uploadURL, err := client.CreateImageUpload(context.Background(), ImageCreateUploadOptions{
		Region:      "us-ord",
		Label:       "linodego-image-test",
		Description: "An image that does stuff.",
	})
	if err != nil {
		t.Errorf("Failed to create image upload: %v", err)
	}
	defer func() {
		if err := client.DeleteImage(context.Background(), image.ID); err != nil {
			t.Errorf("Failed to delete image %s: %v", image.ID, err)
		}
	}()

	if uploadURL == "" {
		t.Errorf("Expected upload URL, got none")
	}

	if _, err := client.WaitForImageStatus(context.Background(), image.ID, ImageStatusPendingUpload, 60); err != nil {
		t.Errorf("Failed to wait for image pending upload status: %v", err)
	}

	// Because this request currently bypasses the recorder, we should only run it when the recorder is recording
	if testingMode != recorder.ModeReplaying {
		if err := client.UploadImageToURL(context.Background(), uploadURL, bytes.NewReader(testImageBytes)); err != nil {
			t.Errorf("failed to upload image: %v", err)
		}
	}

	if _, err := client.WaitForImageStatus(context.Background(), image.ID, ImageStatusAvailable, 240); err != nil {
		t.Errorf("Failed to wait for image available upload status: %v", err)
	}
}

func TestImage_CreateUpload(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestImage_CreateUpload")
	defer teardown()

	image, uploadURL, err := client.CreateImageUpload(context.Background(), ImageCreateUploadOptions{
		Region:      getRegionsWithCaps(t, client, []string{"Metadata"})[0],
		Label:       "linodego-image-create-upload",
		Description: "An image that does stuff.",
		CloudInit:   true,
	})
	if err != nil {
		t.Errorf("Failed to create image upload: %v", err)
	}
	defer func() {
		if err := client.DeleteImage(context.Background(), image.ID); err != nil {
			t.Errorf("Failed to delete image %s: %v", image.ID, err)
		}
	}()

	assertSliceContains(t, image.Capabilities, "cloud-init")

	if uploadURL == "" {
		t.Errorf("Expected upload URL, got none")
	}
}

func TestImage_CloudInit(t *testing.T) {
	client, instance, teardown, err := setupInstance(
		t, "fixtures/TestImage_CloudInit",
		func(client *Client, options *InstanceCreateOptions) {
			options.Region = getRegionsWithCaps(t, client, []string{"Metadata"})[0]
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(teardown)

	instanceDisks, err := client.ListInstanceDisks(
		context.Background(),
		instance.ID,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	image, err := client.CreateImage(context.Background(), ImageCreateOptions{
		DiskID:    instanceDisks[0].ID,
		Label:     "linodego-test-cloud-init",
		CloudInit: true,
	})
	if err != nil {
		t.Errorf("Failed to create image upload: %v", err)
	}
	t.Cleanup(func() {
		if err := client.DeleteImage(context.Background(), image.ID); err != nil {
			t.Errorf("Failed to delete image %s: %v", image.ID, err)
		}
	})

	assertSliceContains(t, image.Capabilities, "cloud-init")
}
