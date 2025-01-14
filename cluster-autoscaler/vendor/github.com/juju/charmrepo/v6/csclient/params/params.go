// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

// The params package holds types that are a part of the charm store's external
// contract - they will be marshalled (or unmarshalled) as JSON
// and delivered through the HTTP API.
package params // import "github.com/juju/charmrepo/v6/csclient/params"

import (
	"encoding/json"
	"time"

	"github.com/juju/charm/v8"
	"gopkg.in/macaroon.v2"
)

const (
	// ContentHashHeader specifies the header attribute
	// that will hold the content hash for archive GET responses.
	ContentHashHeader = "Content-Sha384"

	// EntityIdHeader specifies the header attribute that will hold the
	// id of the entity for archive GET responses.
	EntityIdHeader = "Entity-Id"
)

// Special user/group names.
const (
	Everyone = "everyone"
	Admin    = "admin"
)

// Channel is the name of a channel in which an entity may be published.
type Channel string

const (
	// EdgeChannel is the channel used for charms or bundles under development.
	EdgeChannel Channel = "edge"

	// BetaChannel is the channel used for beta charms or bundles.
	BetaChannel Channel = "beta"

	// CandidateChannel is the channel used for charms or bundles release
	// candidates.
	CandidateChannel Channel = "candidate"

	// StableChannel is the channel used for stable charms or bundles.
	StableChannel Channel = "stable"

	// UnpublishedChannel is the default channel to which charms are uploaded.
	UnpublishedChannel Channel = "unpublished"

	// NoChannel represents where no channel has been specifically requested.
	NoChannel Channel = ""

	// DevelopmentChannel is only defined for backward compatibility.
	DevelopmentChannel Channel = "development"
)

// OrderedChannels holds the list of valid channels in order of publishing
// status, most stable first.
var OrderedChannels = []Channel{
	StableChannel,
	CandidateChannel,
	BetaChannel,
	EdgeChannel,
	UnpublishedChannel,
}

// ValidChannels holds the set of all allowed channels for an entity.
var ValidChannels = func() map[Channel]bool {
	channels := make(map[Channel]bool, len(OrderedChannels))
	for _, ch := range OrderedChannels {
		channels[ch] = true
	}
	return channels
}()

// MetaAnyResponse holds the result of a meta/any request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaany
type MetaAnyResponse EntityResult

// ArchiveUploadResponse holds the result of a post or a put to /id/archive.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#post-idarchive
type ArchiveUploadResponse struct {
	Id            *charm.URL
	PromulgatedId *charm.URL `json:",omitempty"`
}

// Constants for the StatsUpdateRequest
type StatsUpdateType string

const (
	UpdateDownload StatsUpdateType = "download" // Accesses with non listed clients and web browsers.
	UpdateTraffic  StatsUpdateType = "traffic"  // Bots and unknown clients.
	UpdateDeploy   StatsUpdateType = "deploy"   // known clients like juju client.
)

// StatsUpdateRequest holds the parameters for a put to /stats/update.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#stats-update
type StatsUpdateRequest struct {
	Entries []StatsUpdateEntry
}

// StatsUpdateEntry holds an entry of the StatsUpdateRequest for a put to /stats/update.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#stats-update
type StatsUpdateEntry struct {
	Timestamp      time.Time       // Time when the update did happen.
	Type           StatsUpdateType // One of the constant Download, Traffic or Deploy.
	CharmReference *charm.URL      // The charm to be updated.
}

// ExpandedId holds a charm or bundle fully qualified id.
// A slice of ExpandedId is used as response for
// id/expand-id GET requests.
type ExpandedId struct {
	Id string
}

// ArchiveSizeResponse holds the result of an
// id/meta/archive-size GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaarchive-size
type ArchiveSizeResponse struct {
	Size int64
}

// HashResponse holds the result of id/meta/hash and id/meta/hash256 GET
// requests.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetahash
// and https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetahash256
type HashResponse struct {
	Sum string
}

// ManifestFile holds information about a charm or bundle file.
// A slice of ManifestFile is used as response for
// id/meta/manifest GET requests.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetamanifest
type ManifestFile struct {
	Name string
	Size int64
}

// ArchiveUploadTimeResponse holds the result of an id/meta/archive-upload-time
// GET request. See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaarchive-upload-time
type ArchiveUploadTimeResponse struct {
	UploadTime time.Time
}

// RelatedResponse holds the result of an id/meta/charm-related GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetacharm-related
type RelatedResponse struct {
	// Requires holds an entry for each interface provided by
	// the charm, containing all charms that require that interface.
	Requires map[string][]EntityResult `json:",omitempty"`

	// Provides holds an entry for each interface required by the
	// the charm, containing all charms that provide that interface.
	Provides map[string][]EntityResult `json:",omitempty"`
}

// RevisionInfoResponse holds the result of an id/meta/revision-info GET
// request. See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetarevision-info
type RevisionInfoResponse struct {
	Revisions []*charm.URL
}

// SupportedSeries holds the result of an id/meta/supported-series GET
// request. See See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetasupported-series
type SupportedSeriesResponse struct {
	SupportedSeries []string
}

