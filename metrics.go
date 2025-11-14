package main

// getMetrics returns all metric definitions
func getMetrics() map[string]MetricDefinition {
	// Label mappings
	hostportStatsLabels := map[string]string{"durable-id": "port"}
	diskLabels := map[string]string{"location": "location", "serial-number": "serial"}
	volumeLabels := map[string]string{"volume-name": "volume"}
	poolStatsLabels := map[string]string{"pool": "pool", "serial-number": "serial"}
	poolLabels := map[string]string{"name": "pool", "serial-number": "serial"}
	tierLabels := map[string]string{"tier": "tier", "pool": "pool", "serial-number": "serial"}
	controllerLabels := map[string]string{"durable-id": "controller"}
	psuLabels := map[string]string{"durable-id": "psu", "serial-number": "serial"}

	return map[string]MetricDefinition{
		"hostport_data_read": {
			Description: "Data Read",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "data-read-numeric",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_data_written": {
			Description: "Data Written",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "data-written-numeric",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_avg_resp_time_read": {
			Description: "Read Response Time",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "avg-read-rsp-time",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_avg_resp_time_write": {
			Description: "Write Response Time",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "avg-write-rsp-time",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_avg_resp_time": {
			Description: "I/O Response Time",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "avg-rsp-time",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_queue_depth": {
			Description: "Queue Depth",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "queue-depth",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_reads": {
			Description: "Reads",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "number-of-reads",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"hostport_writes": {
			Description: "Writes",
			Sources: []MetricSource{{
				Path:              "host-port-statistics",
				ObjectSelector:    "host-port-statistics",
				PropertySelector:  "number-of-writes",
				PropertiesAsLabel: hostportStatsLabels,
			}},
		},
		"disk_temperature": {
			Description: "Temperature",
			Sources: []MetricSource{{
				Path:              "disks",
				ObjectSelector:    "drive",
				PropertySelector:  "temperature-numeric",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_iops": {
			Description: "IOPS",
			Sources: []MetricSource{{
				Path:              "disk-statistics",
				ObjectSelector:    "disk-statistics",
				PropertySelector:  "iops",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_power_on_hours": {
			Description: "Power on hours",
			Sources: []MetricSource{{
				Path:              "disk-statistics",
				ObjectSelector:    "disk-statistics",
				PropertySelector:  "power-on-hours",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_bps": {
			Description: "Bytes per second",
			Sources: []MetricSource{{
				Path:              "disks",
				ObjectSelector:    "disk-statistics",
				PropertySelector:  "bytes-per-second-numeric",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_avg_resp_time": {
			Description: "Average I/O Response Time",
			Sources: []MetricSource{{
				Path:              "disks",
				ObjectSelector:    "drive",
				PropertySelector:  "avg-rsp-time",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_ssd_life_left": {
			Description: "SSD Life Remaining",
			Sources: []MetricSource{{
				Path:              "disks",
				ObjectSelector:    "drive",
				PropertySelector:  "ssd-life-left-numeric",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_health": {
			Description: "Health",
			Sources: []MetricSource{{
				Path:              "disks",
				ObjectSelector:    "drive",
				PropertySelector:  "health-numeric",
				PropertiesAsLabel: diskLabels,
			}},
		},
		"disk_errors": {
			Description: "Errors",
			Sources: []MetricSource{
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "smart-count-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "smart", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "smart-count-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "smart", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "io-timeout-count-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "io-timeout", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "io-timeout-count-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "io-timeout", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "no-response-count-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "no-response", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "no-response-count-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "no-response", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "spinup-retry-count-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "spinup-retry", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "spinup-retry-count-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "spinup-retry", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-media-errors-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "media-errors", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-media-errors-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "media-errors", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-nonmedia-errors-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "nonmedia-errors", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-nonmedia-errors-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "nonmedia-errors", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-block-reassigns-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "block-reassigns", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-block-reassigns-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "block-reassigns", "port": 2},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-bad-blocks-1",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "bad-blocks", "port": 1},
				},
				{
					Path:              "disk-statistics",
					ObjectSelector:    "disk-statistics",
					PropertySelector:  "number-of-bad-blocks-2",
					PropertiesAsLabel: diskLabels,
					Labels:            map[string]interface{}{"type": "bad-blocks", "port": 2},
				},
			},
		},
		"volume_health": {
			Description: "Health",
			Sources: []MetricSource{{
				Path:              "volumes",
				ObjectSelector:    "volume",
				PropertySelector:  "health-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_iops": {
			Description: "IOPS",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "iops",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_bps": {
			Description: "Bytes per second",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "bytes-per-second-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_reads": {
			Description: "Reads",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "number-of-reads",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_writes": {
			Description: "Writes",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "number-of-writes",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_data_read": {
			Description: "Data Read",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "data-read-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_data_written": {
			Description: "Data Written",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "data-written-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_shared_pages": {
			Description: "Shared Pages",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "shared-pages",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_read_hits": {
			Description: "Read-Cache Hits",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "read-cache-hits",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_read_misses": {
			Description: "Read-Cache Misses",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "read-cache-misses",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_write_hits": {
			Description: "Read-Cache Hits",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "write-cache-hits",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_write_misses": {
			Description: "Read-Cache Misses",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "write-cache-misses",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_small_destage": {
			Description: "Small Destages",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "small-destages",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_full_stripe_write_destages": {
			Description: "Full Stripe Write Destages",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "full-stripe-write-destages",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_read_ahead_ops": {
			Description: "Read-Ahead Operations",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "read-ahead-operations",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_write_cache_space": {
			Description: "Write Cache Space",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "write-cache-space",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_write_cache_percent": {
			Description: "Write Cache Percentage",
			Sources: []MetricSource{{
				Path:              "volume-statistics",
				ObjectSelector:    "volume-statistics",
				PropertySelector:  "write-cache-percent",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_size": {
			Description: "Size",
			Sources: []MetricSource{{
				Path:              "volumes",
				ObjectSelector:    "volume",
				PropertySelector:  "size-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_total_size": {
			Description: "Total Size",
			Sources: []MetricSource{{
				Path:              "volumes",
				ObjectSelector:    "volume",
				PropertySelector:  "total-size-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_allocated_size": {
			Description: "Total Size",
			Sources: []MetricSource{{
				Path:              "volumes",
				ObjectSelector:    "volume",
				PropertySelector:  "allocated-size-numeric",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_blocks": {
			Description: "Blocks",
			Sources: []MetricSource{{
				Path:              "volumes",
				ObjectSelector:    "volume",
				PropertySelector:  "blocks",
				PropertiesAsLabel: volumeLabels,
			}},
		},
		"volume_tier_distribution": {
			Description: "Volume tier distribution",
			Sources: []MetricSource{
				{
					Path:              "volume-statistics",
					ObjectSelector:    "volume-statistics",
					PropertySelector:  "percent-tier-ssd",
					PropertiesAsLabel: volumeLabels,
					Labels:            map[string]interface{}{"tier": "Performance"},
				},
				{
					Path:              "volume-statistics",
					ObjectSelector:    "volume-statistics",
					PropertySelector:  "percent-tier-sas",
					PropertiesAsLabel: volumeLabels,
					Labels:            map[string]interface{}{"tier": "Standard"},
				},
				{
					Path:              "volume-statistics",
					ObjectSelector:    "volume-statistics",
					PropertySelector:  "percent-tier-sata",
					PropertiesAsLabel: volumeLabels,
					Labels:            map[string]interface{}{"tier": "Archive"},
				},
				{
					Path:              "volume-statistics",
					ObjectSelector:    "volume-statistics",
					PropertySelector:  "percent-allocated-rfc",
					PropertiesAsLabel: volumeLabels,
					Labels:            map[string]interface{}{"tier": "RFC"},
				},
			},
		},
		"pool_data_read": {
			Description: "Data Read",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "pool-statistics",
				PropertySelector:  "data-read-numeric",
				PropertiesAsLabel: poolStatsLabels,
			}},
		},
		"pool_data_written": {
			Description: "Data Written",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "pool-statistics",
				PropertySelector:  "data-written-numeric",
				PropertiesAsLabel: poolStatsLabels,
			}},
		},
		"pool_avg_resp_time": {
			Description: "I/O Response Time",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "pool-statistics",
				PropertySelector:  "avg-rsp-time",
				PropertiesAsLabel: poolStatsLabels,
			}},
		},
		"pool_avg_resp_time_read": {
			Description: "Read Response Time",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "pool-statistics",
				PropertySelector:  "avg-read-rsp-time",
				PropertiesAsLabel: poolStatsLabels,
			}},
		},
		"pool_total_size": {
			Description: "Total Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "total-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_available_size": {
			Description: "Available Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "total-avail-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_snapshot_size": {
			Description: "Snapshot Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "snap-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_allocated_pages": {
			Description: "Allocated Pages",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "allocated-pages",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_available_pages": {
			Description: "Available Pages",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "available-pages",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_metadata_volume_size": {
			Description: "Metadata Volume Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "metadata-vol-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_total_rfc_size": {
			Description: "Total RFC Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "total-rfc-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_available_rfc_size": {
			Description: "Available RFC Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "available-rfc-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_reserved_size": {
			Description: "Reserved Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "reserved-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"pool_unallocated_reserved_size": {
			Description: "Unallocated Reserved Size",
			Sources: []MetricSource{{
				Path:              "pools",
				ObjectSelector:    "pools",
				PropertySelector:  "reserved-unalloc-size-numeric",
				PropertiesAsLabel: poolLabels,
			}},
		},
		"tier_reads": {
			Description: "Reads",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "number-of-reads",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"tier_writes": {
			Description: "Writes",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "number-of-writes",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"tier_data_read": {
			Description: "Data Read",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "data-read-numeric",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"tier_data_written": {
			Description: "Data Written",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "data-written-numeric",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"tier_avg_resp_time": {
			Description: "I/O Response Time",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "avg-rsp-time",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"tier_avg_resp_time_read": {
			Description: "Read Response Time",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "avg-read-rsp-time",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"tier_avg_resp_time_write": {
			Description: "Write Response Time",
			Sources: []MetricSource{{
				Path:              "pool-statistics",
				ObjectSelector:    "tier-statistics",
				PropertySelector:  "avg-write-rsp-time",
				PropertiesAsLabel: tierLabels,
			}},
		},
		"enclosure_power": {
			Description: "Power consumption in watts",
			Sources: []MetricSource{{
				Path:              "enclosures",
				ObjectSelector:    "enclosures",
				PropertySelector:  "enclosure-power",
				PropertiesAsLabel: map[string]string{"enclosure-id": "id", "enclosure-wwn": "wwn"},
			}},
		},
		"controller_cpu": {
			Description: "CPU Load",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "cpu-load",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"controller_iops": {
			Description: "IOPS",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "iops",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"controller_bps": {
			Description: "Bytes per second",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "bytes-per-second-numeric",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"controller_read_hits": {
			Description: "Read-Cache Hits",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "read-cache-hits",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"controller_read_misses": {
			Description: "Read-Cache Misses",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "read-cache-misses",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"controller_write_hits": {
			Description: "Write-Cache Hits",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "write-cache-hits",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"controller_write_misses": {
			Description: "Write-Cache Misses",
			Sources: []MetricSource{{
				Path:              "controller-statistics",
				ObjectSelector:    "controller-statistics",
				PropertySelector:  "write-cache-misses",
				PropertiesAsLabel: controllerLabels,
			}},
		},
		"psu_health": {
			Description: "Power-supply unit health",
			Sources: []MetricSource{{
				Path:              "enclosure",
				ObjectSelector:    "power-supplies",
				PropertySelector:  "health-numeric",
				PropertiesAsLabel: psuLabels,
			}},
		},
		"psu_status": {
			Description: "Power-supply unit status",
			Sources: []MetricSource{{
				Path:              "enclosure",
				ObjectSelector:    "power-supplies",
				PropertySelector:  "status-numeric",
				PropertiesAsLabel: psuLabels,
			}},
		},
		"system_health": {
			Description: "System health",
			Sources: []MetricSource{{
				Path:              "system",
				ObjectSelector:    "system-information",
				PropertySelector:  "health-numeric",
				PropertiesAsLabel: map[string]string{},
			}},
		},
	}
}
