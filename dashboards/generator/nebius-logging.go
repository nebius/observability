package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
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
	WithPanel(timeseries.NewPanelBuilder().
		Title("Ingest requests").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(logging_ingest_requests_total{status="ok"}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("requests"),
		).
		Unit(units.RequestsPerSecond).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Ingest requests errors").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(status) (rate(logging_ingest_requests_total{status!="ok"}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("{{status}}"),
		).
		Unit(units.RequestsPerSecond).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Ingest requests duration").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, sum by(le)(rate(logging_ingest_duration_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, sum by(le)(rate(logging_ingest_duration_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, sum by(le)(rate(logging_ingest_duration_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, sum by(le)(rate(logging_ingest_duration_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, sum by(le)(rate(logging_ingest_duration_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Seconds).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Logs save lag").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, sum by(le)(rate(logging_storage_save_lag_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, sum by(le)(rate(logging_storage_save_lag_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, sum by(le)(rate(logging_storage_save_lag_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, sum by(le)(rate(logging_storage_save_lag_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, sum by(le)(rate(logging_storage_save_lag_seconds_bucket{}[$__rate_interval])))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Seconds).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Ingest logs").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(logging_ingest_logs_total{}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("logs"),
		).
		Unit(units.Short).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Ingest bytes").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(logging_ingest_logs_bytes_total{}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("logs"),
		).
		Unit(units.BytesPerSecondIEC).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	/*WithPanel(timeseries.NewPanelBuilder().
		Title("Ingest rejected logs").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(logging_ingest_rejected_logs_total{}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("logs"),
		).
		Unit(units.Short).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).*/
	WithRow(
		dashboard.NewRowBuilder("Read"),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Read requests").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(logging_read_requests_total{status="ok"}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("requests"),
		).
		Unit(units.RequestsPerSecond).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Read requests errors").
		Description("").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(status) (rate(logging_read_requests_total{status!="ok"}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("{{status}}"),
		).
		Unit(units.RequestsPerSecond).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(panelHeight).
		Span(panelSpan),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()
