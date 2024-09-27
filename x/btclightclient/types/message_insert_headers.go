package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgInsertHeaders{}

func NewMsgInsertHeaders(creator string, header string, other string) *MsgInsertHeaders {
	return &MsgInsertHeaders{
		Creator: creator,
		Header:  header,
		Other: other,
	}
}

func (msg *MsgInsertHeaders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
