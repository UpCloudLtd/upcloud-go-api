package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
	"github.com/davecgh/go-spew/spew"
)

var username, password string

func init() {
	flag.StringVar(&username, "username", "", "UpCloud username")
	flag.StringVar(&password, "password", "", "UpCloud password")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	if password == "" {
		password = os.Getenv("UPCLOUD_PASSWORD")
	}
	if username == "" {
		username = os.Getenv("UPCLOUD_USERNAME")
	}

	command := flag.Arg(0)

	if len(username) == 0 {
		fmt.Fprintln(os.Stderr, "Username must be specified")
		return 1
	}

	if len(password) == 0 {
		fmt.Fprintln(os.Stderr, "Password must be specified")
		return 2
	}

	fmt.Println("Creating new client")
	c := client.New(username, password)
	s := service.New(c)

	switch command {
	case "deleteservers":
		if err := deleteServers(s); err != nil {
			return 1
		}
	case "deletestorage":
		if err := deleteStorage(s); err != nil {
			return 2
		}
	case "createserver":
		if err := createServer(s); err != nil {
			return 3
		}
	default:
		fmt.Fprintln(os.Stderr, "Unknown command: ", command)
		return 99
	}

	return 0
}

func createServer(s *service.Service) error {
	ctx := context.TODO()
	fmt.Println("Creating server")
	details, err := s.CreateServer(ctx, &request.CreateServerRequest{
		Hostname: "stuart.example.com",
		Title:    "example-cli-server",
		Zone:     "fi-hel2",
		Plan:     "1xCPU-1GB",
		StorageDevices: []request.CreateServerStorageDevice{
			{
				Action:  request.CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000020060100",
				Title:   "Debian GNU/Linux 11 (Bullseye) from a template",
				Size:    10,
				Tier:    upcloud.StorageTierMaxIOPS,
			},
		},
		Networking: &request.CreateServerNetworking{
			Interfaces: []request.CreateServerInterface{
				{
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.NetworkTypeUtility,
				},
			},
		},
		Labels: &upcloud.LabelSlice{
			upcloud.Label{
				Key:   "env",
				Value: "dev",
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create server: %#v\n", err)
		return err
	}
	spew.Println(details)
	if len(details.UUID) == 0 {
		fmt.Fprintf(os.Stderr, "UUID missing")
		return errors.New("UUID too short")
	}
	details, err = s.WaitForServerState(ctx, &request.WaitForServerStateRequest{
		UUID:         details.UUID,
		DesiredState: upcloud.ServerStateStarted,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to wait for server: %#v", err)
		return err
	}

	fmt.Printf("Created server: %#v\n", details)

	return nil
}

func deleteServers(s *service.Service) error {
	ctx := context.TODO()
	fmt.Println("Getting servers")
	servers, err := s.GetServers(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get servers: %#v\n", err)
		return err
	}

	fmt.Printf("Retrieved %d servers\n", len(servers.Servers))

	if len(servers.Servers) > 0 {
		fmt.Println("Deleting all servers")
		for _, server := range servers.Servers {
			if server.State != upcloud.ServerStateStopped {
				fmt.Printf("Server %s (%s) is not stopped. Stopping\n", server.Title, server.UUID)
				_, err := s.StopServer(ctx, &request.StopServerRequest{
					UUID:     server.UUID,
					StopType: request.ServerStopTypeHard,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to stop server: %#v\n", err)
					return err
				}
				_, err = s.WaitForServerState(ctx, &request.WaitForServerStateRequest{
					UUID:         server.UUID,
					DesiredState: upcloud.ServerStateStopped,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to wait for server to reach desired state: %#v", err)
					return err
				}
			}
			fmt.Printf("Deleting %s (%s)\n", server.Title, server.UUID)
			err := s.DeleteServerAndStorages(ctx, &request.DeleteServerAndStoragesRequest{
				UUID: server.UUID,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete server: %#v\n", err)
				return err
			}
			fmt.Printf("Successfully deleted %s (%s)\n", server.Title, server.UUID)
		}
	}

	return nil
}

func deleteStorage(s *service.Service) error {
	ctx := context.TODO()
	fmt.Println("Getting storage")
	storages, err := s.GetStorages(ctx, &request.GetStoragesRequest{
		Access: upcloud.StorageAccessPrivate,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get storages: %#v\n", err)
		return err
	}

	fmt.Printf("Retrieved %d storages\n", len(storages.Storages))

	if len(storages.Storages) > 0 {
		fmt.Println("Deleting all storages")
		for _, storage := range storages.Storages {
			err := errors.New("Dummy")
			for i := 0; err != nil && i < 5; i++ {
				fmt.Printf("%d: Deleting %s (%s)\n", i, storage.Title, storage.UUID)
				err = s.DeleteStorage(ctx, &request.DeleteStorageRequest{
					UUID: storage.UUID,
				})
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete storage: %#v (%s)\n", err, err.Error())
				return err
			}

			fmt.Printf("Successfully deleted %s (%s)\n", storage.Title, storage.UUID)
		}
	}

	return nil
}
