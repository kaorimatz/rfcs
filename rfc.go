package main

import (
	"fmt"
	"time"
)

type RFC struct {
	Number          int
	DocumentID      string
	Title           string
	PublicationDate RFCPublicationDate
}

type RFCRepository interface {
	FindAll() ([]*RFC, error)
	FindNonObsolete() ([]*RFC, error)
	FindObsoletedBy(number int) ([]*RFC, error)
	FindObsolete(number int) ([]*RFC, error)
	FindUpdatedBy(number int) ([]*RFC, error)
	FindUpdate(number int) ([]*RFC, error)
	FindBySTDNumber(number int) ([]*RFC, error)
	FindByBCPNumber(number int) ([]*RFC, error)
	FindByFYINumber(number int) ([]*RFC, error)
	FindByCategory(category RFCCategory) ([]*RFC, error)
	FindByStream(stream RFCStream) ([]*RFC, error)
}

type RFCCategory int

const (
	RFCCategoryProposedStandard RFCCategory = iota
	RFCCategoryDraftStandard
	RFCCategoryInternetStandard
	RFCCategoryExperimental
	RFCCategoryInformational
	RFCCategoryHistoric
	RFCCategoryBestCurrentPractice
	RFCCategoryUnknown
)

type RFCStream int

const (
	RFCStreamIetf RFCStream = iota
	RFCStreamIab
	RFCStreamIrtf
	RFCStreamIndependentSubmission
	RFCStreamLegacy
)

type RFCPublicationDate struct {
	Year  int
	Month time.Month
	Day   int
}

func (d RFCPublicationDate) Before(o RFCPublicationDate) bool {
	if d.Year != o.Year {
		return d.Year < o.Year
	} else if d.Month != o.Month {
		return d.Month < o.Month
	} else {
		return d.Day < o.Day
	}
}

func (d RFCPublicationDate) String() string {
	if d.Day == 0 {
		return fmt.Sprintf("%s %d", d.Month, d.Year)
	}
	return fmt.Sprintf("%d %s %d", d.Day, d.Month, d.Year)
}

type RFCContentRepository interface {
	FindByNumber(number int) ([]byte, error)
}
