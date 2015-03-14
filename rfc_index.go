package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type RFCIndexDocumentID string

func (id RFCIndexDocumentID) Number() (int, error) {
	return strconv.Atoi(string(id)[3:])
}

type RFCIndexStatus string

type RFCIndexFileFormat string

type RFCIndexDayOfMonth int

type RFCIndexMonthName string

type RFCIndexDocumentRef struct {
	DocIDs []RFCIndexDocumentID `xml:"doc-id"`
}

type RFCIndexStream string

type RFCIndexAuthor struct {
	Name         string `xml:"name"`
	Title        string `xml:"title,omitempty"`
	Organization string `xml:"organization,omitempty"`
	OrgAbbrev    string `xml:"org-abbrev,omitempty"`
}

type RFCIndexDate struct {
	Month RFCIndexMonthName  `xml:"month"`
	Day   RFCIndexDayOfMonth `xml:"day,omitempty"`
	Year  int                `xml:"year"`
}

func (d RFCIndexDate) ToRFCPublicationDate() (RFCPublicationDate, error) {
	time, err := time.Parse("January", string(d.Month))
	if err != nil {
		return RFCPublicationDate{}, err
	}

	date := RFCPublicationDate{
		Year:  d.Year,
		Month: time.Month(),
		Day:   int(d.Day),
	}

	return date, nil
}

type RFCIndexFormat struct {
	FileFormat RFCIndexFileFormat `xml:"file-format"`
	CharCount  int                `xml:"char-count"`
	PageCount  int                `xml:"page-count,omitempty"`
}

type RFCIndexKeywords struct {
	Kws []string `xml:"kw"`
}

type RFCIndexAbstract struct {
	Ps []string `xml:"p"`
}

type RFCIndexSTDEntry struct {
	DocID  RFCIndexDocumentID   `xml:"doc-id"`
	Title  string               `xml:"title,omitempty"`
	IsAlso *RFCIndexDocumentRef `xml:"is-also,omitempty"`
}

func (e *RFCIndexSTDEntry) Include(rfcEntry *RFCIndexRFCEntry) bool {
	for _, docID := range e.IsAlso.DocIDs {
		if docID == rfcEntry.DocID {
			return true
		}
	}
	return false
}

type RFCIndexSTDEntries []*RFCIndexSTDEntry

func (es RFCIndexSTDEntries) Get(docID RFCIndexDocumentID) *RFCIndexSTDEntry {
	for _, entry := range es {
		if entry.DocID == docID {
			return entry
		}
	}
	return nil
}

type RFCIndexBCPEntry struct {
	DocID  RFCIndexDocumentID   `xml:"doc-id"`
	Title  string               `xml:"title,omitempty"`
	IsAlso *RFCIndexDocumentRef `xml:"is-also,omitempty"`
}

func (e *RFCIndexBCPEntry) Include(rfcEntry *RFCIndexRFCEntry) bool {
	for _, docID := range e.IsAlso.DocIDs {
		if docID == rfcEntry.DocID {
			return true
		}
	}
	return false
}

type RFCIndexBCPEntries []*RFCIndexBCPEntry

func (es RFCIndexBCPEntries) Get(docID RFCIndexDocumentID) *RFCIndexBCPEntry {
	for _, entry := range es {
		if entry.DocID == docID {
			return entry
		}
	}
	return nil
}

type RFCIndexFYIEntry struct {
	DocID  RFCIndexDocumentID   `xml:"doc-id"`
	Title  string               `xml:"title,omitempty"`
	IsAlso *RFCIndexDocumentRef `xml:"is-also,omitempty"`
}

func (e *RFCIndexFYIEntry) Include(rfcEntry *RFCIndexRFCEntry) bool {
	for _, docID := range e.IsAlso.DocIDs {
		if docID == rfcEntry.DocID {
			return true
		}
	}
	return false
}

type RFCIndexFYIEntries []*RFCIndexFYIEntry

func (es RFCIndexFYIEntries) Get(docID RFCIndexDocumentID) *RFCIndexFYIEntry {
	for _, entry := range es {
		if entry.DocID == docID {
			return entry
		}
	}
	return nil
}

