package types

import (
	"context"
	"path/filepath"
	"time"

	"github.com/hashicorp/nomad/plugins/device"
)

type License struct {
	// Unique ID that corresponds to a particular license instance. This ID is
	// only specific to this plugin, and doesn't hold relevance in other contexts.
	ID string

	// Vendor that provides the license. Must be all lowercase.
	Vendor string

	// Name of the license. Must be all lowercase, dashes converted to underscores.
	Feature string

	// One of:
	// - FREE
	// - RESERVED
	// - IN_USE
	Status string

	// The last timestamp when `Status` above changed
	LastUpdateTime time.Time

	// The name of the Nomad node that made the last status change. If status is
	// FREE, then this is null.
	UserNode *string

	// The name of the job for which the last status change was made. If status is
	// FREE, then this is null.
	UserProcess *string
}

func (l *License) MountInfo(root string) *device.Mount {
	return &device.Mount{
		HostPath: filepath.Join(root, l.Vendor, l.Feature, l.ID),
		TaskPath: filepath.Join("/tmp/license_handles", l.Vendor, l.Feature, l.ID),
		ReadOnly: true,
	}
}

type Reserver interface {
	Reserve(ctx context.Context, licenseIDs []string, node string) ([]*License, error)

	Use(ctx context.Context, licenseID string, node string, user string) error

	Free(ctx context.Context, licenseID string, node string) error
}

type Notifier interface {
	GetCurrent(ctx context.Context) ([]*License, error)

	Chan(ctx context.Context) <-chan struct{}
}
