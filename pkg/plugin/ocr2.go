package plugin

import (
	"context"
	"math/big"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type OCR2ProviderRPC struct {
	broker *plugin.MuxBroker
	client *rpc.Client
}

func (o *OCR2ProviderRPC) Start() error {
	var resp error
	err := o.client.Call("Plugin.Start", new(interface{}), &resp)
	if err != nil {
		return errors.Wrap(err, "plugin call failed")
	}
	return resp
}

func (o *OCR2ProviderRPC) Close() error {
	var resp error
	err := o.client.Call("Plugin.Close", new(interface{}), &resp)
	if err != nil {
		return errors.Wrap(err, "plugin call failed")
	}
	return resp
}

func (o *OCR2ProviderRPC) Ready() error {
	var resp error
	err := o.client.Call("Plugin.Ready", new(interface{}), &resp)
	if err != nil {
		return errors.Wrap(err, "plugin call failed")
	}
	return resp
}

func (o *OCR2ProviderRPC) Healthy() error {
	var resp error
	err := o.client.Call("Plugin.Healthy", new(interface{}), &resp)
	if err != nil {
		return errors.Wrap(err, "plugin call failed")
	}
	return resp
}

func (o *OCR2ProviderRPC) ContractTransmitter() types.ContractTransmitter {
	var resp uint32
	MustCall(o.client.Call("Plugin.ContractTransmitter", new(interface{}), &resp))
	conn := MustDial(o.broker.Dial(resp))
	return &ContractTransmitterRPC{client: rpc.NewClient(conn)}
}

func (o *OCR2ProviderRPC) ContractConfigTracker() types.ContractConfigTracker {
	var resp uint32
	MustCall(o.client.Call("Plugin.ContractConfigTracker", new(interface{}), &resp))
	conn := MustDial(o.broker.Dial(resp))
	return &ContractConfigTrackerRPC{client: rpc.NewClient(conn)}
}

func (o *OCR2ProviderRPC) OffchainConfigDigester() types.OffchainConfigDigester {
	var resp uint32
	MustCall(o.client.Call("Plugin.OffchainConfigDigester", new(interface{}), &resp))
	conn := MustDial(o.broker.Dial(resp))
	return &OffchainConfigDigesterRPC{client: rpc.NewClient(conn)}
}

func (o *OCR2ProviderRPC) ReportCodec() median.ReportCodec {
	var resp uint32
	MustCall(o.client.Call("Plugin.ReportCodec", new(interface{}), &resp))
	conn := MustDial(o.broker.Dial(resp))
	return &ReportCodecRPC{client: rpc.NewClient(conn)}
}

func (o *OCR2ProviderRPC) MedianContract() median.MedianContract {
	var resp uint32
	MustCall(o.client.Call("Plugin.MedianContract", new(interface{}), &resp))
	conn := MustDial(o.broker.Dial(resp))
	return &MedianContractRPC{client: rpc.NewClient(conn)}
}

type OCR2ProviderRPCServer struct {
	broker *plugin.MuxBroker
	Impl   OCR2Provider
}

func (o *OCR2ProviderRPCServer) Start(args interface{}, resp *error) error {
	*resp = o.Impl.Start()
	return nil
}

func (o *OCR2ProviderRPCServer) Close(args interface{}, resp *error) error {
	*resp = o.Impl.Close()
	return nil
}

func (o *OCR2ProviderRPCServer) Ready(args interface{}, resp *error) error {
	*resp = o.Impl.Ready()
	return nil
}

func (o *OCR2ProviderRPCServer) Healthy(args interface{}, resp *error) error {
	*resp = o.Impl.Healthy()
	return nil
}

func (o *OCR2ProviderRPCServer) ContractTransmitter(args interface{}, resp *uint32) error {
	*resp = o.broker.NextId()
	go o.broker.AcceptAndServe(*resp, &ContractTransmitterRPCServer{Impl: o.Impl.ContractTransmitter()})
	return nil
}

func (o *OCR2ProviderRPCServer) ContractConfigTracker(args interface{}, resp *uint32) error {
	*resp = o.broker.NextId()
	go o.broker.AcceptAndServe(*resp, &ContractConfigTrackerRPCServer{o.Impl.ContractConfigTracker()})
	return nil
}

func (o *OCR2ProviderRPCServer) OffchainConfigDigester(args interface{}, resp *uint32) error {
	*resp = o.broker.NextId()
	go o.broker.AcceptAndServe(*resp, &OffchainConfigDigesterRPCServer{Impl: o.Impl.OffchainConfigDigester()})
	return nil
}

func (o *OCR2ProviderRPCServer) ReportCodec(args interface{}, resp *uint32) error {
	*resp = o.broker.NextId()
	go o.broker.AcceptAndServe(*resp, &ReportCodecRPCServer{Impl: o.Impl.ReportCodec()})
	return nil
}

func (o *OCR2ProviderRPCServer) MedianContract(args interface{}, resp *uint32) error {
	*resp = o.broker.NextId()
	go o.broker.AcceptAndServe(*resp, &MedianContractRPCServer{Impl: o.Impl.MedianContract()})
	return nil
}

type ContractTransmitterRPC struct {
	client *rpc.Client
}

func (c *ContractTransmitterRPC) Transmit(ctx context.Context, context types.ReportContext, report types.Report, signatures []types.AttributedOnchainSignature) error {
	var resp error
	err := c.client.Call("Plugin.Transmit", TransmitRequest{
		Context: context, Report: report, Signatures: signatures,
	}, &resp)
	if err != nil {
		return errors.Wrap(err, "plugin Call failed")
	}
	return resp
}

func (c *ContractTransmitterRPC) LatestConfigDigestAndEpoch(ctx context.Context) (configDigest types.ConfigDigest, epoch uint32, err error) {
	var resp LatestConfigDigestAndEpochResponse
	err = c.client.Call("Plugin.LatestConfigDigestAndEpoch", new(interface{}), &resp)
	if err != nil {
		err = errors.Wrap(err, "plugin Call failed")
		return
	}
	return resp.ConfigDigest, resp.Epoch, nil
}

func (c *ContractTransmitterRPC) FromAccount() (resp types.Account) {
	MustCall(c.client.Call("Plugin.Transmit", new(interface{}), &resp))
	return
}

type TransmitRequest struct {
	Context    types.ReportContext
	Report     types.Report
	Signatures []types.AttributedOnchainSignature
}

type LatestConfigDigestAndEpochResponse struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint32
}

