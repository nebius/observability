package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

var DatasourceVar = dashboard.NewDatasourceVariableBuilder("datasource").
	Type("prometheus").
	Current(dashboard.VariableOption{
		Text: dashboard.StringOrArrayOfString{
			String: New("Nebius Services"),
		},
		Value: dashboard.StringOrArrayOfString{
			String: New("Nebius Services"),
		},
	}).
	AllowCustomValue(false)

var DatasourceRef = dashboard.DataSourceRef{
	Type: New("prometheus"),
	Uid:  New("${datasource}"),
}

var DatasourceLoggingVar = dashboard.NewDatasourceVariableBuilder("datasource").
	Type("loki").
	Hide(dashboard.VariableHideHideVariable).
	Current(dashboard.VariableOption{
		Text: dashboard.StringOrArrayOfString{
			String: New("Nebius Logging"),
		},
		Value: dashboard.StringOrArrayOfString{
			String: New("Nebius Logging"),
		},
	}).
	AllowCustomValue(false)

var DatasourceLoggingRef = dashboard.DataSourceRef{
	Type: New("loki"),
	Uid:  New("${datasource}"),
}
