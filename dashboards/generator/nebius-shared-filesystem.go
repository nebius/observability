package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

var NebiusSharedFilesystem = dashboard.NewDashboardBuilder("Nebius Shared Filesystem").
	Uid("nebius-shared-filesystem").
	Description("Dashboard provides an overview of the Nebius Shared Filesystems.").
	Tags([]string{"Nebius", "NBS"}).
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
	WithVariable(
		dashboard.NewQueryVariableBuilder("filestore").
			Datasource(DatasourceRef).
			Query(dashboard.StringOrMap{
				String: New("label_values(filestore)"),
			}).
			AllowCustomValue(false),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS read latency (quantiles)").
		Description("Percentiles of the filesystem write requests latency. Measured in milliseconds.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, sum by(le) (rate(filestore_read_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, sum by(le) (rate(filestore_read_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, sum by(le) (rate(filestore_read_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, sum by(le) (rate(filestore_read_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, sum by(le) (rate(filestore_read_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Milliseconds).
		AxisSoftMin(0).
		LineWidth(1.5).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS write latency (quantiles)").
		Description("Percentiles of the filesystem read requests latency. Measured in milliseconds.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.1, sum by(le) (rate(filestore_write_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p10").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, sum by(le) (rate(filestore_write_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, sum by(le) (rate(filestore_write_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, sum by(le) (rate(filestore_write_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, sum by(le) (rate(filestore_write_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, sum by(le) (rate(filestore_write_latency_bucket{filestore="$filestore"}[$__rate_interval])))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Milliseconds).
		AxisSoftMin(0).
		LineWidth(1.5).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS read operations").
		Description("Average read IOPS. Measured in operations per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_read_ops{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Read").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(filestore_read_ops_burst{filestore="$filestore"})`).
			LegendFormat("Read Burst").
			Range(),
		).
		Unit(units.IOOpsPerSecond).
		AxisSoftMin(0).
		LineWidth(1.5).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS write operations").
		Description("Average write IOPS. Measured in operations per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_write_ops{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Write").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(filestore_write_ops_burst{filestore="$filestore"})`).
			LegendFormat("Write Burst").
			Range(),
		).
		Unit(units.IOOpsPerSecond).
		AxisSoftMin(0).
		LineWidth(1.5).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS read bytes").
		Description("Average read throughput. Measured in bytes per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_read_bytes{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Read").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(filestore_read_bytes_burst{filestore="$filestore"})`).
			LegendFormat("Read Burst").
			Range(),
		).
		Unit(units.BytesPerSecondIEC).
		AxisSoftMin(0).
		LineWidth(1.5).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS write bytes").
		Description("Average write throughput. Measured in bytes per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_write_bytes{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Write").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(filestore_write_bytes_burst{filestore="$filestore"})`).
			LegendFormat("Write Burst").
			Range(),
		).
		Unit(units.BytesPerSecondIEC).
		AxisSoftMin(0).
		LineWidth(1.5).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS read errors").
		Description("Number of times a filesystem fails to read data.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_read_errors{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Read").
			Range(),
		).
		AxisSoftMin(0).
		LineWidth(1.5).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS write errors").
		Description("Number of times a filesystem fails to write data.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_write_errors{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Write").
			Range(),
		).
		AxisSoftMin(0).
		LineWidth(1.5).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS index operations").
		Description("Number of indexing actions (reads, writes, updates and deletions) performed in a time period.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_index_ops{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Ops").
			Range(),
		).
		AxisSoftMin(0).
		LineWidth(1.5).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("FS index errors").
		Description("Number of failed indexing operations in a time period.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(filestore_index_errors{filestore="$filestore"}[$__rate_interval]))`).
			LegendFormat("Errors").
			Range(),
		).
		AxisSoftMin(0).
		LineWidth(1.5).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()
