package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

const (
	panelSpan   = 12
	panelHeight = 8
)

var NebiusLogging = dashboard.NewDashboardBuilder("Nebius Logging").
	Uid("nebius-logging").
	Description("Dashboard provides an overview of the Nebius Logging. https://docs.nebius.com/observability/logging").
	Tags([]string{"Nebius", "Logging"}).
	Link(dashboard.NewDashboardLinkBuilder("Docs").
		Type(dashboard.DashboardLinkTypeLink).
		Url("https://docs.nebius.com/observability").
		TargetBlank(true).
		Icon("doc"),
	).
	Link(dashboard.NewDashboardLinkBuilder("GitHub").
		Type(dashboard.DashboardLinkTypeLink).
		Url("https://github.com/nebius/observability").
		TargetBlank(true).
		Icon("external link"),
	).
	WithVariable(
		DatasourceVar,
	).
	WithRow(
		dashboard.NewRowBuilder("Ingest"),
	).
	WithPanel(
		createBaseStat(
			"Total ingested logs",
			"Total number of logs ingested over time",
			units.Short,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum(increase(logging_ingest_logs_total{}[$__interval]))`),
			),
	).
	WithPanel(
		createBaseStat(
			"Total ingested bytes",
			"Total size of ingested logs over time",
			units.BytesSI,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum(increase(logging_ingest_logs_bytes_total{}[$__interval]))`),
			),
	).
	WithPanel(
		createBasePanel(
			"Ingested logs per second",
			"Number of logs ingested per second",
			units.Short,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum(rate(logging_ingest_logs_total{}[$__rate_interval])) OR on() vector(0)`).
				LegendFormat("logs"),
			),
	).
	WithPanel(
		createBasePanel(
			"Ingest bytes per second",
			"Volume of log data ingested per second in bytes",
			units.BytesPerSecondSI,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum(rate(logging_ingest_logs_bytes_total{}[$__rate_interval])) OR on() vector(0)`).
				LegendFormat("bytes"),
			),
	).
	WithPanel(
		createBasePanel(
			"Ingest requests",
			"Successful ingestion requests per second",
			units.RequestsPerSecond,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum(rate(logging_ingest_requests_total{status="ok"}[$__rate_interval])) OR on() vector(0)`).
				LegendFormat("success"),
			),
	).
	WithPanel(
		createBasePanel(
			"Failed ingest requests",
			"Failed ingestion requests per second, by status",
			units.RequestsPerSecond,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum by(status) (rate(logging_ingest_requests_total{status!="ok"}[$__rate_interval])) OR on() vector(0)`).
				LegendFormat("{{status}}"),
			),
	).
	WithPanel(
		addQuantileTargets(
			createBasePanel(
				"Ingest requests latency",
				"Ingestion request latency (quantiles)",
				units.Seconds,
			),
			"logging_ingest_duration_seconds_bucket",
		),
	).
	WithPanel(
		addQuantileTargets(
			createBasePanel(
				"Log save latency",
				"Time delay between receiving a log and saving it to storage, shown as quantiles",
				units.Seconds,
			),
			"logging_storage_save_lag_seconds_bucket",
		),
	).
	WithRow(
		dashboard.NewRowBuilder("Read"),
	).
	WithPanel(
		createBasePanel(
			"Successful read requests",
			"Successful log read requests per second",
			units.RequestsPerSecond,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum(rate(logging_read_requests_total{status="ok"}[$__rate_interval])) OR on() vector(0)`).
				LegendFormat("success"),
			),
	).
	WithPanel(
		createBasePanel(
			"Failed read requests",
			"Failed read requests per second, by status",
			units.RequestsPerSecond,
		).
			WithTarget(prometheus.NewDataqueryBuilder().
				Expr(`sum by(status) (rate(logging_read_requests_total{status!="ok"}[$__rate_interval])) OR on() vector(0)`).
				LegendFormat("{{status}}"),
			),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()

func createBasePanel(title, description string, unit string) *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title(title).
		Description(description).
		Datasource(DatasourceRef).
		Unit(unit).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		AxisSoftMax(1).
		AxisSoftMin(0).
		Height(panelHeight).
		Span(panelSpan)
}

func createBaseStat(title, description string, unit string) *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title(title).
		Description(description).
		Datasource(DatasourceRef).
		Unit(unit).
		Height(3).
		Span(panelSpan).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Calcs([]string{"sum"}),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		GraphMode(common.BigValueGraphModeNone)
}

func addQuantileTargets(panel *timeseries.PanelBuilder, metricName string) *timeseries.PanelBuilder {
	quantiles := []struct {
		value  string
		legend string
	}{
		{"0.5", "p50"},
		{"0.75", "p75"},
		{"0.90", "p90"},
		{"0.95", "p95"},
		{"0.99", "p99"},
	}

	for _, q := range quantiles {
		panel = panel.WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(` + q.value + `, sum by(le)(rate(` + metricName + `{}[$__rate_interval])))`).
			LegendFormat(q.legend).
			Range(),
		)
	}

	return panel
}
