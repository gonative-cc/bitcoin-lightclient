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
	store := prefix.NewStore(storeAdapter, []byte(types.StoreKey))
	
	// check hash chain
	latestBlock, _ := k.LatestBlock(ctx, &types.QueryLatestBlockRequest{})
	prevHeader, _ := types.NewBlockHeader([]byte(latestBlock.HeaderHex))
	
	for _, header := range headers {
		prevBlockHash := prevHeader.BlockHash()
		if header.PrevBlock.IsEqual(&prevBlockHash) {
			k.Logger().Error("not equal " + header.PrevBlock.String() + " " +prevBlockHash.String())
			return errors.New("This is not hash chain")
		}
		prevHeader = header
	}


	for id, header := range headers {
		headerBytesTMP, _ := types.ByteFromBlockHeader(header)
		var headerBytes types.BTCHeaderBytes = headerBytesTMP

		if key, err := types.HeaderKey(uint64(id + int(latestBlock.Height)) + 1); err != nil {
			return err;
		} else {
			store.Set(key, headerBytes)
		}
	}

	latestHeight := latestBlock.Height + uint64(len(headers))
	lightBlock  := types.NewBTCLightBlock(latestHeight, prevHeader)
	lightBlockBytes, _ := lightBlock.Marshal()
	store.Set(types.LatestBlockKey, lightBlockBytes)
	return nil
}


func (k Keeper) InitGenesisBTCBlock(ctx context.Context, height uint64, headerStr string) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte(types.StoreKey))

	var headerBytes types.BTCHeaderBytes
	_ = headerBytes.UnmarshalHex(headerStr)
	header, _ := headerBytes.NewBlockHeaderFromBytes()
	
	lightBlock := types.NewBTCLightBlock(height, header)
	value, err := lightBlock.Marshal()
	if err != nil {
		return err
	}
	store.Set(types.LatestBlockKey, value)
	return nil
}
