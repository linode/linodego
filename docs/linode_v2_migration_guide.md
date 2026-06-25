# Migration guide

## 1\. Update imports

linodego v2 uses a Go major-version module path.

Change all imports from:

```go
import "github.com/linode/linodego"
```

to:

```go
import "github.com/linode/linodego/v2"
```

Also update your module dependencies:

```shell
go get github.com/linode/linodego/v2
go mod tidy
```

## 2\. Handle NewClient errors

NewClient changed from returning only a Client to returning (Client, error).

### NewClient in linodego v1

```go
linodeClient := linodego.NewClient(oauth2Client)
```

### NewClient in linodego v2

```go
linodeClient, err := linodego.NewClient(oauth2Client)
if err != nil {
    return err
}
```

Every call site must now handle an initialization error.

## 3\. Replace removed deprecated APIs, aliases, and compatibility types

linodego v2 removes many deprecated fields, methods, aliases, and compatibility wrappers.

Confirmed removals include:

- MarkEventRead  
- ActionVolumeDelte  
- ActionCreateCardUpdated  
- AccountMaintenance.When  
- Domain.Group  
- LinodeKernel.XEN  
- MonitorAlertDefinition  
- CapabilityObjectStorageRegions  
- MutateInstance  
- LKEClusterDashboard  
- GetLKEClusterDashboard  
- ObjectStorageCluster  
- ListObjectStorageClusters  
- GetObjectStorageCluster  
- legacy paged response compatibility types  
- deprecated LKE cluster pool compatibility aliases and methods

If your code still references any of these, it must be updated before moving to v2.

Common replacements include:

- MarkEventRead → MarkEventsSeen  
- ActionVolumeDelte → ActionVolumeDelete  
- ActionCreateCardUpdated → ActionCreditCardUpdated  
- MutateInstance → UpgradeInstance  
- LKEClusterPool\* → LKENodePool\*

## 4\. Replace removed preview/temporary V2 compatibility APIs

Some preview or temporary V2\-suffixed APIs and types were removed in linodego v2.

**Migration rule:** in these areas, the old V2 forms were removed and the supported v2 surface now uses the original non-V2 names. In some cases this is effectively the same V2 behavior under the old name; in others, there are also related parameter, field-shape, or surrounding API changes.  
If you were already using the V2 form in v1, migrate to the corresponding supported non-V2 name in v2 and review any nearby type or behavior changes.

Affected removed V2 preview/temporary compatibility APIs and types include:

- GetInstanceTransferMonthlyV2  
- MonthlyInstanceTransferStatsV2  
- ObjectStorageBucketCertV2  
- UploadObjectStorageBucketCertV2  
- GetObjectStorageBucketCertV2  
- ObjectStorageObjectACLConfigV2  
- GetObjectStorageObjectACLConfigV2  
- UpdateObjectStorageObjectACLConfigV2  
- IPAddressUpdateOptionsV2  
- UpdateIPAddressV2

## 5\. Update methods that now take options structs

Many public methods that previously accepted individual primitive parameters for request bodies now take a single typed **options struct** instead.

In practice, request-body attributes are now generally grouped into dedicated ...Options structs and passed as one argument.

This affects a range of APIs, including:

- instance actions  
- disk operations  
- IP operations  
- snapshot creation  
- and similar endpoints

Migration approach:

- find compile errors  
- identify the corresponding ...Options type  
- move the request-body values into that struct  
- pass the struct instead of separate primitive values

This is mostly a mechanical migration.

### Notable exception

One breaking change goes in the opposite direction:

- CloneInstanceDisk(ctx, linodeID, diskID, opts InstanceDiskCloneOptions)  
- became  
- CloneInstanceDisk(ctx, linodeID, diskID)

So if you call CloneInstanceDisk, remove the now-deleted options argument.

## 6\. Migrate firewall APIs

The firewall API had substantial type cleanup.

### Type renames and splits

- FirewallRuleSet → FirewallRules  
- RuleSet → FirewallRuleSet

FirewallRule was split differently depending on which API you are using.

### If you use /firewalls APIs

Replace:

- FirewallRuleSet → FirewallRules  
- FirewallRule → FirewallRuleInbound / FirewallRuleOutbound

