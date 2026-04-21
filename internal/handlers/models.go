package handlers

type URLMetric struct {
	Type, Name string
	Value      interface{}
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

//func (m *Metrics) UnmarshalJSON(data []byte) (err error) {
//	type MetricsAlias Metrics
//
//	aliasValue := &struct {
//		*MetricsAlias
//		//Delta *int64   `json:"delta,omitempty"`
//		//Value *float64 `json:"value,omitempty"`
//	}{
//		MetricsAlias: (*MetricsAlias)(m),
//	}
//
//	if err = json.Unmarshal(data, aliasValue); err != nil {
//		return
//	}
//
//	switch m.MType {
//	case models.Gauge:
//		if m.Value == nil {
//			return fmt.Errorf("no value provided")
//		}
//		m.Delta = nil
//	case models.Counter:
//		if m.Delta == nil {
//			return fmt.Errorf("no delta provided")
//		}
//		m.Value = nil
//	default:
//		return fmt.Errorf("unknown metric type: %s, expected '%s' or '%s'", m.MType, models.Gauge, models.Counter)
//	}
//
//	//m.Delta = aliasValue.Delta
//	//m.Value = aliasValue.Value
//	return
//}
