# Event and Status Polling

Many Linode API calls finish on the server after the HTTP response returns. Creating an instance, booting it, resizing a volume, and similar work is asynchronous.

linodego helps you wait for that work in two ways:

1. **Status polling** — keep checking a resource until its status looks right (for example, an instance is `running`).
2. **Event polling** — watch account events until a specific action finishes (for example, `linode_boot` reaches `finished`).

Every wait helper uses the client's poll interval and stops when:

- the condition is met
- an error is returned
- the `context.Context` is canceled or times out

## Quick start: which helper should I use?

| What you want | Use this |
| --- | --- |
| A resource reached a known status | A `WaitFor*Status` helper |
| A specific action finished on an entity | `NewEventPoller` + `WaitForFinished` |
| A create finished, but you did not know the ID yet | `NewEventPollerWithoutEntity` |
| An action on a nested resource (for example, a disk on an instance) | `NewEventPollerWithSecondary` |
| No in-progress events left on a resource | `WaitForResourceFree` |
| A timestamp-based event wait (lower-level alternative) | `WaitForEventFinished` |

**Rule of thumb**

- Use **status waits** when you care whether the resource looks ready.
- Use **event waits** when you care whether a particular action completed.
- Prefer `EventPoller` over `WaitForEventFinished` when you can create the poller before the mutating API call.

## Configuration

### How often does it poll?

By default, the client polls every 3 seconds (`APISecondsPerPoll`).

```go
client.SetPollDelay(5 * time.Second)
delay := client.GetPollDelay()
```

`SetPollDelay` controls how often wait helpers and event pollers check for progress.

By default, request retries also start with a 3 second minimum wait (`SetRetryWaitTime`), matching the poll delay. Those are separate settings: changing the poll delay does not automatically change retry timing.

Shorter poll delays notice completion sooner but create more API traffic. Longer delays are quieter but slower.

### How do timeouts work?

Always pass a deadline-aware context into the wait call itself (`WaitForFinished`, `WaitForInstanceStatus`, and so on). Without a deadline, a wait can block forever if the expected status or event never appears.

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
defer cancel()

instance, err := client.WaitForInstanceStatus(ctx, instanceID, linodego.InstanceRunning)
if err != nil {
	log.Fatal(err)
}
```

When the context expires, wait helpers return an error that wraps `ctx.Err()`.

For `NewEventPoller` and `NewEventPollerWithSecondary`, the create-poller call also takes a context because those helpers list existing events first. For `NewEventPollerWithoutEntity`, only the later wait call needs the deadline context.

## Status polling

Status helpers repeatedly fetch a resource and return once its status matches what you asked for.

### Available helpers

| Method | Waits for |
| --- | --- |
| `WaitForInstanceStatus` | Instance status |
| `WaitForInstanceDiskStatus` | Instance disk status |
| `WaitForVolumeStatus` | Volume status |
| `WaitForVolumeLinodeID` | Volume attach or detach (`LinodeID`) |
| `WaitForVolumeIOReadyStatus` | Volume `IOReady` |
| `WaitForSnapshotStatus` | Instance snapshot status |
| `WaitForImageStatus` | Image status |
| `WaitForImageRegionStatus` | Image replica status in a region |
| `WaitForLKEClusterStatus` | LKE cluster status |
| `WaitForLKEClusterConditions` | Custom LKE conditions |
| `WaitForDatabaseStatus` | Managed database status |
| `WaitForAlertDefinitionStatus` | Monitor alert definition status |

### Example: wait for an instance to become running

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
defer cancel()

instance, err := client.CreateInstance(ctx, linodego.InstanceCreateOptions{
	Region:   "us-east",
	Type:     "g6-nanode-1",
	Label:    "polling-example",
	Image:    "linode/ubuntu22.04",
	RootPass: "replace-with-a-secure-password",
})
if err != nil {
	log.Fatal(err)
}

instance, err = client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("instance %d is %s\n", instance.ID, instance.Status)
```

### Example: wait for a volume to become active

```go
volume, err := client.WaitForVolumeStatus(ctx, volumeID, linodego.VolumeActive)
if err != nil {
	log.Fatal(err)
}
```

### Example: wait for a volume to attach or detach

```go
// Wait until the volume is attached to this instance.
volume, err := client.WaitForVolumeLinodeID(ctx, volumeID, &instanceID)
if err != nil {
	log.Fatal(err)
}

// Wait until the volume is detached.
volume, err = client.WaitForVolumeLinodeID(ctx, volumeID, nil)
if err != nil {
	log.Fatal(err)
}
```

## Event polling with EventPoller

Account events describe actions on entities. `EventPoller` watches for one entity and one action.

Create the poller **before** you trigger the operation so the helper can snapshot existing events (record a baseline of event IDs to ignore). That way old events are skipped, and you only wait for the new one.

This is the best option when:

- you care about a specific action completing
- the same entity may have several similar events over time

### Existing entity

