package commands

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/hyperledger/burrow/execution/state"

	"github.com/hyperledger/burrow/rpc/rpcdump"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/db"

	cli "github.com/jawher/mow.cli"
	"google.golang.org/grpc"
)

var cdc = amino.NewCodec()

func Dump(output Output) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		chainURLOpt := cmd.StringOpt("c chain", "127.0.0.1:10997", "chain to be used in IP:PORT format")
		heightOpt := cmd.IntOpt("h height", 0, "Block height to dump to, defaults to latest block height")
		filename := cmd.StringArg("FILE", "", "Save dump here")
		useJSON := cmd.BoolOpt("j json", false, "Output in json")

		timeoutOpt := cmd.IntOpt("t timeout", 10, "GRPC timeout in seconds (default 10)")
		timeout := time.Duration(*timeoutOpt)

		s := state.NewState(db.NewMemDB())

		cmd.Action = func() {
			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			defer cancel()

			var opts []grpc.DialOption
			opts = append(opts, grpc.WithInsecure())
			conn, err := grpc.Dial(*chainURLOpt, opts...)
			if err != nil {
				output.Fatalf("failed to connect: %v", err)
				return
			}
			dc := rpcdump.NewDumpClient(conn)

			dump, err := dc.GetDump(ctx, &rpcdump.GetDumpParam{Height: uint64(*heightOpt)})
			if err != nil {
				output.Fatalf("failed to retrieve dump: %v", err)
				return
			}

			f, err := os.OpenFile(*filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				output.Fatalf("%s: failed to save dump: %v", *filename, err)
				return
			}

			_, _, err = s.Update(func(ws state.Updatable) error {
				for {
					resp, err := dump.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						output.Fatalf("failed to recv dump: %v", err)
						return err
					}

					// update our temporary state
					if resp.Account != nil {
						ws.UpdateAccount(resp.Account)
					}
					if resp.AccountStorage != nil {
						for _, storage := range resp.AccountStorage.Storage {
							ws.SetStorage(resp.AccountStorage.Address, storage.Key, storage.Value)
						}
					}
					if resp.Name != nil {
						ws.UpdateName(resp.Name)
					}

					var bs []byte
					if *useJSON {
						bs, err = json.Marshal(resp)
						if bs != nil {
							bs = append(bs, []byte("\n")...)
						}
					} else {
						bs, err = cdc.MarshalBinaryLengthPrefixed(resp)
					}
					if err != nil {
						output.Fatalf("failed to marshall dump: %v", *filename, err)
					}

					n, err := f.Write(bs)
					if err == nil && n < len(bs) {
						output.Fatalf("%s: failed to save dump: %v", *filename, err)
					}
				}

				return nil
			})

			if err := f.Close(); err != nil {
				output.Fatalf("%s: failed to save dump: %v", *filename, err)
			}
		}
	}
}
