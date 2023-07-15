package model

import (
	"errors"

	"github.com/cucumber/messages-go/v16"
)

type Step struct {
	Keyword string
	Statement string
	DocString string
	DocStringType string
	DataTable *DataTable
}

type DataTable struct {
	Heads []string
	Rows [][]string
}

func (s *Step) Init(step *messages.Step) error {
	if step == nil {
		return errors.New("expected a 'Given'/'When'/'Then'/'And' statement. Got 'nil'")
	}
	s.Keyword = step.Keyword
	s.Statement = step.Text
	if step.DocString != nil && step.DocString.Content != "" {
		s.DocStringType = step.DocString.MediaType
		s.DocString = step.DocString.Content
	}
	if step.DataTable != nil {
		s.DataTable = &DataTable{}
		s.DataTable.Init(step.DataTable)
	}

	return nil
}

func (d *DataTable) Init(t *messages.DataTable) {
	for _, cell := range t.Rows[0].Cells {
		d.Heads = append(d.Heads, cell.Value)
	}

	for i := 1; i < len(t.Rows); i++ {
		var dataRow []string
		for _, cell := range t.Rows[i].Cells {
			dataRow = append(dataRow, cell.Value)
		}
		d.Rows = append(d.Rows, dataRow)
	}
}