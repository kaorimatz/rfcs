package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func usageListRFCs(f *flag.FlagSet) func() {
	return func() {
		fmt.Println("Usage: rfcs list [options]")
		fmt.Println("")
		fmt.Println("Options:")
		f.PrintDefaults()
	}
}

func listRFCs(Args []string) error {
	var selectOptions SelectOptions
	var displayOptions DisplayOptions
	var category string
	var stream string

	f := flag.NewFlagSet("list", flag.ContinueOnError)
	f.Usage = usageListRFCs(f)

	f.BoolVar(&selectOptions.ExcludeObsolete, "exclude-obsolete", false, "Exclude obsolete RFCs")
	f.IntVar(&selectOptions.ObsoletedBy, "obsoleted-by", 0, "List RFCs obsoleted by the specified RFC")
	f.IntVar(&selectOptions.Obsolete, "obsolete", 0, "List RFCs obsoleting the specified RFC")
	f.IntVar(&selectOptions.UpdatedBy, "updated-by", 0, "List RFCs updated by the specified RFC")
	f.IntVar(&selectOptions.Update, "update", 0, "List RFCs updating the specified RFC")
	f.IntVar(&selectOptions.STDNumber, "std", 0, "List RFCs labeled with the specified STD")
	f.IntVar(&selectOptions.BCPNumber, "bcp", 0, "List RFCs labeled with the specified BCP")
	f.IntVar(&selectOptions.FYINumber, "fyi", 0, "List RFCs labeled with the specified FYI")
	f.StringVar(&category, "category", "", "List RFCs with the specified category")
	f.StringVar(&stream, "stream", "", "List RFCs in the specified document stream")
	f.BoolVar(&displayOptions.SortByPublicationDate, "sort-by-date", false, "Sort by publication date")
	f.StringVar(&displayOptions.OutputTemplate, "format", "{{.DocumentID}} {{.Title}}", "Format output of each RFC using the given Go template")

	if err := f.Parse(Args); err != nil {
		return nil
	}

	if category != "" {
		rfcCategory, err := toRFCCategory(category)
		if err != nil {
			return err
		}

		selectOptions.Category = &rfcCategory
	}

	if stream != "" {
		rfcStream, err := toRFCStream(stream)
		if err != nil {
			return err
		}

		selectOptions.Stream = &rfcStream
	}

	repository, err := NewRFCIndexRFCRepository()
	if err != nil {
		return err
	}

	command := ListCommand{
		RFCRepository:  repository,
		SelectOptions:  selectOptions,
		DisplayOptions: displayOptions,
	}

	return command.Execute()
}

func toRFCCategory(category string) (RFCCategory, error) {
	switch category {
	case "proposed-standard":
		return RFCCategoryProposedStandard, nil
	case "draft-standard":
		return RFCCategoryDraftStandard, nil
	case "internet-standard":
		return RFCCategoryInternetStandard, nil
	case "experimental":
		return RFCCategoryExperimental, nil
	case "informational":
		return RFCCategoryInformational, nil
	case "historic":
		return RFCCategoryHistoric, nil
	case "bcp":
		return RFCCategoryBestCurrentPractice, nil
	case "unknown":
		return RFCCategoryUnknown, nil
	default:
		return RFCCategory(0), fmt.Errorf("unknown category: %s", category)
	}
}

func toRFCStream(stream string) (RFCStream, error) {
	switch stream {
	case "ietf":
		return RFCStreamIetf, nil
	case "iab":
		return RFCStreamIab, nil
	case "irtf":
		return RFCStreamIrtf, nil
	case "independent":
		return RFCStreamIndependentSubmission, nil
	case "legacy":
		return RFCStreamLegacy, nil
	default:
		return RFCStream(0), fmt.Errorf("unknown stream: %s", stream)
	}
}

func usageGetRFC(f *flag.FlagSet) func() {
	return func() {
		fmt.Println("Usage: rfcs get <RFC number>")
	}
}

func getRFC(Args []string) error {
	f := flag.NewFlagSet("get", flag.ContinueOnError)
	f.Usage = usageGetRFC(f)

	if err := f.Parse(Args); err != nil {
		return nil
	}

	if f.NArg() < 1 {
		f.Usage()
		return nil
	}

	rfcNumber, err := strconv.Atoi(Args[0])
	if err != nil {
		fmt.Println(err)
		f.Usage()
		return nil
	}

	repository := NewDefaultRFCContentRepository()

	command := GetCommand{
		RFCContentRepository: repository,
		RFCNumber:            rfcNumber,
	}

	return command.Execute()
}

func usage() {
	fmt.Println("Usage: rfcs command [command options] [argument]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  list    List RFCs")
	fmt.Println("  get     Fetch RFC")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	var err error
	if os.Args[1] == "list" {
		err = listRFCs(os.Args[2:])
	} else if os.Args[1] == "get" {
		err = getRFC(os.Args[2:])
	} else {
		usage()
	}

	if err != nil {
		fmt.Println(err)
	}
}
