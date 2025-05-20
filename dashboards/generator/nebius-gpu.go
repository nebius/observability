package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/gauge"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

var NebiusGPU = dashboard.NewDashboardBuilder("Nebius GPU").
	Uid("nebius-gpu").
	Description("Dashboard to visualize data from the NVIDIA Data Center GPU Manager (DCGM).").
	Tags([]string{"Nebius", "Compute", "GPU"}).
	WithVariable(
		dashboard.NewDatasourceVariableBuilder("datasource").
			Type("prometheus").
			Current(dashboard.VariableOption{
				Text: dashboard.StringOrArrayOfString{
					String: New("o11y-sandbox"),
				},
				Value: dashboard.StringOrArrayOfString{
					String: New("o11y-sandbox"),
				},
			}).
			AllowCustomValue(false),
	).
	WithVariable(
		dashboard.NewQueryVariableBuilder("hostname").
			Label("instance").
			Datasource(dashboard.DataSourceRef{
				Type: New("prometheus"),
				Uid:  New("${datasource}"),
			}).
			Query(dashboard.StringOrMap{
				String: New("label_values(node_uname_info, instance_id)"),
			}).
			AllowCustomValue(false),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Host Load").
		Description("Host load averages indicate system processing demand over 1, 5, and 15-minute intervals. Values reflect the number of processes waiting for resources.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_load1{instance_id="$hostname"}`).
			LegendFormat("1min").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_load5{instance_id="$hostname"}`).
			LegendFormat("5min").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_load15{instance_id="$hostname"}`).
			LegendFormat("15min").
			Range(),
		).
		Unit(units.Short).
		FillOpacity(10).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(6),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("GPU Power Usage").
		Description("Tracks real-time power consumption of each GPU in watts.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`DCGM_FI_DEV_POWER_USAGE{instance_id="$hostname"}`).
			LegendFormat("{{uuid}}").
			Range(),
		).
		Unit(units.Watt).
		Stacking(common.NewStackingConfigBuilder().
			Mode(common.StackingModeNormal),
		).
		FillOpacity(10).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(8),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("GPU Total Power").
		Description("Displays the combined power consumption of all GPUs in the system, measured in watts.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(DCGM_FI_DEV_POWER_USAGE{instance_id="$hostname"})`).
			Instant(),
		).
		Unit(units.Watt).
		Min(0).
		Max(2400).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(1800.0),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(2200.0),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("NVLINK Bandwidth").
		Description("Measures data transfer rates between GPUs over NVLINK interconnects in bytes per second.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`sum(rate(DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL{instance_id="$hostname"}[$__rate_interval]))`).
			LegendFormat("Total").
			Range(),
		).
		Unit(units.BytesPerSecondSI).
		FillOpacity(10).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(7),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("System Memory Usage").
		Description("Shows physical memory consumption, calculated as (Total Memory - Free Memory - Buffers - Cached).").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_memory_MemTotal_bytes{instance_id="$hostname"} - node_memory_MemFree_bytes{instance_id="$hostname"} - node_memory_Buffers_bytes{instance_id="$hostname"} - node_memory_Cached_bytes{instance_id="$hostname"}`).
			LegendFormat("Used").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_memory_Buffers_bytes{instance="$hostname"}`).
			LegendFormat("Buffered").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_memory_Cached_bytes{instance="$hostname"}`).
			LegendFormat("Cached").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_memory_MemFree_bytes{instance="$hostname"}`).
			LegendFormat("Free").
			Range(),
		).
		Unit(units.BytesIEC).
		Stacking(common.NewStackingConfigBuilder().
			Mode(common.StackingModeNormal),
		).
		FillOpacity(10).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(6),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("GPU Temperature").
		Description("Monitors the core temperature of each GPU in degrees Celsius.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`DCGM_FI_DEV_GPU_TEMP{instance_id="$hostname"}`).
			LegendFormat("{{uuid}}").
			Range(),
		).
		Unit(units.Celsius).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(8),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("GPU Avg. Temperature").
		Description("Displays the average temperature across all GPUs in the system.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg(DCGM_FI_DEV_GPU_TEMP{instance_id="$hostname"})`).
			Instant(),
		).
		Unit(units.Celsius).
		Min(0).
		Max(90).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(83.0),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(87.0),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("PCIe Throughput").
		Description("Tracks data transfer rates between GPUs and the host system over PCIe connections in MB/s, showing both transmit (Tx) and receive (Rx) traffic.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(DCGM_FI_PROF_PCIE_TX_BYTES{instance_id="$hostname"}[$__interval])`).
			LegendFormat("{{uuid}} Tx").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(DCGM_FI_PROF_PCIE_RX_BYTES{instance_id="$hostname"}[$__interval])`).
			LegendFormat("{{uuid}} Rx").
			Range(),
		).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "custom.transform",
				Value: common.GraphTransformNegativeY,
			},
		}).
		Unit(units.KilobytesPerSecond).
		FillOpacity(10).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(7),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("Disk Usage").
		Description("Shows total disk consumption as a percentage of total capacity.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg((node_filesystem_size_bytes{instance_id="$hostname", device!="rootfs"} - node_filesystem_avail_bytes{instance_id="$hostname", device!="rootfs"}) / node_filesystem_size_bytes{instance_id="$hostname", device!="rootfs"})`).
			Instant(),
		).
		Unit(units.PercentUnit).
		Min(0).
		Max(1).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(0.75),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(0.90),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("Memory Usage").
		Description("Displays the percentage of system RAM actively in use, excluding cache and buffers.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`(node_memory_MemTotal_bytes{instance_id="$hostname"} - node_memory_MemFree_bytes{instance_id="$hostname"} - node_memory_Buffers_bytes{instance_id="$hostname"} - node_memory_Cached_bytes{instance_id="$hostname"}) / node_memory_MemTotal_bytes{instance_id="$hostname"}`).
			Instant(),
		).
		Unit(units.PercentUnit).
		Min(0).
		Max(1).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(0.80),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(0.90),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("GPU Utilization").
		Description("Tracks the percentage of time the GPU is actively processing workloads.").
		Unit(units.Percent).
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`DCGM_FI_DEV_GPU_UTIL{instance_id="$hostname"}`).
			LegendFormat("{{uuid}}").
			Range(),
		).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(8),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("GPU Total Utilization").
		Description("Displays the utilization across all GPUs in the system.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg(DCGM_FI_DEV_GPU_UTIL{instance_id="$hostname"})`).
			Instant(),
		).
		Unit(units.Percent).
		Min(0).
		Max(100).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(80.0),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(90.0),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Network Throughput").
		Description("Monitors host network traffic across all network interfaces, displaying incoming (In) and outgoing (Out) data rates in kilobits per second.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(node_network_receive_bytes_total{instance_id="$hostname"}[$__interval]) or irate(node_network_receive_bytes{instance_id="$hostname"}[$__interval])`).
			LegendFormat("{{device}} In").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(node_network_transmit_bytes_total{instance_id="$hostname"}[$__interval]) or irate(node_network_transmit_bytes{instance_id="$hostname"}[$__interval])`).
			LegendFormat("{{device}} Out").
			Range(),
		).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "custom.transform",
				Value: common.GraphTransformNegativeY,
			},
		}).
		Unit(units.BytesPerSecondSI).
		FillOpacity(10).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(7),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("Disk Throughput").
		Description("Measures host disk I/O activity in MB/s, tracking both read and write operations.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(node_disk_written_bytes_total{instance_id="$hostname"}[$__interval]) or irate(node_disk_sectors_written{instance_id="$hostname"}[$__interval]) * 512`).
			LegendFormat("{{device}} write").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(node_disk_read_bytes_total{instance_id="$hostname"}[$__interval]) or irate(node_disk_sectors_read{instance_id="$hostname"}[$__interval]) * 512`).
			LegendFormat("{{device}} read").
			Range(),
		).
		Unit(units.BytesPerSecondSI).
		FillOpacity(10).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(6),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("GPU Memory Copy Utilization").
		Description("Measures the percentage of time the GPU's copy engines are actively transferring data between host and device memory.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`DCGM_FI_DEV_MEM_COPY_UTIL{instance_id="$hostname"}`).
			LegendFormat("{{uuid}}").
			Range(),
		).
		Unit(units.Percent).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(8),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("GPU Total Memory Copy Utilization").
		Description("Displays the average utilization of GPU memory copy engines across all GPUs, showing the percentage of time spent transferring data between host and device memory.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg(DCGM_FI_DEV_MEM_COPY_UTIL{instance_id="$hostname"})`).
			Instant(),
		).
		Unit(units.Percent).
		Min(0).
		Max(100).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(70.0),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(90.0),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("InfiniBand Throughput").
		Description("Tracks data transfer rates over InfiniBand network interfaces, a high-performance, low-latency interconnect commonly used in HPC and GPU clusters.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(node_infiniband_port_data_transmitted_bytes_total{instance_id="$hostname", device!~"hfi.+"}[$__interval]) or irate(node_infiniband_port_data_transmitted_bytes_total{instance_id="$hostname", device=~"hfi.+"}[$__interval]) * 2`).
			LegendFormat("{{device}} Out").
			Range(),
		).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`irate(node_infiniband_port_data_received_bytes_total{instance_id="$hostname", device!~"hfi.+"}[$__interval]) or irate(node_infiniband_port_data_received_bytes_total{instance_id="$hostname", device=~"hfi.+"}[$__interval]) * 2`).
			LegendFormat("{{device}} In").
			Range(),
		).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "custom.transform",
				Value: common.GraphTransformNegativeY,
			},
		}).
		Unit(units.BytesPerSecondSI).
		FillOpacity(10).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Legend(common.NewVizLegendOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			ShowLegend(true),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(7),
	).
	WithPanel(stat.NewPanelBuilder().
		Title("Compute instance name").
		Description("Displays the full system identifier.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_uname_info{instance_id="$hostname"}`).
			Format(prometheus.PromQueryFormatTable).
			Instant(),
		).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			Fields("nodename"),
		).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		ColorMode(common.BigValueColorModeNone).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(2).
		Span(3),
	).
	WithPanel(stat.NewPanelBuilder().
		Title("Kernel").
		Description("Displays the Linux kernel version running on the host system.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`node_uname_info{instance_id="$hostname"}`).
			Format(prometheus.PromQueryFormatTable).
			Instant(),
		).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Calcs([]string{"lastNotNull"}).
			Fields("release"),
		).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		ColorMode(common.BigValueColorModeNone).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(2).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("GPU Memory Usage").
		Description("Displays the percentage of GPU memory currently allocated, calculated as used memory divided by total available memory.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`DCGM_FI_DEV_FB_USED{instance_id="$hostname"} / (DCGM_FI_DEV_FB_USED{instance_id="$hostname"} + DCGM_FI_DEV_FB_FREE{instance_id="$hostname"})`).
			LegendFormat("{{uuid}}").
			Range(),
		).
		Unit(units.Percent).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(8),
	).
	WithPanel(gauge.NewPanelBuilder().
		Title("GPU Power Draw").
		Description("Displays the average power consumption of GPUs in watts.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg(DCGM_FI_DEV_POWER_USAGE{instance_id="$hostname"})`).
			Instant(),
		).
		Unit(units.Watt).
		Min(0).
		Max(3200).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().
			Steps([]dashboard.Threshold{
				{
					Color: "rgb(41, 156, 70)",
				},
				{
					Value: New(2600.0),
					Color: "rgb(237, 129, 40)",
				},
				{
					Value: New(3000.0),
					Color: "rgb(212, 74, 58)",
				},
			}),
		).
		Height(5).
		Span(3),
	).
	WithPanel(timeseries.NewPanelBuilder().
		Title("GPU Power Draw").
		Description("Tracks total GPU power consumption over time.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`DCGM_FI_DEV_POWER_USAGE{instance_id="$hostname", uuid=~"GPU-.*"}`).
			LegendFormat("{{uuid}}").
			Range(),
		).
		Unit(units.Watt).
		FillOpacity(10).
		LineWidth(2).
		ShowPoints(common.VisibilityModeNever).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone),
		).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(5).
		Span(7),
	).
	WithPanel(stat.NewPanelBuilder().
		Title("GPU SM Clocks").
		Description("Displays the current Streaming Multiprocessor (SM) clock frequency in MHz.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg(DCGM_FI_DEV_SM_CLOCK{instance_id="$hostname"}) * 1000000`).
			Instant(),
		).
		Unit(units.Hertz).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		ColorMode(common.BigValueColorModeNone).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(3).
		Span(3),
	).
	WithPanel(stat.NewPanelBuilder().
		Title("GPU Memory Clocks").
		Description("Displays the current memory clock frequency of the GPU in GHz.").
		Datasource(dashboard.DataSourceRef{
			Type: New("prometheus"),
			Uid:  New("${datasource}"),
		}).
		WithTarget(prometheus.NewDataqueryBuilder().
			Expr(`avg(DCGM_FI_DEV_MEM_CLOCK{instance_id="$hostname"}) * 1000000`).
			Instant(),
		).
		Unit(units.Hertz).
		Mappings([]dashboard.ValueMapping{
			{
				SpecialValueMap: &dashboard.SpecialValueMap{
					Type: dashboard.MappingTypeSpecialValue,
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: dashboard.SpecialValueMatchNull,
						Result: dashboard.ValueMappingResult{
							Text: New("N/A"),
						},
					},
				},
			},
		}).
		ColorMode(common.BigValueColorModeNone).
		Thresholds(dashboard.NewThresholdsConfigBuilder()).
		Height(3).
		Span(3),
	).
	Time("now-24h", "now").
	Refresh("1m").
	Readonly()
