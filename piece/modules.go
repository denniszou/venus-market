package piece

import (
	"context"
	"github.com/filecoin-project/go-fil-markets/piecestore"
	"github.com/filecoin-project/venus-market/builder"
	"github.com/filecoin-project/venus-market/config"
	"github.com/filecoin-project/venus-market/models"
	"github.com/filecoin-project/venus-market/utils"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

// NewProviderPieceStore creates a statestore for storing metadata about pieces
// shared by the piecestorage and retrieval providers
func NewProviderPieceStore(lc fx.Lifecycle, ds models.MetadataDS) (piecestore.PieceStore, error) {
	ps, err := NewDsPieceStore(namespace.Wrap(ds, datastore.NewKey("/storagemarket")))
	if err != nil {
		return nil, err
	}

	ps.OnReady(utils.ReadyLogger("piecestore"))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ps.Start(ctx)
		},
	})
	return ps, nil
}

func NewPieceStorage(pieceStrorageCfg *config.PieceStorage) (IPieceStorage, error) {
	switch pieceStrorageCfg.Type {
	case "local":
		return NewPieceFileStorage(pieceStrorageCfg.Path)
	default:
		return nil, xerrors.Errorf("unsupport piece piecestorage type %s", pieceStrorageCfg.Type)
	}
}

var PieceOpts = func(cfg *config.MarketConfig) builder.Option {
	return builder.Options(
		//piece
		builder.Override(new(IPieceStorage), NewPieceStorage), //save read peiece data
		builder.Override(new(ExtendPieceStore), NewDsPieceStore),
		builder.Override(new(piecestore.PieceStore), builder.From(new(ExtendPieceStore))), //save piece metadata(location)   save to metadata /storagemarket
	)
}
