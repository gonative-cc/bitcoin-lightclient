package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/btcsuite/btcd/wire"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitcoin-lightclient/x/btclightclient/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InsertHeader(ctx context.Context, headers []*wire.BlockHeader) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	
	// check hash chain
	prevHeader := headers[0]
	for _, header := range headers[1:] {
		prevBlockHash := prevHeader.BlockHash()
		if header.PrevBlock != prevBlockHash {
			return errors.New("This is not hash chain")
		}
		prevHeader = header
	}

	for id, header := range headers {
		headerBytes, _ := types.ByteFromBlockHeader(header)
		if key, err := types.HeaderKey(uint64(id)); err != nil {
			return err;
		} else {
			store.Set(key, headerBytes)
		}
	}

	lightBlock  := types.NewBTCLightBlock(10, prevHeader)
	lightBlockBytes, _ := lightBlock.Marshal()
	store.Set(types.LatestBlockKey, lightBlockBytes)
	return nil
}
