package mid

import "time"

type NatsConfig struct {
	// The URL of the Connect server.
	Url        string
	Name       string
	StreamName string // 持久化名称

	ReconnectWait int // 秒
	Timeout       int // 秒
	// MaxReconnect sets the number of reconnect attempts that will be
	// tried before giving up. If negative, then it will never give up
	// trying to reconnect.
	// Defaults to 60.
	MaxReconnects int
	// The username to use when authenticating.
	Username string
	// The password to use when authenticating.
	Password string

	// NoRandomize configures whether we will randomize the
	// server pool.
	NoRandomize bool

	// NoEcho configures whether the server will echo back messages
	// that are sent on this connection if we also have matching subscriptions.
	// Note this is supported on servers >= version 1.2. Proto 1 or greater.
	NoEcho bool

	// Verbose signals the server to send an OK ack for commands
	// successfully processed by the server.
	Verbose bool

	// Pedantic signals the server whether it should be doing further
	// validation of subjects.
	Pedantic bool

	// Secure enables TLS secure connections that skip server
	// verification by default. NOT RECOMMENDED.
	Secure bool

	// AllowReconnect enables reconnection logic to be used when we
	// encounter a disconnect from the current server.
	AllowReconnect bool

	// Token sets the token to be used when connecting to a server.
	Token        string
	StreamConfig StreamConfig
}

type StreamConfig struct {

	// Name is the name of the stream. It is required and must be unique
	// across the JetStream account.
	//
	// Name Names cannot contain whitespace, ., *, >, path separators
	// (forward or backwards slash), and non-printable characters.
	Name string

	// Description is an optional description of the stream.
	Description string

	// Subjects is a list of subjects that the stream is listening on.
	// Wildcards are supported. Subjects cannot be set if the stream is
	// created as a mirror.
	Subjects []string

	// Retention defines the message retention policy for the stream.
	// Defaults to LimitsPolicy.
	Retention int

	// MaxConsumers specifies the maximum number of consumers allowed for
	// the stream.
	MaxConsumers int

	// MaxMsgs is the maximum number of messages the stream will store.
	// After reaching the limit, stream adheres to the discard policy.
	// If not set, server default is -1 (unlimited).
	MaxMsgs int64

	// MaxBytes is the maximum total size of messages the stream will store.
	// After reaching the limit, stream adheres to the discard policy.
	// If not set, server default is -1 (unlimited).
	MaxBytes int64

	// Discard defines the policy for handling messages when the stream
	// reaches its limits in terms of number of messages or total bytes.
	Discard int

	// DiscardNewPerSubject is a flag to enable discarding new messages per
	// subject when limits are reached. Requires DiscardPolicy to be
	// DiscardNew and the MaxMsgsPerSubject to be set.
	DiscardNewPerSubject bool

	// MaxAge is the maximum age of messages that the stream will retain.
	MaxAge time.Duration

	// MaxMsgsPerSubject is the maximum number of messages per subject that
	// the stream will retain.
	MaxMsgsPerSubject int64

	// MaxMsgSize is the maximum size of any single message in the stream.
	MaxMsgSize int32

	// Storage specifies the type of storage backend used for the stream
	// (file or memory).
	Storage int

	// Replicas is the number of stream replicas in clustered JetStream.
	// Defaults to 1, maximum is 5.
	Replicas int

	// NoAck is a flag to disable acknowledging messages received by this
	// stream.
	//
	// If set to true, publish methods from the JetStream client will not
	// work as expected, since they rely on acknowledgements. Core NATS
	// publish methods should be used instead. Note that this will make
	// message delivery less reliable.
	NoAck bool

	// Sealed streams do not allow messages to be published or deleted via limits or API,
	// sealed streams can not be unsealed via configuration update. Can only
	// be set on already created streams via the Update API.
	Sealed bool
	// DenyDelete restricts the ability to delete messages from a stream via
	// the API. Defaults to false.
	DenyDelete bool

	// DenyPurge restricts the ability to purge messages from a stream via
	// the API. Defaults to false.
	DenyPurge bool

	// AllowRollup allows the use of the Nats-Rollup header to replace all
	// contents of a stream, or subject in a stream, with a single new
	// message.
	AllowRollup bool

	// FirstSeq is the initial sequence number of the first message in the
	// stream.
	FirstSeq uint64

	// AllowDirect enables direct access to individual messages using direct
	// get API. Defaults to false.
	AllowDirect bool

	// MirrorDirect enables direct access to individual messages from the
	// origin stream using direct get API. Defaults to false.
	MirrorDirect bool

	// Metadata is a set of application-defined key-value pairs for
	// associating metadata on the stream. This feature requires nats-server
	// v2.10.0 or later.
	Metadata map[string]string

	// Template identifies the template that manages the Stream. Deprecated:
	// This feature is no longer supported.
	Template    string
	Compression uint
}