type ContractTransmitterRPCServer struct {
	Impl types.ContractTransmitter
}

func (c *ContractTransmitterRPCServer) Transmit(args TransmitRequest, resp *error) error {
	*resp = c.Impl.Transmit(context.Background(), args.Context, args.Report, args.Signatures)
	return nil
}

func (c *ContractTransmitterRPCServer) LatestConfigDigestAndEpoch(args interface{}, resp *LatestConfigDigestAndEpochResponse) (err error) {
	resp.ConfigDigest, resp.Epoch, err = c.Impl.LatestConfigDigestAndEpoch(context.Background())
	return
}

func (c *ContractTransmitterRPCServer) FromAccount(args interface{}, resp *types.Account) error {
	*resp = c.Impl.FromAccount()
	return nil
}

type ContractConfigTrackerRPC struct {
	client *rpc.Client
}

func (c *ContractConfigTrackerRPC) Notify() <-chan struct{} {
	return nil
}

func (c *ContractConfigTrackerRPC) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	var resp LatestConfigDetailsResponse
	err = c.client.Call("Plugin.LatestConfigDetails", new(interface{}), &resp)
	if err != nil {
		err = errors.Wrap(err, "plugin Call failed")
		return
	}
	return resp.ChangedInBlock, resp.ConfigDigest, nil
}

func (c *ContractConfigTrackerRPC) LatestConfig(ctx context.Context, changedInBlock uint64) (resp types.ContractConfig, err error) {
	err = c.client.Call("Plugin.LatestConfig", new(interface{}), &resp)
	if err != nil {
		err = errors.Wrap(err, "plugin Call failed")
	}
	return
}

func (c *ContractConfigTrackerRPC) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	err = c.client.Call("Plugin.LatestBlockHeight", new(interface{}), &blockHeight)
	if err != nil {
		err = errors.Wrap(err, "plugin Call failed")
	}
	return
}

type LatestConfigDetailsResponse struct {
	ChangedInBlock uint64
	ConfigDigest   types.ConfigDigest
}

type ContractConfigTrackerRPCServer struct {
	Impl types.ContractConfigTracker
}

func (c *ContractConfigTrackerRPCServer) LatestConfigDetails(args interface{}, resp *LatestConfigDetailsResponse) (err error) {
	resp.ChangedInBlock, resp.ConfigDigest, err = c.Impl.LatestConfigDetails(context.Background())
	return
}

