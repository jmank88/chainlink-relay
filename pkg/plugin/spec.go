package plugin

import "github.com/gagliardetto/solana-go"

//TODO convert to proto and generate
type SolanaSpec struct {
	ID          int32
	IsBootstrap bool

	// network data
	NodeEndpointHTTP string

	// on-chain program + 2x state accounts (state + transmissions) + store program
	ProgramID       solana.PublicKey
	StateID         solana.PublicKey
	StoreProgramID  solana.PublicKey
	TransmissionsID solana.PublicKey

	// transaction + state parameters [optional]
	UsePreflight bool
	Commitment   string

	// polling configuration [optional]
	PollingInterval   string
	PollingCtxTimeout string
	StaleTimeout      string

	TransmissionSigner TransmissionSigner
}
