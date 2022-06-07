package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
	"github.com/davecgh/go-spew/spew"
)

var username string
var password string

func init() {
	username = *flag.String("username", os.Getenv("UPCLOUD_USERNAME"), "UpCloud username")
	password = *flag.String("password", os.Getenv("UPCLOUD_PASSWORD"), "UpCloud password")
	rand.Seed(time.Now().Unix())
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

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
	fmt.Println("Creating server")
	details, err := s.CreateServer(&request.CreateServerRequest{
		Hostname: "stuart.example.com",
		Title:    fmt.Sprintf("example-cli-server-%04d", rand.Int31n(1000)),
		Zone:     "fi-hel2",
		Plan:     "1xCPU-1GB",
		StorageDevices: []request.CreateServerStorageDevice{
			{
				Action:  request.CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000050010400",
				Title:   "Centos8 from a template",
				Size:    50,
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
	details, err = s.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         details.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      1 * time.Minute,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to wait for server: %#v", err)
		return err
	}

	fmt.Printf("Created server: %#v\n", details)

	return nil
}

func deleteServers(s *service.Service) error {
	fmt.Println("Getting servers")
	servers, err := s.GetServers()
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
				_, err := s.StopServer(&request.StopServerRequest{
					UUID:     server.UUID,
					StopType: request.ServerStopTypeHard,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to stop server: %#v\n", err)
					return err
				}
				_, err = s.WaitForServerState(&request.WaitForServerStateRequest{
					UUID:         server.UUID,
					DesiredState: upcloud.ServerStateStopped,
					Timeout:      1 * time.Minute,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to wait for server to reach desired state: %#v", err)
					return err
				}
			}
			fmt.Printf("Deleting %s (%s)\n", server.Title, server.UUID)
			err := s.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
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
	fmt.Println("Getting storage")
	storages, err := s.GetStorages(&request.GetStoragesRequest{
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
				err = s.DeleteStorage(&request.DeleteStorageRequest{
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
