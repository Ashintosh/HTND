package libhoosatwallet

import (
	"encoding/hex"

	"github.com/Hoosat-Oy/hoosatd/app/appmessage"
	"github.com/Hoosat-Oy/hoosatd/cmd/hoosatwallet/daemon/pb"
	"github.com/Hoosat-Oy/hoosatd/domain/consensus/model/externalapi"
	"github.com/Hoosat-Oy/hoosatd/domain/consensus/utils/transactionid"
	"github.com/Hoosat-Oy/hoosatd/domain/consensus/utils/utxo"
)

// HoosatwalletdUTXOsTolibhoosatwalletUTXOs converts a  []*pb.UtxosByAddressesEntry to a []*libhoosatwallet.UTXO
func HoosatwalletdUTXOsTolibhoosatwalletUTXOs(hoosatwalletdUtxoEntires []*pb.UtxosByAddressesEntry) ([]*UTXO, error) {
	UTXOs := make([]*UTXO, len(hoosatwalletdUtxoEntires))
	for i, entry := range hoosatwalletdUtxoEntires {
		script, err := hex.DecodeString(entry.UtxoEntry.ScriptPublicKey.ScriptPublicKey)
		if err != nil {
			return nil, err
		}
		transactionID, err := transactionid.FromString(entry.Outpoint.TransactionId)
		if err != nil {
			return nil, err
		}
		UTXOs[i] = &UTXO{
			UTXOEntry: utxo.NewUTXOEntry(
				entry.UtxoEntry.Amount,
				&externalapi.ScriptPublicKey{
					Script:  script,
					Version: uint16(entry.UtxoEntry.ScriptPublicKey.Version),
				},
				entry.UtxoEntry.IsCoinbase,
				entry.UtxoEntry.BlockDaaScore,
			),
			Outpoint: &externalapi.DomainOutpoint{
				TransactionID: *transactionID,
				Index:         entry.Outpoint.Index,
			},
		}
	}
	return UTXOs, nil
}

// AppMessageUTXOToHoosatwalletdUTXO converts an appmessage.UTXOsByAddressesEntry to a  pb.UtxosByAddressesEntry
func AppMessageUTXOToHoosatwalletdUTXO(appUTXOsByAddressesEntry *appmessage.UTXOsByAddressesEntry) *pb.UtxosByAddressesEntry {
	return &pb.UtxosByAddressesEntry{
		Outpoint: &pb.Outpoint{
			TransactionId: appUTXOsByAddressesEntry.Outpoint.TransactionID,
			Index:         appUTXOsByAddressesEntry.Outpoint.Index,
		},
		UtxoEntry: &pb.UtxoEntry{
			Amount: appUTXOsByAddressesEntry.UTXOEntry.Amount,
			ScriptPublicKey: &pb.ScriptPublicKey{
				Version:         uint32(appUTXOsByAddressesEntry.UTXOEntry.ScriptPublicKey.Version),
				ScriptPublicKey: appUTXOsByAddressesEntry.UTXOEntry.ScriptPublicKey.Script,
			},
			BlockDaaScore: appUTXOsByAddressesEntry.UTXOEntry.BlockDAAScore,
			IsCoinbase:    appUTXOsByAddressesEntry.UTXOEntry.IsCoinbase,
		},
	}
}
