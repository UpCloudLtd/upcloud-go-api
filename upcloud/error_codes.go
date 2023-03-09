package upcloud

// List of possible error codes
// This can be used for comparison against `Problem.ErrorCode()`
const (
	ErrCodeAuthenticationFailed                           string = "AUTHENTICATION_FAILED"
	ErrCodeActionInvalid                                  string = "ACTION_INVALID"
	ErrCodeActionMissing                                  string = "ACTION_MISSING"
	ErrCodeBootOrderInvalid                               string = "BOOT_ORDER_INVALID"
	ErrCodeCoreMemoryUnsupported                          string = "CORE_MEMORY_UNSUPPORTED"
	ErrCodeCoreNumberInvalid                              string = "CORE_NUMBER_INVALID"
	ErrCodeCreatePasswordInvalid                          string = "CREATE_PASSWORD_INVALID"
	ErrCodeFirewallInvalid                                string = "FIREWALL_INVALID"
	ErrCodeHostnameInvalid                                string = "HOSTNAME_INVALID"
	ErrCodeHostnameMissing                                string = "HOSTNAME_MISSING"
	ErrCodeMemoryAmountInvalid                            string = "MEMORY_AMOUNT_INVALID"
	ErrCodeNicModelInvalid                                string = "NIC_MODEL_INVALID"
	ErrCodePasswordDeliveryInvalid                        string = "PASSWORD_DELIVERY_INVALID"
	ErrCodeServerTitleInvalid                             string = "SERVER_TITLE_INVALID"
	ErrCodeServerTitleMissing                             string = "SERVER_TITLE_MISSING"
	ErrCodeSimpleBackupInvalid                            string = "SIMPLE_BACKUP_INVALID"
	ErrCodeSizeInvalid                                    string = "SIZE_INVALID"
	ErrCodeSizeMissing                                    string = "SIZE_MISSING"
	ErrCodeStorageDeviceInvalid                           string = "STORAGE_DEVICE_INVALID"
	ErrCodeStorageDeviceMissing                           string = "STORAGE_DEVICE_MISSING"
	ErrCodeStorageDevicesInvalid                          string = "STORAGE_DEVICES_INVALID"
	ErrCodeStorageDevicesMissing                          string = "STORAGE_DEVICES_MISSING"
	ErrCodeStorageInvalid                                 string = "STORAGE_INVALID"
	ErrCodeStorageMissing                                 string = "STORAGE_MISSING"
	ErrCodeStorageTitleInvalid                            string = "STORAGE_TITLE_INVALID"
	ErrCodeStorageTitleMissing                            string = "STORAGE_TITLE_MISSING"
	ErrCodeTimezoneInvalid                                string = "TIMEZONE_INVALID"
	ErrCodeTypeInvalid                                    string = "TYPE_INVALID"
	ErrCodeTierInvalid                                    string = "TIER_INVALID"
	ErrCodeUserDataInvalid                                string = "USER_DATA_INVALID"
	ErrCodeVideoModelInvalid                              string = "VIDEO_MODEL_INVALID"
	ErrCodeVncInvalid                                     string = "VNC_INVALID"
	ErrCodeVncPasswordInvalid                             string = "VNC_PASSWORD_INVALID"
	ErrCodeZoneInvalid                                    string = "ZONE_INVALID"
	ErrCodeZoneMissing                                    string = "ZONE_MISSING"
	ErrCodeInvalidLabelKey                                string = "INVALID_LABEL_KEY"
	ErrCodeInvalidLabelValue                              string = "INVALID_LABEL_VALUE"
	ErrCodeInsufficientCredits                            string = "INSUFFICIENT_CREDITS"
	ErrCodeStorageForbidden                               string = "STORAGE_FORBIDDEN"
	ErrCodeHostForbidden                                  string = "HOST_FORBIDDEN"
	ErrCodePlanCoreNumberIllegal                          string = "PLAN_CORE_NUMBER_ILLEGAL"
	ErrCodePlanMemoryAmountIllegal                        string = "PLAN_MEMORY_AMOUNT_ILLEGAL"
	ErrCodeTrialPlan                                      string = "TRIAL_PLAN"
	ErrCodeStorageNotFound                                string = "STORAGE_NOT_FOUND"
	ErrCodeHostNotFound                                   string = "HOST_NOT_FOUND"
	ErrCodeZoneNotFound                                   string = "ZONE_NOT_FOUND"
	ErrCodeCdromDeviceInUse                               string = "CDROM_DEVICE_IN_USE"
	ErrCodeDeviceAddressInUse                             string = "DEVICE_ADDRESS_IN_USE"
	ErrCodeIpAddressResourcesUnavailable                  string = "IP_ADDRESS_RESOURCES_UNAVAILABLE"
	ErrCodeMultipleTemplates                              string = "MULTIPLE_TEMPLATES"
	ErrCodePublicStorageAttach                            string = "PUBLIC_STORAGE_ATTACH"
	ErrCodeServerResourcesUnavailable                     string = "SERVER_RESOURCES_UNAVAILABLE"
	ErrCodeStorageAttachedAsCdrom                         string = "STORAGE_ATTACHED_AS_CDROM"
	ErrCodeStorageAttachedAsDisk                          string = "STORAGE_ATTACHED_AS_DISK"
	ErrCodeStorageDeviceLimitReached                      string = "STORAGE_DEVICE_LIMIT_REACHED"
	ErrCodeStorageInUse                                   string = "STORAGE_IN_USE"
	ErrCodeStorageResourcesUnavailable                    string = "STORAGE_RESOURCES_UNAVAILABLE"
	ErrCodeStorageStateIllegal                            string = "STORAGE_STATE_ILLEGAL"
	ErrCodeStorageTypeIllegal                             string = "STORAGE_TYPE_ILLEGAL"
	ErrCodeZoneMismatch                                   string = "ZONE_MISMATCH"
	ErrCodeInvalidUsername                                string = "INVALID_USERNAME"
	ErrCodeServerCreatingLimitReached                     string = "SERVER_CREATING_LIMIT_REACHED"
	ErrCodeTooManyBootDisks                               string = "TOO_MANY_BOOT_DISKS"
	ErrCodeWindowsNotAvailable                            string = "WINDOWS_NOT_AVAILABLE"
	ErrCodeMsLmaRequired                                  string = "MS_LMA_REQUIRED"
	ErrCodeBackupRuleConflict                             string = "BACKUP_RULE_CONFLICT"
	ErrCodeCanNotCreatePasswordForTheSelectedDistribution string = "CAN_NOT_CREATE_PASSWORD_FOR_THE_SELECTED_DISTRIBUTION"
	ErrCodeMetadataDisabledOnCloudInit                    string = "METADATA_DISABLED_ON_CLOUD-INIT"
	ErrCodeServerInvalid                                  string = "SERVER_INVALID"
	ErrCodeServerNotFound                                 string = "SERVER_NOT_FOUND"
	ErrCodeServerForbidden                                string = "SERVER_FORBIDDEN"
	ErrCodeNoStoragesAttached                             string = "NO_STORAGES_ATTACHED"
	ErrCodeServerStateIllegal                             string = "SERVER_STATE_ILLEGAL"
	ErrCodeStopTypeInvalid                                string = "STOP_TYPE_INVALID"
	ErrCodeTimeoutInvalid                                 string = "TIMEOUT_INVALID"
	ErrCodeTimeoutMissing                                 string = "TIMEOUT_MISSING"
	ErrCodeTimeoutActionInvalid                           string = "TIMEOUT_ACTION_INVALID"
	ErrCodeTitleInvalid                                   string = "TITLE_INVALID"
	ErrCodeFailedToAddMemory                              string = "FAILED_TO_ADD_MEMORY"
	ErrCodeHotResizeNotEnabled                            string = "HOT_RESIZE_NOT_ENABLED"
	ErrCodeHotResizeUnavailable                           string = "HOT_RESIZE_UNAVAILABLE"
	ErrCodeNoAvailableMemorySlots                         string = "NO_AVAILABLE_MEMORY_SLOTS"
	ErrCodeRequestInvalid                                 string = "REQUEST_INVALID"
	ErrCodeStorageDeletionPolicyInvalid                   string = "STORAGE_DELETION_POLICY_INVALID"
	ErrCodeGroupInvalid                                   string = "GROUP_INVALID"
	ErrCodeGroupForbidden                                 string = "GROUP_FORBIDDEN"
	ErrCodeGroupNotFound                                  string = "GROUP_NOT_FOUND"
	ErrCodeDuplicateLabelKeys                             string = "DUPLICATE_LABEL_KEYS"
	ErrCodeAlreadyInServerGroup                           string = "ALREADY_IN_SERVER_GROUP"
	ErrCodeBackupRuleInvalid                              string = "BACKUP_RULE_INVALID"
	ErrCodeIntervalInvalid                                string = "INTERVAL_INVALID"
	ErrCodeIntervalMissing                                string = "INTERVAL_MISSING"
	ErrCodeHourInvalid                                    string = "HOUR_INVALID"
	ErrCodeHourMissing                                    string = "HOUR_MISSING"
	ErrCodeRetentionInvalid                               string = "RETENTION_INVALID"
	ErrCodeRetentionMissing                               string = "RETENTION_MISSING"
	ErrCodeTitleMissing                                   string = "TITLE_MISSING"
	ErrCodeStorageAttached                                string = "STORAGE_ATTACHED"
	ErrCodeResizeFailed                                   string = "RESIZE_FAILED"
	ErrCodeAddressInvalid                                 string = "ADDRESS_INVALID"
	ErrCodeCDROMHotplugUnsupported                        string = "CDROM_HOTPLUG_UNSUPPORTED"
	ErrCodeIdeHotplugUnsupported                          string = "IDE_HOTPLUG_UNSUPPORTED"
	ErrCodeAddressMissing                                 string = "ADDRESS_MISSING"
	ErrCodeAddressOrUuidRequired                          string = "ADDRESS_OR_UUID_REQUIRED"
	ErrCodeDeviceAddressNotInUse                          string = "DEVICE_ADDRESS_NOT_IN_USE"
	ErrCodeHotplugFailed                                  string = "HOTPLUG_FAILED"
	ErrCodeNoCdromDevice                                  string = "NO_CDROM_DEVICE"
	ErrCodeCdromEjectFailed                               string = "CDROM_EJECT_FAILED"
	ErrCodeImportUnavailable                              string = "IMPORT_UNAVAILABLE"
	ErrCodeStorageInconsistent                            string = "STORAGE_INCONSISTENT"
	ErrCodeStorageImportNotFound                          string = "STORAGE_IMPORT_NOT_FOUND"
	ErrCodeStorageImportNotInProgress                     string = "STORAGE_IMPORT_NOT_IN_PROGRESS"
	ErrCodeUnableToCancel                                 string = "UNABLE_TO_CANCEL"
	ErrCodeStorageTierIllegal                             string = "STORAGE_TIER_ILLEGAL"
	ErrCodeBackupDeletionPolicyInvalid                    string = "BACKUP_DELETION_POLICY_INVALID"
	ErrCodeIpAddressLimitReached                          string = "IP_ADDRESS_LIMIT_REACHED"
	ErrCodeIpAddressInvalid                               string = "IP_ADDRESS_INVALID"
	ErrCodePtrRecordInvalid                               string = "PTR_RECORD_INVALID"
	ErrCodeMacInvalid                                     string = "MAC_INVALID"
	ErrCodeIpAddressForbidden                             string = "IP_ADDRESS_FORBIDDEN"
	ErrCodeIpAddressNotFound                              string = "IP_ADDRESS_NOT_FOUND"
	ErrCodePtrRecordNotSupported                          string = "PTR_RECORD_NOT_SUPPORTED"
	ErrCodeFloatingIpNotAvailable                         string = "FLOATING_IP_NOT_AVAILABLE"
	ErrCodeCannotDeletePrivateAddress                     string = "CANNOT_DELETE_PRIVATE_ADDRESS"
	ErrCodeDirectionInvalid                               string = "DIRECTION_INVALID"
	ErrCodeDirectionMissing                               string = "DIRECTION_MISSING"
	ErrCodeICMPTypeInvalid                                string = "ICMP_TYPE_INVALID"
	ErrCodeDestinationAddressOrderIllegal                 string = "DESTINATION_ADDRESS_ORDER_ILLEGAL"
	ErrCodeDestinationAddressStartInvalid                 string = "DESTINATION_ADDRESS_START_INVALID"
	ErrCodeDestinationAddressEndInvalid                   string = "DESTINATION_ADDRESS_END_INVALID"
	ErrCodeDestinationPortOrderIllegal                    string = "DESTINATION_PORT_ORDER_ILLEGAL"
	ErrCodeDestinationPortStartInvalid                    string = "DESTINATION_PORT_START_INVALID"
	ErrCodeDestinationPortEndInvalid                      string = "DESTINATION_PORT_END_INVALID"
	ErrCodeICMPTypeProtocolMismatch                       string = "ICMP_TYPE_PROTOCOL_MISMATCH"
	ErrCodePortProtocolMismatch                           string = "PORT_PROTOCOL_MISMATCH"
	ErrCodePositionInvalid                                string = "POSITION_INVALID"
	ErrCodeProtocolInvalid                                string = "PROTOCOL_INVALID"
	ErrCodeSourceAddressOrderIllegal                      string = "SOURCE_ADDRESS_ORDER_ILLEGAL"
	ErrCodeSourceAddressStartInvalid                      string = "SOURCE_ADDRESS_START_INVALID"
	ErrCodeSourceAddressEndInvalid                        string = "SOURCE_ADDRESS_END_INVALID"
	ErrCodeSourcePortOrderIllegal                         string = "SOURCE_PORT_ORDER_ILLEGAL"
	ErrCodeSourcePortStartInvalid                         string = "SOURCE_PORT_START_INVALID"
	ErrCodeSourcePortEndInvalid                           string = "SOURCE_PORT_END_INVALID"
	ErrCodeCommentInvalid                                 string = "COMMENT_INVALID"
	ErrCodeFirewallRuleExists                             string = "FIREWALL_RULE_EXISTS"
	ErrCodeFirewallRuleLimitReached                       string = "FIREWALL_RULE_LIMIT_REACHED"
	ErrCodeFirewallRuleNotFound                           string = "FIREWALL_RULE_NOT_FOUND"
	ErrCodeTagInvalid                                     string = "TAG_INVALID"
	ErrCodeTagForbidden                                   string = "TAG_FORBIDDEN"
	ErrCodeTagExists                                      string = "TAG_EXISTS"
	ErrCodeTagNotFound                                    string = "TAG_NOT_FOUND"
	ErrCodeNetworkNameInvalid                             string = "NETWORK_NAME_INVALID"
	ErrCodeNetworkNameMissing                             string = "NETWORK_NAME_MISSING"
	ErrCodeNetworkTypeInvalid                             string = "NETWORK_TYPE_INVALID"
	ErrCodeNetworkTypeMissing                             string = "NETWORK_TYPE_MISSING"
	ErrCodeNetworkForbidden                               string = "NETWORK_FORBIDDEN"
	ErrCodeNetworkNotFound                                string = "NETWORK_NOT_FOUND"
	ErrCodeNetworkInvalid                                 string = "NETWORK_INVALID"
	ErrCodeUnknownAttribute                               string = "UNKNOWN_ATTRIBUTE"
	ErrCodeNetworkNotEmpty                                string = "NETWORK_NOT_EMPTY"
	ErrCodeInterfaceExists                                string = "INTERFACE_EXISTS"
	ErrCodeInterfaceNotFound                              string = "INTERFACE_NOT_FOUND"
	ErrCodeInterfaceForbidden                             string = "INTERFACE_FORBIDDEN"
	ErrCodeAddressAttached                                string = "ADDRESS_ATTACHED"
	ErrCodeDualStackInterface                             string = "DUAL_STACK_INTERFACE"
	ErrCodeRouterNotFound                                 string = "ROUTER_NOT_FOUND"
	ErrCodeRouterNameMissing                              string = "ROUTER_NAME_MISSING"
	ErrCodeDuplicateRoute                                 string = "DUPLICATE_ROUTE"
	ErrCodeNexthopInvalid                                 string = "NEXTHOP_INVALID"
	ErrCodeStaticRouteTargetInvalid                       string = "STATIC_ROUTE_TARGET_INVALID"
	ErrCodeInvalidRoute                                   string = "INVALID_ROUTE"
	ErrCodeInvalidRouteFamily                             string = "INVALID_ROUTE_FAMILY"
	ErrCodeStaticRouteLimitReached                        string = "STATIC_ROUTE_LIMIT_REACHED"
	ErrCodeRouterAttached                                 string = "ROUTER_ATTACHED"
	ErrCodeZoneHostForbidden                              string = "ZONE_HOST_FORBIDDEN"
	ErrCodeAccessKeyInvalid                               string = "ACCESS_KEY_INVALID"
	ErrCodeAccessKeyMissing                               string = "ACCESS_KEY_MISSING"
	ErrCodeSecretKeyInvalid                               string = "SECRET_KEY_INVALID"
	ErrCodeSecretKeyMissing                               string = "SECRET_KEY_MISSING"
	ErrCodeDescriptionInvalid                             string = "DESCRIPTION_INVALID"
	ErrCodeObjectStorageForbidden                         string = "OBJECT_STORAGE_FORBIDDEN"
	ErrCodeObjectStorageNotFound                          string = "OBJECT_STORAGE_NOT_FOUND"
	ErrCodeNameInvalid                                    string = "NAME_INVALID"
	ErrCodeServiceExists                                  string = "SERVICE_EXISTS"
	ErrCodeValidationError                                string = "VALIDATION_ERROR"
	ErrCodeServiceError                                   string = "SERVICE_ERROR"
	ErrCodeServiceNotFound                                string = "SERVICE_NOT_FOUND"
	ErrCodeInvalidRequest                                 string = "INVALID_REQUEST"
	ErrCodeDBExists                                       string = "DB_EXISTS"
	ErrCodeDBNotFound                                     string = "DB_NOT_FOUND"
	ErrCodeUserNotFound                                   string = "USER_NOT_FOUND"
	ErrCodeConnectionPoolNotFound                         string = "CONNECTION_POOL_NOT_FOUND"
	ErrCodeResourceNotFound                               string = "RESOURCE_NOT_FOUND"
	ErrCodeTargetTypeInvalid                              string = "TARGET_TYPE_INVALID"
	ErrCodeTargetIdentifierInvalid                        string = "TARGET_IDENTIFIER_INVALID"
	ErrCodeUserInvalid                                    string = "USER_INVALID"
	ErrCodeInvalidOptions                                 string = "INVALID_OPTIONS"
	ErrCodeActionForbidden                                string = "ACTION_FORBIDDEN"
	ErrCodeAccountForbidden                               string = "ACCOUNT_FORBIDDEN"
	ErrCodeResourceAlreadyExists                          string = "RESOURCE_ALREADY_EXISTS"
	ErrCodeMethodNotAllowed                               string = "METHOD_NOT_ALLOWED"
	ErrCodeNotFound                                       string = "NOT_FOUND"
	ErrCodePeeringNotFound                                string = "PEERING_NOT_FOUND"
	ErrCodeLocalNetworkNoRouter                           string = "LOCAL_NETWORK_NO_ROUTER"
	ErrCodePeerNetworkNotFound                            string = "PEER_NETWORK_NOT_FOUND"
	ErrCodePeeringAccoundInvalid                          string = "PEERING_ACCOUNT_INVALID"
	ErrCodePeeringConflict                                string = "PEERING_CONFLICT"
	ErrPeeringNotDisabled                                 string = "PEERING_NOT_DISABLED"
	ErrCodeDuplicateResource                              string = "DUPLICATE_RESOURCE"
	ErrCodeServerIPLimitReached                           string = "SERVER_IP_LIMIT_REACHED"
	ErrCodeServerCoresLimitReached                        string = "SERVER_CORE_LIMIT_REACHED"
	ErrCodeServerMemoryLimitReached                       string = "SERVER_MEMORY_LIMIT_REACHED"
	ErrCodeMaxiOpsStorageLimitReached                     string = "MAXIOPS_STORAGE_LIMIT_REACHED"
)
