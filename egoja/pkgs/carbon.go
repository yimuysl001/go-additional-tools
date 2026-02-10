package pkgs

import "github.com/dromara/carbon/v2"

func Carbon() {
	RegisterImport("carbon", map[string]any{
		"Default":     carbon.Default{},
		"Carbon":      carbon.Carbon{},
		"SetDefault":  carbon.SetDefault,
		"SetTimezone": carbon.SetTimezone,
		"SetLocation": carbon.SetLocation,
		"NewCarbon":   carbon.NewCarbon,

		"Now":                      carbon.Now,
		"Tomorrow":                 carbon.Tomorrow,
		"Yesterday":                carbon.Yesterday,
		"Parse":                    carbon.Parse,
		"ParseByFormat":            carbon.ParseByFormat,
		"CreateFromTimestamp":      carbon.CreateFromTimestamp,
		"CreateFromTimestampMilli": carbon.CreateFromTimestampMilli,
		"CreateFromTimestampMicro": carbon.CreateFromTimestampMicro,
		"CreateFromTimestampNano":  carbon.CreateFromTimestampNano,
		"CreateFromDateTime":       carbon.CreateFromDateTime,
		"CreateFromDateTimeMilli":  carbon.CreateFromDateTimeMilli,
		"CreateFromDateTimeMicro":  carbon.CreateFromDateTimeMicro,
		"CreateFromDateTimeNano":   carbon.CreateFromDateTimeNano,
		"CreateFromDate":           carbon.CreateFromDate,
		"CreateFromDateMilli":      carbon.CreateFromDateMilli,
		"CreateFromDateMicro":      carbon.CreateFromDateMicro,
		"CreateFromDateNano":       carbon.CreateFromDateNano,
		"CreateFromTime":           carbon.CreateFromTime,
		"CreateFromTimeMilli":      carbon.CreateFromTimeMilli,
		"CreateFromTimeMicro":      carbon.CreateFromTimeMicro,
		"CreateFromTimeNano":       carbon.CreateFromTimeNano,
		"Max":                      carbon.Max,
		"Min":                      carbon.Min,
		"CreateFromStdTime":        carbon.CreateFromStdTime,
		"Local":                    carbon.Local,
		//"CET":                      carbon.CET,
		"Shanghai":        carbon.Shanghai,
		"CreateFromLunar": carbon.CreateFromLunar,
	})
}