// BundleCount holds the result of an id/meta/bundle-unit-count
// or bundle-machine-count GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetabundle-unit-count
// and https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetabundle-machine-count
type BundleCount struct {
	Count int
}

// TagsResponse holds the result of an id/meta/tags GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetatags
type TagsResponse struct {
	Tags []string
}

// Published holds the result of a changes/published GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-changespublished
type Published struct {
	Id          *charm.URL
	PublishTime time.Time
}

// DebugStatus holds the result of the status checks.
// This is defined for backward compatibility: new clients should use
// debugstatus.CheckResult directly.
type DebugStatus struct {
	// Name is the human readable name for the check.
	Name string

	// Value is the check result.
	Value string

	// Passed reports whether the check passed.
	Passed bool

	// Duration holds the duration that the
	// status check took to run.
	Duration time.Duration
}

// EntityResult holds a the resolved entity ID along with any requested metadata.
type EntityResult struct {
	Id *charm.URL
	// Meta holds at most one entry for each meta value
	// specified in the include flags, holding the
	// data that would be returned by reading /meta/meta?id=id.
	// Metadata not relevant to a particular result will not
	// be included.
	Meta map[string]interface{} `json:",omitempty"`
}

// SearchResponse holds the response from a search operation.
type SearchResponse struct {
	SearchTime time.Duration
	Total      int
	Results    []EntityResult
}

// ListResponse holds the response from a list operation.
type ListResponse struct {
	Results []EntityResult
}

// IdUserResponse holds the result of an id/meta/id-user GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaid-user
type IdUserResponse struct {
	User string
}

// IdSeriesResponse holds the result of an id/meta/id-series GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaid-series
type IdSeriesResponse struct {
	Series string
}

// IdNameResponse holds the result of an id/meta/id-name GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaid-name
type IdNameResponse struct {
	Name string
}

// IdRevisionResponse holds the result of an id/meta/id-revision GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaid-revision
type IdRevisionResponse struct {
	Revision int
}

// IdResponse holds the result of an id/meta/id GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaid
type IdResponse struct {
	Id       *charm.URL
	User     string `json:",omitempty"`
	Series   string `json:",omitempty"`
	Name     string
	Revision int
}

// AllPermsResponse holds the resource of an id/allperms GET
// request.
type AllPermsResponse struct {
	Perms map[Channel]PermResponse
}

// PermResponse holds the result of an id/meta/perm GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetaperm
type PermResponse struct {
	Read  []string
	Write []string
}

// PermRequest holds the request of an id/meta/perm PUT request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#put-idmetaperm
type PermRequest struct {
	Read  []string
	Write []string
}

// PromulgatedResponse holds the result of an id/meta/promulgated GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetapromulgated
type PromulgatedResponse struct {
	Promulgated bool
}

// CanIngestResponse holds the result of an id/meta/can-ingest GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetacan-ingest
type CanIngestResponse struct {
	CanIngest bool
}

// CanWriteResponse holds the result of an id/meta/can-write GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-idmetacan-write
type CanWriteResponse struct {
	CanWrite bool
}

// PromulgateRequest holds the request of an id/promulgate PUT request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#put-idpromulgate
type PromulgateRequest struct {
	Promulgated bool
}

// PublishRequest holds the request of an id/publish PUT request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#put-idpublish
type PublishRequest struct {
	Channels []Channel
	// Resources defines the resource revisions to use for the charm.
	// Each resource in the charm's metadata.yaml (if any) must have its
	// name mapped to a revision. That revision must be one of the
	// existing revisions for that resource.
	Resources map[string]int `json:",omitempty"`
}

// PublishResponse holds the result of an id/publish PUT request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#put-idpublish
type PublishResponse struct {
	Id            *charm.URL
	PromulgatedId *charm.URL `json:",omitempty"`
}

// PublishedResponse holds the result of an id/meta/published GET request.
type PublishedResponse struct {
	// Channels holds an entry for each channel that the
	// entity has been published to.
	Info []PublishedInfo
}

// PublishedInfo holds information on a channel that an entity
// has been published to.
type PublishedInfo struct {
	// Channel holds the value of the channel that
	// the entity has been published to.
	// This will never be "unpublished" as entities
	// cannot be published to that channel.
	Channel Channel

	// Current holds whether the entity is the most
	// recently published member of the channel.
	Current bool
}

// WhoAmIResponse holds the result of a whoami GET request.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#whoami
type WhoAmIResponse struct {
	User   string
	Groups []string
}

// Resource describes a resource in the charm store.
type Resource struct {
	// Name identifies the resource.
	Name string

	// Type is the name of the resource type.
	Type string

	// Path is where the resource will be stored.
	Path string

	// Description contains user-facing info about the resource.
	Description string `json:",omitempty"`

	// Revision is the revision, if applicable.
	Revision int

	// Fingerprint is the SHA-384 checksum for the resource blob.
	Fingerprint []byte

	// Size is the size of the resource, in bytes.
	Size int64
}

// ResourceUploadResponse holds the result of a post or a put to /id/resources/name.
type ResourceUploadResponse struct {
	Revision int
}

