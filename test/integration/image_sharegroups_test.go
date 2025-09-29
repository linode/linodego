package integration

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/linode/linodego"
)

func TestImageSharing_Suite(t *testing.T) {
	client, instance, teardown, err := setupInstance(
		t, "fixtures/TestImageSharing_Suite", true,
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.Region = getRegionsWithCaps(t, client, []string{"Linodes"})[0]
			options.Image = "linode/alpine3.22"
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

	// First, create a Private Image and verify that the `IsShared` and `ImageSharing` fields are as expected
	image, err := client.CreateImage(context.Background(), linodego.ImageCreateOptions{
		DiskID: instanceDisks[0].ID,
		Label:  "linodego-test-image-sharing-image",
	})
	if err != nil {
		t.Errorf("Failed to create image: %v", err)
	}
	t.Cleanup(func() {
		if err := client.DeleteImage(context.Background(), image.ID); err != nil {
			t.Errorf("Failed to delete image %s: %v", image.ID, err)
		}
	})

	if _, err := client.WaitForImageStatus(context.Background(), image.ID, linodego.ImageStatusAvailable, 120); err != nil {
		t.Errorf("Failed to wait for image available status: %v", err)
	}

	if image.IsShared == nil {
		t.Errorf("Expected field 'IsShared' to be present (false) for a new image, got nil")
	} else if *image.IsShared {
		t.Errorf("Expected field 'IsShared' to be false for a new image, got true")
	}

	if image.ImageSharing.SharedWith == nil {
		t.Errorf("Expected ImageSharing.SharedWith to be present, got nil")
	} else {
		expected := linodego.ImageSharingSharedWith{
			ShareGroupCount:   0,
			ShareGroupListURL: fmt.Sprintf("/images/%s/sharegroups", image.ID),
		}
		if *image.ImageSharing.SharedWith != expected {
			t.Errorf("Expected SharedWith to be %v, got %v", expected, *image.ImageSharing.SharedWith)
		}
	}

	if image.ImageSharing.SharedBy != nil {
		t.Errorf("Expected ImageSharing.SharedBy to be nil, got %v", image.ImageSharing.SharedBy)
	}

	// Next, create an empty ImageShareGroup and ensure it has been created successfully
	imageShareGroup, err := client.CreateImageShareGroup(context.Background(), linodego.ImageShareGroupCreateOptions{
		Label: "linodego-test-image-sharing-image-share-group",
	})
	if err != nil {
		t.Errorf("Failed to create image share group: %v", err)
	}
	t.Cleanup(func() {
		if err := client.DeleteImageShareGroup(context.Background(), imageShareGroup.ID); err != nil {
			t.Errorf("Failed to delete image share group %d: %v", imageShareGroup.ID, err)
		}
	})

	// Next, add the previously created Private Image to the ImageShareGroup
	imagesToAdd := linodego.ImageShareGroupAddImagesOptions{
		Images: []linodego.ImageShareGroupImage{
			{
				ID:          image.ID,
				Label:       linodego.Pointer("A label."),
				Description: linodego.Pointer("A description."),
			},
		},
	}

	imagesAdded, err := client.ImageShareGroupAddImages(context.Background(), imageShareGroup.ID, imagesToAdd)
	if err != nil {
		t.Errorf("Failed to add image to image share group: %v", err)
	}

	if len(imagesAdded) != 1 {
		t.Errorf("Expected to add 1 image to imageShareGroup, got %d", len(imagesAdded))
	}

	imageShare := imagesAdded[0]

	if imageShare.IsShared != nil {
		t.Errorf("Expected field 'IsShared' to be nil for an image_share row, got %v", imageShare.IsShared)
	}

	if imageShare.ImageSharing.SharedWith != nil {
		t.Errorf("Expected field 'SharedWith' to be nil for an image_share row, got %v", imageShare.ImageSharing.SharedWith)
	}

	if imageShare.ImageSharing.SharedBy == nil {
		t.Errorf("Expected ImageSharing.SharedBy to be set after sharing, got nil")
	} else {
		sb := imageShare.ImageSharing.SharedBy
		if sb.ShareGroupID != imageShareGroup.ID {
			t.Errorf("Expected SharedBy.ShareGroupID to be %d, got %d", imageShareGroup.ID, sb.ShareGroupID)
		}
		if sb.ShareGroupUUID != imageShareGroup.UUID {
			t.Errorf("Expected SharedBy.ShareGroupUUID to be %s, got %s", imageShareGroup.UUID, sb.ShareGroupUUID)
		}
		if sb.ShareGroupLabel != imageShareGroup.Label {
			t.Errorf("Expected SharedBy.ShareGroupLabel to be %s, got %s", imageShareGroup.Label, sb.ShareGroupLabel)
		}
		if sb.SourceImageID != nil && *sb.SourceImageID != image.ID {
			t.Errorf("Expected SharedBy.SourceImageID to be %s, got %s", image.ID, *sb.SourceImageID)
		}
	}

	image, err = client.GetImage(context.Background(), image.ID)
	if err != nil {
		t.Errorf("Failed to get image: %v", err)
	}

	if image.IsShared == nil {
		t.Errorf("Expected field 'IsShared' to be present (true) for a shared image, got nil")
	} else if !*image.IsShared {
		t.Errorf("Expected field 'IsShared' to be true for a new image, got false")
	}

	if image.ImageSharing.SharedWith == nil {
		t.Errorf("Expected ImageSharing.SharedWith to be present, got nil")
	} else {
		expected := linodego.ImageSharingSharedWith{
			ShareGroupCount:   1,
			ShareGroupListURL: fmt.Sprintf("/images/%s/sharegroups", image.ID),
		}
		if *image.ImageSharing.SharedWith != expected {
			t.Errorf("Expected SharedWith to be %+v, got %+v", expected, *image.ImageSharing.SharedWith)
		}
	}

	if image.ImageSharing.SharedBy != nil {
		t.Errorf("Expected ImageSharing.SharedBy to be nil after sharing, got %v", *image.ImageSharing.SharedBy)
	}

	images, err := client.ListImages(context.Background(), nil)
	if err != nil {
		t.Errorf("Failed to list images: %v", err)
	}
	if !slices.ContainsFunc(images, func(img linodego.Image) bool { return img.ID == image.ID }) {
		t.Errorf("Expected to find image with ID %s in listed images: ", image.ID)
	}

	// Next, list Share Groups where the previously created image exists
	shareGroups, err := client.ListImageShareGroupsContainingPrivateImage(context.Background(), image.ID, nil)
	if err != nil {
		t.Errorf("Failed to list Share Groups: %v", err)
	}
	if !slices.ContainsFunc(shareGroups, func(sg linodego.ProducerImageShareGroup) bool { return sg.ID == imageShareGroup.ID }) {
		t.Errorf("Expected to find Share Group with ID %d in listed Share Groups: ", imageShareGroup.ID)
	}

	// Next, list all Share Groups in the account
	shareGroups, err = client.ListImageShareGroups(context.Background(), nil)
	if err != nil {
		t.Errorf("Failed to list Share Groups: %v", err)
	}
	if !slices.ContainsFunc(shareGroups, func(sg linodego.ProducerImageShareGroup) bool { return sg.ID == imageShareGroup.ID }) {
		t.Errorf("Expected to find Share Group with ID %d in listed Share Groups: ", imageShareGroup.ID)
	}

	// Next, get the Share Group
	shareGroup, err := client.GetImageShareGroup(context.Background(), imageShareGroup.ID)
	if err != nil {
		t.Errorf("Failed to get Share Group: %v", err)
	}
	if shareGroup.ID != imageShareGroup.ID {
		t.Errorf("Expected Share Group with ID %d:", imageShareGroup.ID)
	}

	// Next, list the images shared in the Share Group
	sharedImages, err := client.ImageShareGroupListImages(context.Background(), imageShareGroup.ID, nil)
	if err != nil {
		t.Errorf("Failed to list shared images: %v", err)
	}
	if !slices.ContainsFunc(sharedImages, func(img linodego.Image) bool { return img.ID == imageShare.ID }) {
		t.Errorf("Expected to find shared image with ID %s in listed images: ", imageShare.ID)
	}

	// Next, list the members of the Share Group (should be empty)
	members, err := client.ImageShareGroupListMembers(context.Background(), imageShareGroup.ID, nil)
	if err != nil {
		t.Errorf("Failed to list members: %v", err)
	}
	if len(members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(members))
	}

	// Next, update the Share Group
	imageShareGroup, err = client.UpdateImageShareGroup(context.Background(), imageShareGroup.ID, linodego.ImageShareGroupUpdateOptions{
		Label:       linodego.Pointer("Updated label"),
		Description: linodego.Pointer("Updated description"),
	})
	if err != nil {
		t.Errorf("Failed to update Share Group: %v", err)
	}
	if imageShareGroup.Label != "Updated label" || imageShareGroup.Description != "Updated description" {
		t.Errorf("Failed to update Share Group with ID %d: ", imageShareGroup.ID)
	}

	// Next, update the Shared Image
	imageShareUpdated, err := client.ImageShareGroupUpdateImage(context.Background(), imageShareGroup.ID, imageShare.ID, linodego.ImageShareGroupUpdateImageOptions{
		Label:       linodego.Pointer("Updated label"),
		Description: linodego.Pointer("Updated description"),
	})
	if err != nil {
		t.Errorf("Failed to update Shared Image: %v", err)
	}
	if imageShareUpdated.Label != "Updated label" || imageShareUpdated.Description != "Updated description" {
		t.Errorf("Failed to update Image Share with ID %s: ", imageShareUpdated.ID)
	}

	// Next, remove the image from the Share Group
	err = client.ImageShareGroupRemoveImage(context.Background(), imageShareGroup.ID, imageShareUpdated.ID)
	if err != nil {
		t.Errorf("Failed to remove image from Share Group: %v", err)
	}

	// Check that the image has been removed
	image, err = client.GetImage(context.Background(), image.ID)
	if err != nil {
		t.Errorf("Failed to get image: %v", err)
	}

	if image.ImageSharing.SharedWith == nil {
		t.Errorf("Expected ImageSharing.SharedWith to be present, got nil")
	} else {
		expected := linodego.ImageSharingSharedWith{
			ShareGroupCount:   0,
			ShareGroupListURL: fmt.Sprintf("/images/%s/sharegroups", image.ID),
		}
		if *image.ImageSharing.SharedWith != expected {
			t.Errorf("Expected SharedWith to be %v, got %v", expected, *image.ImageSharing.SharedWith)
		}
	}

	if image.ImageSharing.SharedBy != nil {
		t.Errorf("Expected ImageSharing.SharedBy to be nil, got %v", image.ImageSharing.SharedBy)
	}

	sharedImages, err = client.ImageShareGroupListImages(context.Background(), imageShareGroup.ID, nil)
	if err != nil {
		t.Errorf("Failed to list shared images: %v", err)
	}
	if len(sharedImages) != 0 {
		t.Errorf("Expected 0 shared images, got %d", len(sharedImages))
	}
}