type RFCIndexRFCEntry struct {
	DocID             RFCIndexDocumentID   `xml:"doc-id"`
	Title             string               `xml:"title"`
	Authors           []RFCIndexAuthor     `xml:"author"`
	Date              RFCIndexDate         `xml:"date"`
	Formats           []RFCIndexFormat     `xml:"format,omitempty"`
	Keywords          *RFCIndexKeywords    `xml:"keywords,omitempty"`
	Abstract          *RFCIndexAbstract    `xml:"abstract,omitempty"`
	Draft             string               `xml:"draft,omitempty"`
	Notes             string               `xml:"notes,omitempty"`
	Obsoletes         *RFCIndexDocumentRef `xml:"obsoletes,omitempty"`
	ObsoletedBy       *RFCIndexDocumentRef `xml:"obsoleted-by,omitempty"`
	Updates           *RFCIndexDocumentRef `xml:"updates,omitempty"`
	UpdatedBy         *RFCIndexDocumentRef `xml:"updated-by,omitempty"`
	IsAlso            *RFCIndexDocumentRef `xml:"is-also,omitempty"`
	SeeAlso           *RFCIndexDocumentRef `xml:"see-also,omitempty"`
	CurrentStatus     RFCIndexStatus       `xml:"current-status"`
	PublicationStatus RFCIndexStatus       `xml:"publication-status"`
	Stream            RFCIndexStream       `xml:"stream,omitempty"`
	Area              string               `xml:"area,omitempty"`
	WgAcronym         string               `xml:"wg_acronym,omitempty"`
	ErrataURL         string               `xml:"errata-url,omitempty"`
}

func (e *RFCIndexRFCEntry) IsObsolete() bool {
	return e.ObsoletedBy != nil && len(e.ObsoletedBy.DocIDs) > 0
}

func (e *RFCIndexRFCEntry) IsObsoletedBy(other *RFCIndexRFCEntry) bool {
	if e.ObsoletedBy == nil {
		return false
	}

	for _, docID := range e.ObsoletedBy.DocIDs {
		if docID == other.DocID {
			return true
		}
	}
	return false
}

func (e *RFCIndexRFCEntry) IsUpdatedBy(other *RFCIndexRFCEntry) bool {
	if e.UpdatedBy == nil {
		return false
	}

	for _, docID := range e.UpdatedBy.DocIDs {
		if docID == other.DocID {
			return true
		}
	}
	return false
}

func (e *RFCIndexRFCEntry) ToRFC() (*RFC, error) {
	number, err := e.DocID.Number()
	if err != nil {
		return nil, err
	}

	date, err := e.Date.ToRFCPublicationDate()
	if err != nil {
		return nil, err
	}

	rfc := RFC{
		Number:          number,
		DocumentID:      string(e.DocID),
		Title:           e.Title,
		PublicationDate: date,
	}

	return &rfc, nil
}

type RFCIndexRFCEntries []*RFCIndexRFCEntry

func (es RFCIndexRFCEntries) ToRFCs() ([]*RFC, error) {
	rfcs := make([]*RFC, len(es))

	for i, entry := range es {
		rfc, err := entry.ToRFC()
		if err != nil {
			return nil, err
		}
		rfcs[i] = rfc
	}

	return rfcs, nil
}

func (es RFCIndexRFCEntries) Get(docID RFCIndexDocumentID) *RFCIndexRFCEntry {
	for _, entry := range es {
		if entry.DocID == docID {
			return entry
		}
	}
	return nil
}

func (es RFCIndexRFCEntries) Select(predicate RFCIndexRFCEntryPredicate) RFCIndexRFCEntries {
	var entries RFCIndexRFCEntries

	for _, entry := range es {
		if predicate(entry) {
			entries = append(entries, entry)
		}
	}

	return entries
}

type RFCIndexRFCEntryPredicate func(*RFCIndexRFCEntry) bool

type RFCIndexRFCNotIssuedEntry struct {
	DocID RFCIndexDocumentID `xml:"doc-id"`
}

type RFCIndexRFCNotIssuedEntries []*RFCIndexRFCNotIssuedEntry

