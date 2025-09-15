package report

import (
	"encoding/xml"
	"strconv"

	"github.com/aljoshare/commala/internal/utils"
	"github.com/aljoshare/commala/internal/validator"
	"github.com/charmbracelet/log"
)

type TestCase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Time    string   `xml:"time,attr"`
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
	ts := TestSuite{
		Name:       "Commala Validation",
		Time:       "0",
		Failures:   0,
		Errors:     0,
		Skipped:    0,
		Assertions: 1,
		Tc:         make([]*TestCase, 0, 1),
	}
	tss.Ts = append(tss.Ts, &ts)

	od := int64(0)
	of := 0
	oa := 0

	for _, v := range vr {
		tc := TestCase{
			Name: v.Validator,
			Time: v.Duration.String(),
		}
		od = od + v.Duration.Nanoseconds()
		if !v.Valid {
			of = of + 1
		}
		oa = oa + 1
		ts.Tc = append(ts.Tc, &tc)
	}

	ts.Time = strconv.FormatInt(od, 10) + "ns"
	ts.Failures = of
	ts.Assertions = oa
	tss.Time = strconv.FormatInt(od, 10) + "ns"

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