```go
// 1. Create the poller first so current events are recorded and ignored.
poller, err := client.NewEventPoller(
	ctx,
	instance.ID,
	linodego.EntityLinode,
	linodego.ActionLinodeBoot,
)
if err != nil {
	log.Fatal(err)
}

// 2. Trigger the operation.
if err := client.BootInstance(ctx, instance.ID, linodego.InstanceBootOptions{}); err != nil {
	log.Fatal(err)
}

// 3. Wait for the matching event to finish.
event, err := poller.WaitForFinished(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("boot event %d finished with status %s\n", event.ID, event.Status)
```

What `WaitForFinished` does:

1. waits for the next unseen matching event
2. polls that event until its status is `finished`
3. returns `nil` and an error if the event status becomes `failed`

If you need the failed event object itself, use `WaitForEventFinished` instead. That helper returns both the event and an error on failure.

If you only need the next matching event, and not necessarily a finished one, call `WaitForLatestUnknownEvent`:

```go
event, err := poller.WaitForLatestUnknownEvent(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("observed event %d with status %s\n", event.ID, event.Status)
```

### Create operations (ID not known yet)

When you create a resource, you usually do not know its ID until the create call returns. Use `NewEventPollerWithoutEntity`, set `EntityID` right after create, then wait.

Unlike `NewEventPoller` and `NewEventPollerWithSecondary`, this helper does **not** snapshot existing events up front. `previousEvents` starts empty. You can create the poller at any time, but you must set `EntityID` before calling `WaitForFinished`.

```go
poller, err := client.NewEventPollerWithoutEntity(
	linodego.EntityLinode,
	linodego.ActionLinodeCreate,
)
if err != nil {
	log.Fatal(err)
}

instance, err := client.CreateInstance(ctx, linodego.InstanceCreateOptions{
	Region: "us-east",
	Type:   "g6-nanode-1",
	Label:  "create-poll-example",
	Booted: linodego.Pointer(false),
})
if err != nil {
	log.Fatal(err)
}

// Set this before waiting. Even if create finished quickly, the poller can
// still match the create event because previousEvents starts empty.
poller.EntityID = instance.ID

event, err := poller.WaitForFinished(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("create event %d finished\n", event.ID)
```

### Secondary entities

Some events have both a primary entity and a secondary entity. Deleting a disk on an instance is a common example: the instance is primary, and the disk is secondary.

`NewEventPollerWithSecondary` takes the secondary ID as an `int`, which fits nested resources such as disks.

```go
poller, err := client.NewEventPollerWithSecondary(
	ctx,
	instance.ID, // primary entity
	linodego.EntityLinode,
	disk.ID, // secondary entity
	linodego.ActionDiskDelete,
)
if err != nil {
	log.Fatal(err)
}

if err := client.DeleteInstanceDisk(ctx, instance.ID, disk.ID); err != nil {
	log.Fatal(err)
}

event, err := poller.WaitForFinished(ctx)
if err != nil {
	log.Fatal(err)
}
```

## Other event helpers

### WaitForEventFinished

`WaitForEventFinished` is a lower-level helper. It finds events matching an entity and action that were created at or after a given timestamp, then waits until one reaches `finished` status.

```go
event, err := client.WaitForEventFinished(
	ctx,
	instance.ID,
	linodego.EntityLinode,
	linodego.ActionLinodeCreate,
	*instance.Created,
)
if err != nil {
	// On failure, event may still be non-nil.
	log.Fatal(err)
}
```

Notes:

- Prefer `EventPoller` when you can create the poller before the mutating call. It avoids timestamp edge cases and ignores events that already exist.
- If the matched event fails, this helper returns both the event and an error.
- Entity filtering is optimized for disk, database, linode, domain, and nodebalancer entities. Other entity types may be less precise.

### WaitForResourceFree

Use this when you want a resource to settle before starting another long-running operation. It waits until the entity has no events in `started` or `scheduled` status.

```go
if err := client.WaitForResourceFree(ctx, linodego.EntityLinode, instance.ID); err != nil {
	log.Fatal(err)
}
```

## Practical tips

1. **Put the deadline on the wait call.** Status waits, `WaitForFinished`, and `WaitForEventFinished` all need a context that can expire.
2. **Create snapshotting pollers before the mutating call.** `NewEventPoller` and `NewEventPollerWithSecondary` snapshot existing events first. `NewEventPollerWithoutEntity` does not.
3. **Pick the wait that matches your goal.** Status waits answer "is it ready?" Event waits answer "did this action finish?"
4. **Tune the poll delay carefully.** Faster polling is more responsive; slower polling is gentler on the API.
5. **Handle failed events.** `EventPoller.WaitForFinished` returns `nil, error` on failure. `WaitForEventFinished` returns the failed event along with an error.

## Reference

- Implementation: [waitfor.go](../waitfor.go)
- Event, entity, and action constants: [account_events.go](../account_events.go)
- Common event statuses: `EventScheduled`, `EventStarted`, `EventFinished`, `EventFailed`, `EventNotification`, `EventCanceled`