type RFCIndex struct {
	XMLName             xml.Name                    `xml:"rfc-index"`
	STDEntries          RFCIndexSTDEntries          `xml:"std-entry"`
	BCPEntries          RFCIndexBCPEntries          `xml:"bcp-entry"`
	FYIEntries          RFCIndexFYIEntries          `xml:"fyi-entry"`
	RFCEntries          RFCIndexRFCEntries          `xml:"rfc-entry"`
	RFCNotIssuedEntries RFCIndexRFCNotIssuedEntries `xml:"rfc-not-issued-entry"`
}

func ParseRFCIndex(doc []byte) (*RFCIndex, error) {
	var rfcIndex RFCIndex

	if err := xml.Unmarshal(doc, &rfcIndex); err != nil {
		return nil, err
	}

	return &rfcIndex, nil
}

type RFCIndexRFCRepository struct {
	RFCIndex *RFCIndex
}

func NewRFCIndexRFCRepository() (*RFCIndexRFCRepository, error) {
	cache := RFCIndexCacheStore{}

	doc, err := cache.Get(RFCIndexDataFormatXML)
	if doc == nil {
		fetcher := RFCIndexFetcher{DataFormat: RFCIndexDataFormatXML}

		doc, err = fetcher.Fetch()
		if err != nil {
			return nil, err
		}
	}

	cache.Put(doc, RFCIndexDataFormatXML)

	rfcIndex, err := ParseRFCIndex(doc)
	if err != nil {
		return nil, err
	}

	repository := RFCIndexRFCRepository{
		RFCIndex: rfcIndex,
	}

	return &repository, nil
}

func (r *RFCIndexRFCRepository) FindAll() ([]*RFC, error) {
	return r.RFCIndex.RFCEntries.ToRFCs()
}

