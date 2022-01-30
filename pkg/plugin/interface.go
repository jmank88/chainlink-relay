package plugin

import (
	"github.com/gagliardetto/solana-go"
	"github.com/hashicorp/go-plugin"
	"github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var SolanaHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  1, // increment with changes
	MagicCookieKey:   "CL_RELAY_PLUGIN",
	MagicCookieValue: "9CEC5406C74D4A63BB4E7FE87709277B", // permanent GUID
}

type Solana interface {
	NewOCR2Provider(externalJobID uuid.UUID, spec SolanaSpec) (OCR2Provider, error)
}

type OCR2Provider interface {
	Start() error
	Close() error
	Ready() error
	Healthy() error
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}

type TransmissionSigner interface {
	Sign(msg []byte) ([]byte, error)
	PublicKey() solana.PublicKey
}
