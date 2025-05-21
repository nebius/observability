package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

var NebiusObjectStorage = dashboard.NewDashboardBuilder("Nebius Object Storage").
	Uid("nebius-object-storage").
	Description("Dashboard provides an overview of the Nebius Object Storage. https://docs.nebius.com/object-storage").
	Tags([]string{"Nebius", "Object Storage"}).
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
		dashboard.NewQueryVariableBuilder("bucket").
			Datasource(DatasourceRef).
			Query(dashboard.StringOrMap{
				String: New("label_values(buckets_stat_quantity, bucket)"),
			}).
			Multi(true).
			AllowCustomValue(false).
			IncludeAll(true).
			AllValue(".*"),
	).
	WithRow(
		dashboard.NewRowBuilder("Requests for $bucket").
			Repeat("bucket"),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Read requests").
		Description("Number of requests made to retrieve object contents from storage. Measured in requests per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(handler) (rate(request_rate{bucket="$bucket", operation_type="read"}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("{{handler}}"),
		).
		Unit(units.RequestsPerSecond).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(6),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Modify requests").
		Description("Number of requests made to upload objects or modify object contents. Measured in requests per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(handler) (rate(request_rate{bucket="$bucket", operation_type="mutate"}[$__rate_interval])) OR on() vector(0)`).
			LegendFormat("{{handler}}"),
		).
		Unit(units.RequestsPerSecond).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(6),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Traffic").
		Description("Amount of data transferred to and from storage. Measured in bytes per second.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(http_bytes_sent{bucket="$bucket"}[$__rate_interval])) OR vector(0)`).
			LegendFormat("Download"),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(http_bytes_received{bucket="$bucket"}[$__rate_interval])) OR vector(0)`).
			LegendFormat("Upload"),
		).
		Unit(units.BytesPerSecondIEC).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(6),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("API error codes").
		Description("Number of errors when accessing our API. Measured in errors per 5 minutes.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(handler, http_code, api_error_code) (increase(http_errors_total{bucket="$bucket"}[5m]))`).
			LegendFormat("{{handler}}:{{http_code}}:{{api_error_code}}"),
		).
		Unit(units.Short).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(6),
	).
	WithRow(
		dashboard.NewRowBuilder("Objects statistics"),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Object counts").
		Description("Number of objects in storage. Single and multipart objects are counted separately.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(bucket) (max by(bucket, counter) (last_over_time(buckets_stat_quantity{bucket=~"$bucket", counter="simple_objects"}[1m])))`).
			LegendFormat("{{bucket}} simple objects"),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(bucket) (max by(bucket, counter) (last_over_time(buckets_stat_quantity{bucket=~"$bucket", counter="multipart_objects"}[1m])))`).
			LegendFormat("{{bucket}} multipart objects"),
		).
		Unit(units.Short).
		FillOpacity(5).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(8),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Space by object type").
		Description("Amount of storage used for objects of different types: single or multipart.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(bucket) (max by(bucket, counter) (last_over_time(buckets_stat_size{bucket=~"$bucket", counter="simple_objects"}[1m])))`).
			LegendFormat("{{bucket}} simple objects"),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(bucket) (max by(bucket, counter) (last_over_time(buckets_stat_size{bucket=~"$bucket", counter="multipart_objects"}[1m])))`).
			LegendFormat("{{bucket}} multipart objects"),
		).
		Unit(units.BytesIEC).
		FillOpacity(70).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(8),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Total bucket size").
		Description("Storage space used by all objects in a bucket.").
		Datasource(DatasourceRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum by(bucket) (max by(bucket, counter) (last_over_time(buckets_stat_size{bucket=~"$bucket"}[1m])))`).
			LegendFormat("{{bucket}}"),
		).
		Unit(units.BytesIEC).
		FillOpacity(70).
		ShowPoints(common.VisibilityModeNever).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(8).
		Span(8),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()
