package monitoring

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/ethereum-optimism/optimism/op-service/pprof"
	krpc "github.com/ethereum-optimism/optimism/op-service/rpc"
)

type metricer interface {
	Serve(context.Context, string, int) error
	StartBalanceMetrics(context.Context, log.Logger, *ethclient.Client, common.Address)
}

// NOTE(pangssu): MaybeStartPprof requires cancelable context to stop http server
func MaybeStartPprof(ctx context.Context, cfg pprof.CLIConfig, l log.Logger) {
	if cfg.Enabled {
		l.Info("starting pprof", "addr", cfg.ListenAddr, "port", cfg.ListenPort)
		go func() {
			if err := pprof.ListenAndServe(ctx, cfg.ListenAddr, cfg.ListenPort); err != nil {
				l.Error("failed to start pprof", "err", err)
			}
		}()
	}
}

// NOTE(pangssu): MaybeStartMetrics requires cancelable context to stop http server
func MaybeStartMetrics(ctx context.Context, cfg metrics.CLIConfig, l log.Logger, m metricer, l1 *ethclient.Client, wallet common.Address) {
	if cfg.Enabled {
		l.Info("starting metrics server", "addr", cfg.ListenAddr, "port", cfg.ListenPort)
		go func() {
			if err := m.Serve(ctx, cfg.ListenAddr, cfg.ListenPort); err != nil {
				l.Error("failed to start metrics server", "err", err)
			}
		}()
		m.StartBalanceMetrics(ctx, l, l1, wallet)
	}
}

func StartRPC(cfg krpc.CLIConfig, version string, opts ...krpc.ServerOption) (*krpc.Server, error) {
	server := krpc.NewServer(cfg.ListenAddr, cfg.ListenPort, version, opts...)
	if err := server.Start(); err != nil {
		return nil, fmt.Errorf("failed to start RPC server: %w", err)
	}

	return server, nil
}
