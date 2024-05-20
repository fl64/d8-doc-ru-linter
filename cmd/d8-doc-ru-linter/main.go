package main

import (
	"fmt"
	"os"

	"d8-doc-ru-linter/internal/crd"

	arg "github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

var args struct {
	Debug          bool   `arg:"--debug" help:"turn on debug mode"`
	Fail           bool   `arg:"--fail" help:"fail if there are diffs"`
	OriginCRD      string `arg:"required,-s,--source" help:"origin CRD" placeholder:"<src crd>"`
	DestinationCRD string `arg:"-d,--destination" help:"destination CRD" placeholder:"<dst crd>"`
	OutputFileName string `arg:"-n,--new" help:"file to save merged CRD" placeholder:"<merged crd>"`
}

func main() {
	arg.MustParse(&args)
	if args.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var originCRD crd.CRD
	if err := originCRD.Load(args.OriginCRD); err != nil {
		log.Fatalf("Can't load source CRD: %v", err)
	}

	var destinationCRD crd.CRD
	if args.DestinationCRD != "" {
		if err := destinationCRD.Load(args.DestinationCRD); err != nil {
			log.Fatalf("Can't load destination CRD: %v", err)
		}
	}

	resultCRD, ops := originCRD.CompareWith(destinationCRD)

	if args.OutputFileName != "" {
		opsReport, err := ops.MarshalJSONReport()
		if err != nil {
			log.Fatalf("Can't marshal operations: %v", err)
		}

		fmt.Println(string(opsReport))
	}
	resultCRDYaml, err := resultCRD.Marshal()
	if err != nil {
		log.Fatalf("Can't marshal source CRD: %v", err)
	}

	if args.OutputFileName == "" {
		// show yaml
		fmt.Println(string(resultCRDYaml))
	} else {
		if err = os.WriteFile(args.OutputFileName, resultCRDYaml, 0644); err != nil {
			log.Error("Can't save result:", err)
		}
	}

	if args.Fail {
		if ops.Len() > 0 {
			os.Exit(33)
		}

	}
}
