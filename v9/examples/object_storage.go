package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/google/uuid"
)

const (
	region      = "europe-1"
	waitTimeout = 10 * time.Minute
)

func main() {
	// credentials from env (UPCLOUD_TOKEN or UPCLOUD_USERNAME/UPCLOUD_PASSWORD) or keyring
	client, err := upcloud.NewFromEnv()
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()
	name := fmt.Sprintf("test-sdk-%d", time.Now().Unix())

	fmt.Println("Listing services...")
	listResp, err := client.ListObjectStoragesWithResponse(ctx, nil)
	if err != nil {
		log.Fatalf("list services: %v", err)
	}
	if listResp.JSON200 == nil {
		log.Fatal("list services failed")
	}
	printServices(*listResp.JSON200)

	fmt.Printf("Creating service %s...\n", name)
	createResp, err := client.CreateObjectStorageWithResponse(ctx, upcloud.ObjectStorage2ServiceCreate{
		Name:             name,
		Region:           region,
		ConfiguredStatus: upcloud.Started,
	})
	if err != nil {
		log.Fatalf("create service: %v", err)
	}
	if createResp.JSON201 == nil {
		log.Fatal("create service failed")
	}
	serviceUUID := createResp.JSON201.Uuid
	if serviceUUID == nil {
		log.Fatal("create response missing uuid")
	}
	fmt.Printf("Created service with UUID: %s\n", *serviceUUID)

	fmt.Println("Waiting for service to start...")
	waitCtx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()
	_, err = client.WaitForObjectStorageOperationalState(waitCtx, *serviceUUID, "running")
	if err != nil {
		log.Fatalf("wait for service start: %v", err)
	}
	fmt.Println("Service running")

	fmt.Println("Listing services...")
	listResp2, err := client.ListObjectStoragesWithResponse(ctx, nil)
	if err != nil {
		log.Fatalf("list services: %v", err)
	}
	if listResp2.JSON200 == nil {
		log.Fatal("list services failed")
	}
	printServices(*listResp2.JSON200)

	fmt.Println("Deleting service...")

	u, err := uuid.Parse(*serviceUUID)
	if err != nil {
		log.Fatalf("cannot parse uuid: %v", err)
	}
	delResp, err := client.DeleteObjectStorageWithResponse(ctx, u, nil)
	if err != nil {
		log.Fatalf("delete service: %v", err)
	}
	if delResp.StatusCode() != 200 && delResp.StatusCode() != 204 {
		if delResp.StatusCode() == 404 {
			fmt.Println("Service already deleted")
		} else {
			log.Fatal("delete service failed")
		}
	}

	fmt.Println("Waiting for deletion...")
	delCtx, delCancel := context.WithTimeout(ctx, waitTimeout)
	defer delCancel()
	if err = client.WaitForObjectStorageDeletion(delCtx, *serviceUUID); err != nil {
		log.Fatalf("wait for deletion: %v", err)
	}
	fmt.Println("Deletion confirmed")
}

func printServices(services []upcloud.ObjectStorage2ServiceDetailResponse) {
	if len(services) == 0 {
		fmt.Println("No services found")
		return
	}

	for _, s := range services {
		name, r := "", ""
		if s.Name != nil {
			name = *s.Name
		}
		fmt.Printf("- %s (%s)\n", name, r)
	}
}
