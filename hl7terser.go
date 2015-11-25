package hl7terser

import (
	"fmt"

	"github.com/kdar/health/hl7"
)

type Query struct {
	Segment               string
	HasSegmentRepeat      bool
	SegmentRepeat         int
	HasField              bool
	Field                 int
	HasFieldRepeat        bool
	FieldRepeat           int
	HasComponent          bool
	Component             int
	HasComponentRepeat    bool
	ComponentRepeat       int
	HasSubComponent       bool
	SubComponent          int
	HasSubComponentRepeat bool
	SubComponentRepeat    int
}

func New(segment string, segmentRepeat, field, fieldRepeat, component, componentRepeat, subComponent, subComponentRepeat int) Query {
	if component == 0 {
		component = 1
	}

	if subComponent == 0 {
		subComponent = 1
	}

	return Query{
		Segment:            segment,
		SegmentRepeat:      segmentRepeat,
		Field:              field,
		FieldRepeat:        fieldRepeat,
		Component:          component,
		ComponentRepeat:    componentRepeat,
		SubComponent:       subComponent,
		SubComponentRepeat: subComponentRepeat,
	}
}

func (q Query) String() string {
	s := q.Segment

	if q.HasSegmentRepeat {
		s += fmt.Sprintf("(%d)", q.SegmentRepeat)
	}

	if !q.HasField {
		return s
	}

	s += fmt.Sprintf("-%d", q.Field)

	if q.HasFieldRepeat {
		s += fmt.Sprintf("(%d)", q.FieldRepeat)
	}

	if !q.HasComponent {
		return s
	}

	s += fmt.Sprintf("-%d", q.Component)

	if q.HasComponentRepeat {
		s += fmt.Sprintf("(%d)", q.ComponentRepeat)
	}

	if !q.HasSubComponent {
		return s
	}

	s += fmt.Sprintf("-%d", q.SubComponent)

	if q.HasComponentRepeat {
		s += fmt.Sprintf("(%d)", q.SubComponentRepeat)
	}

	return s
}

func (q Query) Get(m Message) (hl7.Data, bool) {
	s, ok := m.Segment(q.Segment, q.SegmentRepeat)
	if !ok {
		return nil, false
	}

	if q.Field == 0 {
		return s, true
	}

	f, ok := s.Index(q.Field)
	if !ok {
		return nil, false
	}

	if fr, ok := f.(hl7.Repeated); ok {
		f, ok = fr.Index(q.FieldRepeat)
		if !ok {
			return nil, false
		}
	} else if q.FieldRepeat != 0 {
		return nil, false
	}

	c, ok := f.Index(q.Component - 1)
	if !ok {
		return nil, false
	}

	if cr, ok := c.(hl7.Repeated); ok {
		c, ok = cr.Index(q.ComponentRepeat)
		if !ok {
			return nil, false
		}
	} else if q.ComponentRepeat != 0 {
		return nil, false
	}

	sc, ok := c.Index(q.SubComponent - 1)
	if !ok {
		return nil, false
	}

	if scr, ok := sc.(hl7.Repeated); ok {
		sc, ok = scr.Index(q.SubComponentRepeat)
		if !ok {
			return nil, false
		}
	} else if q.SubComponentRepeat != 0 {
		return nil, false
	}

	return sc, true
}

func (q Query) GetString(m Message) string {
	r, ok := q.Get(m)
	if !ok {
		return ""
	}

	f, ok := r.(hl7.Field)
	if !ok {
		return ""
	}

	return f.String()
}