Also update create/update payload types to:

- FirewallRulesCreateOptions  
- FirewallRulesUpdateOptions

### If you use /firewalls/rulesets APIs

Replace:

- RuleSet → FirewallRuleSet  
- FirewallRule → FirewallRuleSetRule  
- RuleSetCreateOptions → FirewallRuleSetCreateOptions  
- RuleSetUpdateOptions → FirewallRuleSetUpdateOptions

Also note:

- FirewallRuleSetRule no longer includes description  
- FirewallRuleSetRule no longer includes ruleset

### Additional firewall migration note

Also update any code that assumes:

- FirewallDeviceEntity.Label is string instead of \*string

## 7\. Migrate Object Storage APIs to region-based usage

Object Storage APIs were standardized around **regions** instead of older cluster-oriented naming.

Key changes:

- stop using ObjectStorageBucket.Cluster  
- use region IDs in bucket/object/cert/ACL calls  
- replace ListObjectStorageBucketsInCluster with ListObjectStorageBucketsInRegion  
- remove usage of deleted ObjectStorageCluster APIs  
- GetObjectStorageBucketAccess now returns the newer access shape directly, so no separate V2 access type is needed  
- UpdateObjectStorageBucketAccess now uses PUT  
- ModifyObjectStorageBucketAccess was added as a separate POST

Also note that preview/temporary V2 Object Storage compatibility APIs were removed; use the supported non-V2 names in v2.

## 8\. Migrate IP update APIs

Confirmed changes include:

- IPAddressUpdateOptions is now the main IP update type  
- old split between original and V2 preview/temporary update option types removed  
- instance IP collection types changed from pointer element slices to value element slices in some responses

If you were using preview/temporary IP update V2 APIs in v1, switch to the supported non-V2 names in v2.

## 9\. Remove Resty-specific assumptions

The client internals were migrated from Resty to net/http.

Observable changes for advanced users:

- Request now maps to http.Request  
- Response now maps to http.Response  
- Logger is now a local interface, not Resty’s logger  
- OnBeforeRequest / OnAfterResponse now use \*http.Request / \*http.Response

If your code used Resty-specific fields or methods in callbacks, rewrite those hook implementations against net/http.

## 10\. Update region capability usages if needed

In regions.go, region capability constants changed from plain string constants to the custom RegionCapability type.

In many places this will work transparently, but code that relies on untyped string constant behavior may need small adjustments where exact types matter.

## 11\. Audit request/response type shape changes

A broad pattern in v2 is cleanup of pointer-heavy request/response types.

Common examples include:

- \[\]\*T → \[\]T  
- \*\[\]string → \[\]string  
- some string fields becoming pointers where nullable semantics are needed  
- split create/update types for request payloads

Review any code that depends on:

- nil checks  
- pointer identity  
- distinguishing “unset” from “empty”  
- mutation of shared slices

## 12\. Validate request payload behavior

A large number of JSON tags changed from omitempty to omitzero.

If your application depends on subtle update semantics such as “clear vs omit”, run integration tests against affected APIs before rollout.

## 13\. Review retry, logging, and error handling integrations

The retry, logging, and error handling layers changed along with the HTTP client rewrite.

Confirmed changes include:

- retry condition callback types now use \*http.Response  
- retry-after callback types now use \*http.Response  
- retry helper names are now exported/capitalized  
- logger behavior is no longer tied to Resty  
- error handling is now centered around http.Response instead of Resty responses

If your code directly integrates with retry helpers, logging hooks, or response-backed errors, review those call sites carefully.

## 14\. Suggested upgrade workflow

1. Update imports to /v2  
2. Run go get and go mod tidy  
3. Fix NewClient call sites  
4. Fix compile errors from removed deprecated APIs and aliases  
5. Fix usages of removed preview/temporary V2 compatibility APIs by switching to the supported non-V2 names  
6. Fix methods that now take typed options structs  
7. Fix firewall and object-storage type changes  
8. Fix request hook code that assumed Resty types  
9. Review any code depending on region capability constants as plain strings  
10. Run integration tests, especially around update payloads and object/firewall APIs