func (c *ContractConfigTrackerRPCServer) LatestConfig(changedInBlock uint64, resp *types.ContractConfig) (err error) {
	*resp, err = c.Impl.LatestConfig(context.Background(), changedInBlock)
	return
}

func (c *ContractConfigTrackerRPCServer) LatestBlockHeight(args interface{}, blockHeight *uint64) (err error) {
	*blockHeight, err = c.Impl.LatestBlockHeight(context.Background())
	return
}

type OffchainConfigDigesterRPC struct {
	client *rpc.Client
}

func (o *OffchainConfigDigesterRPC) ConfigDigest(config types.ContractConfig) (resp types.ConfigDigest, err error) {
	err = o.client.Call("Plugin.ConfigDigest", config, &resp)
	return
}

func (o *OffchainConfigDigesterRPC) ConfigDigestPrefix() (resp types.ConfigDigestPrefix) {
	MustCall(o.client.Call("Plugin.ConfigDigestPrefix", new(interface{}), &resp))
	return
}

type OffchainConfigDigesterRPCServer struct {
	Impl types.OffchainConfigDigester
}

func (o *OffchainConfigDigesterRPCServer) ConfigDigest(args types.ContractConfig, resp *types.ConfigDigest) (err error) {
	*resp, err = o.Impl.ConfigDigest(args)
	return
}

func (o *OffchainConfigDigesterRPCServer) ConfigDigestPrefix(args interface{}, resp *types.ConfigDigestPrefix) error {
	*resp = o.Impl.ConfigDigestPrefix()
	return nil
}

type ReportCodecRPC struct {
	client *rpc.Client
}

func (r *ReportCodecRPC) BuildReport(observations []median.ParsedAttributedObservation) (resp types.Report, err error) {
	err = r.client.Call("Plugin.BuildReport", observations, &resp)
	return
}

func (r *ReportCodecRPC) MedianFromReport(report types.Report) (resp *big.Int, err error) {
	err = r.client.Call("Plugin.MedianFromReport", report, &resp)
	return
}

type ReportCodecRPCServer struct {
	Impl median.ReportCodec
}

func (r *ReportCodecRPCServer) BuildReport(args []median.ParsedAttributedObservation, resp *types.Report) (err error) {
	*resp, err = r.Impl.BuildReport(args)
	return
}

func (r *ReportCodecRPCServer) MedianFromReport(report types.Report, resp **big.Int) (err error) {
	*resp, err = r.Impl.MedianFromReport(report)
	return
}

type MedianContractRPC struct {
	client *rpc.Client
}

func (m *MedianContractRPC) LatestTransmissionDetails(ctx context.Context) (configDigest types.ConfigDigest, epoch uint32, round uint8, latestAnswer *big.Int, latestTimestamp time.Time, err error) {
	var resp LatestTransmissionDetailsResponse
	err = m.client.Call("Plugin.LatestTransmissionDetails", new(interface{}), &resp)
	if err != nil {
		err = errors.Wrap(err, "plugin Call failed")
		return
	}
	return resp.ConfigDigest, resp.Epoch, resp.Round, resp.LatestAnswer, resp.LatestTimestamp, nil
}

func (m *MedianContractRPC) LatestRoundRequested(ctx context.Context, lookback time.Duration) (configDigest types.ConfigDigest, epoch uint32, round uint8, err error) {
	var resp LatestRoundRequestedResponse
	err = m.client.Call("Plugin.LatestRoundRequested", lookback, &resp)
	if err != nil {
		err = errors.Wrap(err, "plugin Call failed")
		return
	}
	return resp.ConfigDigest, resp.Epoch, resp.Round, nil
}

type LatestTransmissionDetailsResponse struct {
	ConfigDigest    types.ConfigDigest
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp time.Time
}

type LatestRoundRequestedResponse struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint32
	Round        uint8
}

type MedianContractRPCServer struct {
	Impl median.MedianContract
}

func (m *MedianContractRPCServer) LatestTransmissionDetails(args interface{}, resp *LatestTransmissionDetailsResponse) (err error) {
	resp.ConfigDigest, resp.Epoch, resp.Round, resp.LatestAnswer, resp.LatestTimestamp, err = m.Impl.LatestTransmissionDetails(context.Background())
	return
}

func (m *MedianContractRPCServer) LatestRoundRequested(lookback time.Duration, resp *LatestRoundRequestedResponse) (err error) {
	resp.ConfigDigest, resp.Epoch, resp.Round, err = m.Impl.LatestRoundRequested(context.Background(), lookback)
	return
}
