package main

import (
	"fmt"
	
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

const (
	bucketFilter    = `__bucket__="$bucket"`
	clusterFilter   = `k8s_cluster_id=~"$cluster"`
	namespaceFilter = `k8s_namespace_name=~"$namespace"`
	searchFilter    = `|="$search"`
)

var (
	fullFilter = fmt.Sprintf("{%s, %s, %s} %s", bucketFilter, clusterFilter, namespaceFilter, searchFilter)
)

var NebiusMk8sLogs = dashboard.NewDashboardBuilder("Nebius mk8s logs").
	Uid("nebius-mk8s-logs").
	Description("Overview of logs from Nebius Managed Kubernetes (mk8s) clusters. https://docs.nebius.com/observability/logging").
	Tags([]string{"Nebius", "Logging", "mk8s"}).
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
	Link(dashboard.NewDashboardLinkBuilder("View in Explore").
		Type(dashboard.DashboardLinkTypeLink).
		Url(fmt.Sprintf(`/explore?orgId=1&left=["now-1h","now","Loki",{"expr":%q},{"ui":[true,true,true,"none"]}]`, fullFilter)).
		TargetBlank(true).
		KeepTime(true).
		Icon("link"),
	).
	WithVariable(
		DatasourceLoggingVar,
	).
	WithVariable(
		dashboard.NewConstantVariableBuilder("bucket").
			Value(dashboard.StringOrMap{
				String: New("default"),
			}).
			AllowCustomValue(false),
	).
	WithVariable(
		dashboard.NewQueryVariableBuilder("cluster").
			Datasource(DatasourceLoggingRef).
			Query(dashboard.StringOrMap{
				String: New(fmt.Sprintf("label_values({%s}, k8s_cluster_id)", bucketFilter)),
			}).
			Multi(true).
			AllowCustomValue(false).
			IncludeAll(true).
			AllValue(".*"),
	).
	WithVariable(
		dashboard.NewQueryVariableBuilder("namespace").
			Datasource(DatasourceLoggingRef).
			Query(dashboard.StringOrMap{
				String: New(fmt.Sprintf("label_values({%s, %s}, k8s_namespace_name)", bucketFilter, clusterFilter)),
			}).
			Multi(true).
			AllowCustomValue(false).
			IncludeAll(true).
			AllValue(".*"),
	).
	WithVariable(
		dashboard.NewQueryVariableBuilder("pod").
			Datasource(DatasourceLoggingRef).
			Query(dashboard.StringOrMap{
				String: New(fmt.Sprintf("label_values({%s, %s, %s}, k8s_pod_name)", bucketFilter, clusterFilter, namespaceFilter)),
			}).
			Multi(true).
			AllowCustomValue(false).
			IncludeAll(true).
			AllValue(".*"),
	).
	WithVariable(
		dashboard.NewTextBoxVariableBuilder("search").
			Label("Log Search"),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Datasource(DatasourceLoggingRef).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(fmt.Sprintf(`sum(count_over_time(%s [$__interval]))`, fullFilter)).
			LegendFormat("Logs"),
		).DrawStyle(common.GraphDrawStyleBars).
		MaxDataPoints(300).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(false),
		).AxisPlacement(common.AxisPlacementAuto).
		Min(0).
		Height(3).
		Span(24),
	).
	WithPanel(logs.NewPanelBuilder().
		Datasource(DatasourceLoggingRef).
		WithTarget(loki.NewDataqueryBuilder().
			MaxLines(100).
			// direction https://github.com/grafana/grafana-foundation-sdk/pull/840
			Expr(fullFilter)).
		ShowTime(true).
		EnableLogDetails(true).
		Height(20).
		Span(24),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()