func (r *RFCIndexRFCRepository) FindNonObsolete() ([]*RFC, error) {
	predicate := func(entry *RFCIndexRFCEntry) bool {
		return !entry.IsObsolete()
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindObsoletedBy(number int) ([]*RFC, error) {
	docID := toRFCIndexDocumentID(number)

	other := r.RFCIndex.RFCEntries.Get(docID)
	if other == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return entry.IsObsoletedBy(other)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindObsolete(number int) ([]*RFC, error) {
	docID := toRFCIndexDocumentID(number)

	other := r.RFCIndex.RFCEntries.Get(docID)
	if other == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return other.IsObsoletedBy(entry)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindUpdatedBy(number int) ([]*RFC, error) {
	docID := toRFCIndexDocumentID(number)

	other := r.RFCIndex.RFCEntries.Get(docID)
	if other == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return entry.IsUpdatedBy(other)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindUpdate(number int) ([]*RFC, error) {
	docID := toRFCIndexDocumentID(number)

	other := r.RFCIndex.RFCEntries.Get(docID)
	if other == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return other.IsUpdatedBy(entry)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindBySTDNumber(number int) ([]*RFC, error) {
	docID := toSTDIndexDocumentID(number)

	stdEntry := r.RFCIndex.STDEntries.Get(docID)
	if stdEntry == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return stdEntry.Include(entry)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindByBCPNumber(number int) ([]*RFC, error) {
	docID := toBCPIndexDocumentID(number)

	bcpEntry := r.RFCIndex.BCPEntries.Get(docID)
	if bcpEntry == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return bcpEntry.Include(entry)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindByFYINumber(number int) ([]*RFC, error) {
	docID := toFYIIndexDocumentID(number)

	fyiEntry := r.RFCIndex.FYIEntries.Get(docID)
	if fyiEntry == nil {
		return nil, nil
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return fyiEntry.Include(entry)
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindByCategory(category RFCCategory) ([]*RFC, error) {
	rfcIndexStatus, err := toRFCIndexStatus(category)
	if err != nil {
		return nil, err
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return entry.CurrentStatus == rfcIndexStatus
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func (r *RFCIndexRFCRepository) FindByStream(stream RFCStream) ([]*RFC, error) {
	rfcIndexStream, err := toRFCIndexStream(stream)
	if err != nil {
		return nil, err
	}

	predicate := func(entry *RFCIndexRFCEntry) bool {
		return entry.Stream == rfcIndexStream
	}

	return r.RFCIndex.RFCEntries.Select(predicate).ToRFCs()
}

func toRFCIndexDocumentID(number int) RFCIndexDocumentID {
	return RFCIndexDocumentID(fmt.Sprintf("RFC%04d", number))
}

func toSTDIndexDocumentID(number int) RFCIndexDocumentID {
	return RFCIndexDocumentID(fmt.Sprintf("STD%04d", number))
}

func toBCPIndexDocumentID(number int) RFCIndexDocumentID {
	return RFCIndexDocumentID(fmt.Sprintf("BCP%04d", number))
}

func toFYIIndexDocumentID(number int) RFCIndexDocumentID {
	return RFCIndexDocumentID(fmt.Sprintf("FYI%04d", number))
}

func toRFCIndexStatus(category RFCCategory) (RFCIndexStatus, error) {
	switch category {
	case RFCCategoryProposedStandard:
		return "PROPOSED STANDARD", nil
	case RFCCategoryDraftStandard:
		return "DRAFT STANDARD", nil
	case RFCCategoryInternetStandard:
		return "INTERNET STANDARD", nil
	case RFCCategoryExperimental:
		return "EXPERIMENTAL", nil
	case RFCCategoryInformational:
		return "INFORMATIONAL", nil
	case RFCCategoryHistoric:
		return "HISTORIC", nil
	case RFCCategoryBestCurrentPractice:
		return "BEST CURRENT PRACTICE", nil
	case RFCCategoryUnknown:
		return "UNKNOWN", nil
	}
	return "", fmt.Errorf("cannot recognize RFC category: %v", category)
}

func toRFCIndexStream(stream RFCStream) (RFCIndexStream, error) {
	switch stream {
	case RFCStreamIetf:
		return "IETF", nil
	case RFCStreamIab:
		return "IAB", nil
	case RFCStreamIrtf:
		return "IRTF", nil
	case RFCStreamIndependentSubmission:
		return "INDEPENDENT", nil
	case RFCStreamLegacy:
		return "Legacy", nil
	}
	return "", fmt.Errorf("cannot recognize RFC stream: %v", stream)
}

type RFCIndexFetcher struct {
	DataFormat RFCIndexDataFormat
}

func (f *RFCIndexFetcher) Fetch() ([]byte, error) {
	indexURL, err := f.DataFormat.URL()
	if err != nil {
		return nil, err
	}

	response, err := http.Get(indexURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

type RFCIndexDataFormat int

const (
	RFCIndexDataFormatASCII RFCIndexDataFormat = iota
	RFCIndexDataFormatXML
)

func (f RFCIndexDataFormat) URL() (string, error) {
	switch f {
	case RFCIndexDataFormatASCII:
		return "http://www.rfc-editor.org/in-notes/rfc-index.txt", nil
	case RFCIndexDataFormatXML:
		return "http://www.rfc-editor.org/in-notes/rfc-index.xml", nil
	}
	return "", fmt.Errorf("no URL available for file format: %v", f)
}

func (f RFCIndexDataFormat) FileName() (string, error) {
	switch f {
	case RFCIndexDataFormatASCII:
		return "rfc-index.txt", nil
	case RFCIndexDataFormatXML:
		return "rfc-index.xml", nil
	}
	return "", fmt.Errorf("no file name available for file format: %v", f)
}

type RFCIndexCacheStore struct {
	CacheDirectory string
}

func (s *RFCIndexCacheStore) Put(content []byte, format RFCIndexDataFormat) error {
	cacheDir, err := s.cacheDirectory()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	fileName, err := format.FileName()
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(cacheDir, fileName)

	return ioutil.WriteFile(cacheFile, content, 0644)
}

func (s *RFCIndexCacheStore) Get(format RFCIndexDataFormat) ([]byte, error) {
	cacheDir, err := s.cacheDirectory()
	if err != nil {
		return nil, err
	}

	fileName, err := format.FileName()
	if err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(cacheDir, fileName)

	if _, err = os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, nil
	}

	return ioutil.ReadFile(cacheFile)
}

func (s *RFCIndexCacheStore) cacheDirectory() (string, error) {
	if s.CacheDirectory != "" {
		return s.CacheDirectory, nil
	}

	if dir := GetUserCacheDirectory("rfcs"); dir != "" {
		return dir, nil
	}

	return "", fmt.Errorf("cannot determine the cache directory")
}
