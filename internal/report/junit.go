package report

import (
	"encoding/xml"

	"github.com/aljoshare/commala/internal/utils"
	"github.com/aljoshare/commala/internal/validator"
	"github.com/charmbracelet/log"
)

type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

type TestCase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Time    string   `xml:"time,attr"`
	Failure *Failure `xml:"failure,omitempty"`
}

type TestSuite struct {
	XMLName    xml.Name    `xml:"testsuite"`
	Name       string      `xml:"name,attr"`
	Time       string      `xml:"time,attr"`
	Failures   int         `xml:"failures,attr"`
	Errors     int         `xml:"errors,attr"`
	Skipped    int         `xml:"skipped,attr"`
	Assertions int         `xml:"assertions,attr"`
	Tc         []*TestCase `xml:"testcase"`
}

type TestSuites struct {
	XMLName xml.Name     `xml:"testsuites"`
	Time    string       `xml:"time,attr"`
	Ts      []*TestSuite `xml:"testsuite"`
}

func NewJUnitReport(vr []*validator.ValidationResult, path string) error {
	tss := &TestSuites{Ts: make([]*TestSuite, 0, 1), Time: "0"}
	for _, v := range vr {
		ts := TestSuite{
			Name:       v.Validator,
			Time:       "0",
			Failures:   v.Failures,
			Errors:     0,
			Skipped:    0,
			Assertions: v.Assertions,
			Tc:         make([]*TestCase, 0, 1),
		}
		for i, m := range v.Messages {
			tc := &TestCase{
				Name:    i,
				Time:    v.Duration.String(),
				Failure: nil,
			}
			if !m.Valid {
				tc.Failure = &Failure{
					Message: m.Message,
					Type:    v.Validator,
				}
			}
			ts.Tc = append(ts.Tc, tc)
		}
		tss.Ts = append(tss.Ts, &ts)
	}
	output, err := xml.MarshalIndent(tss, "  ", "    ")
	if err != nil {
		log.Errorf("Can't create JUnit report: %v", err)
		return err
	}

	err = utils.WriteFile(path, output)
	if err != nil {
		log.Errorf("Can't create JUnit report: %v", err)
		return err
	}
	return nil
}
