package plugin

import (
	"github.com/pkg/errors"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	uuid "github.com/satori/go.uuid"
)

var _ plugin.Plugin = (*SolanaPlugin)(nil)

type SolanaPlugin struct {
	Impl Solana
}

func (p *SolanaPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &SolanaRPCServer{broker: broker, Impl: p.Impl}, nil
}

func (p *SolanaPlugin) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &SolanaRPC{broker: broker, client: client}, nil
}

var _ Solana = (*SolanaRPC)(nil)

type SolanaRPC struct {
	broker *plugin.MuxBroker
	client *rpc.Client
}

func (r *SolanaRPC) NewOCR2Provider(externalJobID uuid.UUID, spec SolanaSpec) (OCR2Provider, error) {
	var resp uint32
	err := r.client.Call("Plugin.NewOCR2Provider", NewOCR2ProviderRequest{
		ExternalJobID: externalJobID,
		Spec:          spec,
	}, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "plugin call failed")
	}
	conn, err := r.broker.Dial(resp)
	if err != nil {
		return nil, err
	}
	return &OCR2ProviderRPC{broker: r.broker, client: rpc.NewClient(conn)}, nil
}

type NewOCR2ProviderRequest struct {
	ExternalJobID uuid.UUID
	Spec          SolanaSpec
}

type SolanaRPCServer struct {
	broker *plugin.MuxBroker
	Impl   Solana
}

func (r *SolanaRPCServer) NewOCR2Provider(args NewOCR2ProviderRequest, ocr2ProviderID *uint32) error {
	relayer, err := r.Impl.NewOCR2Provider(args.ExternalJobID, args.Spec)
	if err != nil {
		return err
	}
	*ocr2ProviderID = r.broker.NextId()
	go r.broker.AcceptAndServe(*ocr2ProviderID, &OCR2ProviderRPCServer{broker: r.broker, Impl: relayer})

	return nil
}
