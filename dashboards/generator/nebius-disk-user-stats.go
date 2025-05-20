package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

var NebiusDiskUserStats = dashboard.NewDashboardBuilder("Nebius Disk").
	Uid("nebius-disk-user-stats").
	Description("Dashboard provides monitoring and visualization of disk performance metrics for Nebius Disks.").
	Tags([]string{"Nebius", "Compute", "Disk"}).
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
		dashboard.NewQueryVariableBuilder("disk").
			Datasource(DatasourceRef).
			Query(dashboard.StringOrMap{
				String: New("label_values(disk)"),
			}).
			AllowCustomValue(false),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Disk read latency (quantiles)").
		Description("Shows disk read latency quantiles in milliseconds.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, rate(disk_read_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, rate(disk_read_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, rate(disk_read_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, rate(disk_read_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, rate(disk_read_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Milliseconds).
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
		Title("Disk write latency (quantiles)").
		Description("Shows disk write latency quantiles in milliseconds.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, rate(disk_write_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, rate(disk_write_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, rate(disk_write_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, rate(disk_write_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, rate(disk_write_latency_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Milliseconds).
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
		Title("Disk read throttler latency (quantiles)").
		Description("Shows disk read throttler latency quantiles in microseconds.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, rate(disk_read_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, rate(disk_read_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, rate(disk_read_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, rate(disk_read_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, rate(disk_read_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Microseconds).
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
		Title("Disk write throttler latency (quantiles)").
		Description("Shows disk write throttler latency quantiles in microseconds.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.5, rate(disk_write_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p50").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.75, rate(disk_write_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p75").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.90, rate(disk_write_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p90").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.95, rate(disk_write_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p95").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`histogram_quantile(0.99, rate(disk_write_throttler_delay_bucket{disk="$disk"}[$__rate_interval]))`).
			LegendFormat("p99").
			Range(),
		).
		Unit(units.Microseconds).
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
		Title("Disk read operations").
		Description("Shows disk read operations per second and burst limit.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`rate(disk_read_ops{disk="$disk"}[$__rate_interval])`).
			LegendFormat("Read").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`disk_read_ops_burst{disk="$disk"}`).
			LegendFormat("Read Burst").
			Range(),
		).
		Unit(units.IOOpsPerSecond).
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
		Title("Disk write operations").
		Description("Shows disk write operations per second and burst limit.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`rate(disk_write_ops{disk="$disk"}[$__rate_interval])`).
			LegendFormat("Write").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`disk_write_ops_burst{disk="$disk"}`).
			LegendFormat("Write Burst").
			Range(),
		).
		Unit(units.IOOpsPerSecond).
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
		Title("Disk read bytes").
		Description("Shows disk read bytes per second and burst limit.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`rate(disk_read_bytes{disk="$disk"}[$__rate_interval])`).
			LegendFormat("Read").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`disk_read_bytes_burst{disk="$disk"}`).
			LegendFormat("Read Burst").
			Range(),
		).
		Unit(units.BytesPerSecondIEC).
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
		Title("Disk write bytes").
		Description("Shows disk write bytes per second and burst limit.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`rate(disk_write_bytes{disk="$disk"}[$__rate_interval])`).
			LegendFormat("Write").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`disk_write_bytes_burst{disk="$disk"}`).
			LegendFormat("Write Burst").
			Range(),
		).
		Unit(units.BytesPerSecondIEC).
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
		Title("Disk used quota").
		Description("Shows disk quota utilization percentage.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`rate(disk_io_quota_utilization_percentage{disk="$disk"}[$__rate_interval])`).
			LegendFormat("Quota").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`disk_io_quota_utilization_percentage_burst{disk="$disk"}`).
			LegendFormat("Quota Burst").
			Range(),
		).
		Unit(units.Percent).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()
