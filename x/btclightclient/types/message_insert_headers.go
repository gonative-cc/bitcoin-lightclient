package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgInsertHeaders{}

func NewMsgInsertHeaders(creator string, headers []string) *MsgInsertHeaders {
	return &MsgInsertHeaders{
		Creator: creator,
		Headers:  headers,
	}
}

func (msg *MsgInsertHeaders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