// DockerResourceUploadRequest holds the body of a POST to /:id/resources/:name
// when the resource is a docker image.
type DockerResourceUploadRequest struct {
	// ImageName holds the image name when it's an external image not
	// contained within the charm store's registry. If this is empty, the
	// image should have been uploaded to the charm store's registry.
	ImageName string
	// Digest holds the digest of the image, in the form "sha256:hexbytes".
	Digest string
}

// DockerInfoResponse holds the result of a get of /:id/resources/:name/docker-info
type DockerInfoResponse struct {
	// ImageName holds the image name (including host) of the resource in the docker registry.
	ImageName string

	// Username holds the username to use in the docker auth information.
	// (see https://docs.docker.com/registry/spec/auth/token/#requesting-a-token).
	Username string

	// Password holds the password to use in the docker auth information.
	Password string
}

// CharmRevision holds the revision number of a charm and any error
// encountered in retrieving it.
type CharmRevision struct {
	Revision int
	Sha256   string
	Err      error
}

const (
	// BzrDigestKey is the extra-info key used to store the Bazaar digest
	BzrDigestKey = "bzr-digest"

	// LegacyDownloadStats is the extra-info key used to store the legacy
	// download counts, and to retrieve them when
	// charmstore.LegacyDownloadCountsEnabled is set to true.
	// TODO (frankban): remove this constant when removing the legacy counts
	// logic.
	LegacyDownloadStats = "legacy-download-stats"
)

// Log holds the representation of a log message.
// This is used by clients to store log events in the charm store.
type Log struct {
	// Data holds the log message as a JSON-encoded value.
	Data *json.RawMessage

	// Level holds the log level as a string.
	Level LogLevel

	// Type holds the log type as a string.
	Type LogType

	// URLs holds a slice of entity URLs associated with the log message.
	URLs []*charm.URL `json:",omitempty"`
}

// LogResponse represents a single log message and is used in the responses
// to /log GET requests.
// See https://github.com/juju/charmstore/blob/v5-unstable/docs/API.md#get-log
type LogResponse struct {
	// Data holds the log message as a JSON-encoded value.
	Data json.RawMessage

	// Level holds the log level as a string.
	Level LogLevel

	// Type holds the log type as a string.
	Type LogType

	// URLs holds a slice of entity URLs associated with the log message.
	URLs []*charm.URL `json:",omitempty"`

	// Time holds the time of the log.
	Time time.Time
}

// LogLevel defines log levels (e.g. "info" or "error") to be used in log
// requests and responses.
type LogLevel string

const (
	InfoLevel    LogLevel = "info"
	WarningLevel LogLevel = "warning"
	ErrorLevel   LogLevel = "error"
)

// LogType defines log types (e.g. "ingestion") to be used in log requests and
// responses.
type LogType string

const (
	IngestionType        LogType = "ingestion"
	LegacyStatisticsType LogType = "legacyStatistics"

	IngestionStart    = "ingestion started"
	IngestionComplete = "ingestion completed"

	LegacyStatisticsImportStart    = "legacy statistics import started"
	LegacyStatisticsImportComplete = "legacy statistics import completed"
)

// SetAuthCookie holds the parameters used to make a set-auth-cookie request
// to the charm store.
type SetAuthCookie struct {
	// Macaroons holds a slice of macaroons.
	Macaroons macaroon.Slice
}

// NewUploadResponse holds the response from a POST request to /upload.
// TODO remove this when the charmstore code no longer requires it.
type NewUploadResponse UploadInfoResponse

// Parts holds a list of all the parts that are required by a multipart
// upload, as required by a PUT request to /upload/$upload-id.
type Parts struct {
	Parts []Part
}

// Part represents one part of a multipart blob.
// When a set of parts is returned from an upload
// query GET, those with zero sizes should
// be considered non-existent.
type Part struct {
	// Hash holds the SHA384 hash of the part.
	Hash string `json:",omitempty"`
	// Size holds the size of the part.
	Size int64 `json:",omitempty"`
	// Offset holds the offset of the part from the start
	// of the file.
	Offset int64 `json:",omitempty"`
	// Complete holds whether the part has been
	// successfully uploaded.
	Complete bool `json:",omitempty"`
}

func (p Part) Valid() bool {
	return p.Size > 0
}

// FinishUploadResponse holds the response to a put /upload/upload-id/part-number request.
type FinishUploadResponse struct {
	// Hash holds the SHA384 hash of the complete blob. (hex-encoded)
	Hash string
}

// UploadInfoResponse holds the response to a get /upload/upload-id request.
type UploadInfoResponse struct {
	// UploadId holds the id of the upload.
	UploadId string

	// Parts holds all the known parts of the upload.
	// Parts that haven't been uploaded yet will have nil
	// elements.
	Parts Parts

	// Expires holds when the upload will expire.
	Expires time.Time

	// MinPartSize holds the minimum size of a part that may
	// be uploaded (not including the last part).
	MinPartSize int64

	// MaxPartSize holds the maximum size of a part that may
	// be uploaded.
	MaxPartSize int64

	// MaxParts holds the maximum number of parts.
	MaxParts int
}
