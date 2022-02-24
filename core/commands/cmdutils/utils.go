package cmdutils

import (
	"fmt"
	cmds "github.com/ipfs/go-ipfs-cmds"

	"github.com/ipfs/go-cid"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

const (
	AllowBigBlockOptionName = "allow-big-block"
	SoftBlockLimit          = 1024 * 1024 // https://github.com/ipfs/go-ipfs/issues/7421#issuecomment-910833499
)

var AllowBigBlockOption cmds.Option

func init() {
	AllowBigBlockOption = cmds.BoolOption(AllowBigBlockOptionName, "Disable block size check and allow creation of blocks bigger than 1MB. WARNING: such blocks won't be transferable over the standard bitswap.").WithDefault(false)
}

// FIXME(BLOCKING): We need to split between the ones that actually need to retrieve
//  it from the DAG service from the ones that already formed the DAG before insering it
// (ipfs dag put).
func CheckCIDSize(req *cmds.Request, c cid.Cid, dagAPI coreiface.APIDagService) error {
	n, err := dagAPI.Get(req.Context, c)
	if err != nil {
		return err
	}

	nodeSize, err := n.Size()
	if err != nil {
		return err
	}

	return CheckBlockSize(req, nodeSize)
}

func CheckBlockSize(req *cmds.Request, size uint64) error {
	allowAnyBlockSize, _ := req.Options[AllowBigBlockOptionName].(bool)
	if allowAnyBlockSize {
		return nil
	}

	// We do not allow producing blocks bigger than 1 MiB to avoid errors
	// when transmitting them over BitSwap. The 1 MiB constant is an
	// unenforced and undeclared rule of thumb hard-coded here.
	if size > SoftBlockLimit {
		return fmt.Errorf("produced block is over 1MB, object API is deprecated and does not support HAMT-sharding: to create big directories, please use the files API (MFS)")
	}
	return nil

}
