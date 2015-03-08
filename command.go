package main

import (
	"fmt"
	"os"
	"sort"
	"text/template"
)

type SelectOptions struct {
	ExcludeObsolete bool
	ObsoletedBy     int
	Obsolete        int
	UpdatedBy       int
	Update          int
	STDNumber       int
	BCPNumber       int
	FYINumber       int
	Category        *RFCCategory
	Stream          *RFCStream
}

type DisplayOptions struct {
	SortByPublicationDate bool
	OutputTemplate        string
}

type ListCommand struct {
	RFCRepository  RFCRepository
	DisplayOptions DisplayOptions
	SelectOptions  SelectOptions
}

func (c *ListCommand) Execute() error {
	var rfcs []*RFC
	var err error

	if c.SelectOptions.ExcludeObsolete {
		rfcs, err = c.RFCRepository.FindNonObsolete()
	} else if c.SelectOptions.ObsoletedBy != 0 {
		rfcs, err = c.RFCRepository.FindObsoletedBy(c.SelectOptions.ObsoletedBy)
	} else if c.SelectOptions.Obsolete != 0 {
		rfcs, err = c.RFCRepository.FindObsolete(c.SelectOptions.Obsolete)
	} else if c.SelectOptions.UpdatedBy != 0 {
		rfcs, err = c.RFCRepository.FindUpdatedBy(c.SelectOptions.UpdatedBy)
	} else if c.SelectOptions.Update != 0 {
		rfcs, err = c.RFCRepository.FindUpdate(c.SelectOptions.Update)
	} else if c.SelectOptions.STDNumber != 0 {
		rfcs, err = c.RFCRepository.FindBySTDNumber(c.SelectOptions.STDNumber)
	} else if c.SelectOptions.BCPNumber != 0 {
		rfcs, err = c.RFCRepository.FindByBCPNumber(c.SelectOptions.BCPNumber)
	} else if c.SelectOptions.FYINumber != 0 {
		rfcs, err = c.RFCRepository.FindByFYINumber(c.SelectOptions.FYINumber)
	} else if c.SelectOptions.Category != nil {
		rfcs, err = c.RFCRepository.FindByCategory(*c.SelectOptions.Category)
	} else if c.SelectOptions.Stream != nil {
		rfcs, err = c.RFCRepository.FindByStream(*c.SelectOptions.Stream)
	} else {
		rfcs, err = c.RFCRepository.FindAll()
	}
	if err != nil {
		return err
	}

	if c.DisplayOptions.SortByPublicationDate {
		sort.Stable(ByPublicationDate(rfcs))
	}

	tmpl, err := template.New("").Parse(c.DisplayOptions.OutputTemplate)
	if err != nil {
		return err
	}

	for _, rfc := range rfcs {
		tmpl.Execute(os.Stdout, rfc)
		fmt.Println("")
	}

	return nil
}

type ByPublicationDate []*RFC

func (r ByPublicationDate) Len() int {
	return len(r)
}

func (r ByPublicationDate) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ByPublicationDate) Less(i, j int) bool {
	return r[i].PublicationDate.Before(r[j].PublicationDate)
}

type GetCommand struct {
	RFCContentRepository RFCContentRepository
	RFCNumber            int
}

func (c *GetCommand) Execute() error {
	content, err := c.RFCContentRepository.FindByNumber(c.RFCNumber)
	if err != nil {
		return err
	}

	fmt.Println(string(content))

	return nil
}
