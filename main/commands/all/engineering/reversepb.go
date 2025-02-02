package engineering

import (
	"bytes"
	"flag"
	"github.com/golang/protobuf/proto"
	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/v2ray-core/v4/infra/conf/jsonpb"
	"github.com/v2fly/v2ray-core/v4/infra/conf/v2jsonpb"
	"github.com/v2fly/v2ray-core/v4/main/commands/base"
	"io"
	"os"
)

var cmdReversePb = &base.Command{
	UsageLine: "{{.Exec}} engineering reversepb [-f format]",
	Flag: func() flag.FlagSet {
		fs := flag.NewFlagSet("", flag.ExitOnError)
		configFormat = fs.String("f", "v2jsonpb", "")
		return *fs
	}(),
	Run: func(cmd *base.Command, args []string) {
		cmd.Flag.Parse(args)
		configIn := bytes.NewBuffer(nil)
		io.Copy(configIn, os.Stdin)
		var conf core.Config
		if err := proto.Unmarshal(configIn.Bytes(), &conf); err != nil {
			base.Fatalf("%s", err)
		}
		switch *configFormat {
		case "jsonpb":
			if err := jsonpb.DumpJsonPb(&conf, os.Stdout); err != nil {
				base.Fatalf("%s", err)
			}
		case "v2jsonpb":
			if value, err := v2jsonpb.DumpV2JsonPb(&conf); err != nil {
				base.Fatalf("%s", err)
			} else {
				io.Copy(os.Stdout, bytes.NewReader(value))
			}
		}
	}}
