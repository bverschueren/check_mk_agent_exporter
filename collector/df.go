package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

var (
	filesystemLabelNames = []string{"device", "mountpoint", "fstype"}
)

type dfCollector struct {
	SizeDesc       *prometheus.Desc
	UsedDesc       *prometheus.Desc
	AvailDesc      *prometheus.Desc
	PercentageDesc *prometheus.Desc
}

type filesystemLabels struct {
	device, mountPoint, fsType string
}

type filesystemStats struct {
	size, used, avail, percentage float64
	labels                        filesystemLabels
}

func init() {
	registerCollector("df", NewDfCollector)
}

func NewDfCollector() (Collector, error) {
	subsystem := "df"

	SizeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fs_total_size"),
		"Filesystem total size",
		filesystemLabelNames, nil,
	)
	UsedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fs_used_size"),
		"Filesystem used size",
		filesystemLabelNames, nil,
	)
	AvailDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fs_available_size"),
		"Filesystem available size",
		filesystemLabelNames, nil,
	)
	PercentageDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fs_percentage_used"),
		"Filesystem used percentage",
		filesystemLabelNames, nil,
	)
	return dfCollector{
		SizeDesc:       SizeDesc,
		UsedDesc:       UsedDesc,
		AvailDesc:      AvailDesc,
		PercentageDesc: PercentageDesc,
	}, nil

}

func (d dfCollector) Update(unparsedStats *[]string, ch chan<- prometheus.Metric) error {
	stats := d.parseStats(unparsedStats)

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			d.SizeDesc, prometheus.GaugeValue,
			s.size, s.labels.device, s.labels.mountPoint, s.labels.fsType,
		)
		ch <- prometheus.MustNewConstMetric(
			d.UsedDesc, prometheus.GaugeValue,
			s.used, s.labels.device, s.labels.mountPoint, s.labels.fsType,
		)
		ch <- prometheus.MustNewConstMetric(
			d.AvailDesc, prometheus.GaugeValue,
			s.avail, s.labels.device, s.labels.mountPoint, s.labels.fsType,
		)
		ch <- prometheus.MustNewConstMetric(
			d.PercentageDesc, prometheus.GaugeValue,
			s.percentage, s.labels.device, s.labels.mountPoint, s.labels.fsType,
		)
	}
	return nil
}

func (c dfCollector) parseStats(unparsedStats *[]string) []filesystemStats {

	stats := []filesystemStats{}

	for _, stat := range *unparsedStats {
		fields := strings.Fields(stat)
		f_size, _ := strconv.ParseFloat(fields[2], 64)
		f_used, _ := strconv.ParseFloat(fields[3], 64)
		f_avail, _ := strconv.ParseFloat(fields[4], 64)
		f_percentage, _ := strconv.ParseFloat(strings.Trim(fields[5], "%"), 64)

		stats = append(stats, filesystemStats{
			labels: filesystemLabels{
				device:     fields[0],
				fsType:     fields[1],
				mountPoint: fields[6],
			},
			used:       f_used,
			avail:      f_avail,
			percentage: f_percentage,
			size:       f_size,
		})
	}
	return stats
}
