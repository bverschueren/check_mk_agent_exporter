package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

var (
	diskstatLabelNames = []string{"major_number", "minor_number", "device_name"}
)

type diskstatCollector struct {
	ReadsCompletedDesc           *prometheus.Desc
	ReadsMergedDesc              *prometheus.Desc
	SectorsReadDesc              *prometheus.Desc
	TimeSpentReadingDesc         *prometheus.Desc
	WritesCompletedDesc          *prometheus.Desc
	WritesMergedDesc             *prometheus.Desc
	SectorsWrittenDesc           *prometheus.Desc
	TimeSpentWritingDesc         *prometheus.Desc
	IOCurrentlyDesc              *prometheus.Desc
	TimeSpentDoingIODesc         *prometheus.Desc
	WeightedTimeSpentDoingIODesc *prometheus.Desc
	DiscardsCompletedDesc        *prometheus.Desc
	DiscardsMergedDesc           *prometheus.Desc
	SectorsDiscardedDesc         *prometheus.Desc
	TimeSpentDiscardingDesc      *prometheus.Desc
}

type diskLabels struct {
	major_number, minor_number, device_name string
}

type diskstat struct {
	reads_completed, reads_merged, sectors_read, time_spent_reading, writes_completed, writes_merged, sectors_written, time_spent_writing, io_currently, time_spent_doing_io, weighted_time_spent_doing_io, discards_completed, discards_merged, sectors_discarded, time_spent_discarding float64
	labels                                                                                                                                                                                                                                                                                diskLabels
}

func init() {
	registerCollector("diskstat", NewdiskstatCollector)
}

func NewdiskstatCollector() (Collector, error) {
	subsystem := "diskstat"

	ReadsCompletedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "reads_completed_successfully"),
		"reads completed successfully",
		diskstatLabelNames, nil,
	)
	ReadsMergedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "reads_merged"),
		"reads merged",
		diskstatLabelNames, nil,
	)
	SectorsReadDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "sectors_read"),
		"sectors read",
		diskstatLabelNames, nil,
	)
	TimeSpentReadingDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "time_spent_reading"),
		"time spent reading (ms)",
		diskstatLabelNames, nil,
	)
	WritesCompletedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "writes_completed_successfully"),
		"Writes completed successfully",
		diskstatLabelNames, nil,
	)
	WritesMergedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "writes_merged"),
		"writes merged",
		diskstatLabelNames, nil,
	)
	SectorsWrittenDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "sectors_written"),
		"sectors written",
		diskstatLabelNames, nil,
	)
	TimeSpentWritingDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "time_spent_writing"),
		"time spent writing (ms)",
		diskstatLabelNames, nil,
	)
	IOCurrently := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "io_currently_in_progress"),
		"I/Os currently in progress",
		diskstatLabelNames, nil,
	)
	TimeSpentDoingIO := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "time_spent_doing_io"),
		"time spent doing I/Os (ms)",
		diskstatLabelNames, nil,
	)
	WeightedTimeSpentDoingIO := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "weighted_time_spent_doing_io"),
		"weighted time spent doing I/Os (ms)",
		diskstatLabelNames, nil,
	)
	DiscardsCompletedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "discards_completed_successfully"),
		"discards completed successfully",
		diskstatLabelNames, nil,
	)
	DiscardsMergedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "discards_merged"),
		"discards merged",
		diskstatLabelNames, nil,
	)
	SectorsDiscardedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "sectors_discarded"),
		"sectors discarded",
		diskstatLabelNames, nil,
	)
	TimeSpentDiscardingDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "time_spent_discarding"),
		"time spent discarding",
		diskstatLabelNames, nil,
	)
	return diskstatCollector{
		ReadsCompletedDesc:           ReadsCompletedDesc,
		ReadsMergedDesc:              ReadsMergedDesc,
		SectorsReadDesc:              SectorsReadDesc,
		TimeSpentReadingDesc:         TimeSpentReadingDesc,
		WritesCompletedDesc:          WritesCompletedDesc,
		WritesMergedDesc:             WritesMergedDesc,
		SectorsWrittenDesc:           SectorsWrittenDesc,
		TimeSpentWritingDesc:         TimeSpentWritingDesc,
		IOCurrentlyDesc:              IOCurrently,
		TimeSpentDoingIODesc:         TimeSpentDoingIO,
		WeightedTimeSpentDoingIODesc: WeightedTimeSpentDoingIO,
		DiscardsCompletedDesc:        DiscardsCompletedDesc,
		DiscardsMergedDesc:           DiscardsMergedDesc,
		SectorsDiscardedDesc:         SectorsDiscardedDesc,
		TimeSpentDiscardingDesc:      TimeSpentDiscardingDesc,
	}, nil

}

func (d diskstatCollector) Update(unparsedStats *[]string, ch chan<- prometheus.Metric) error {
	stats := d.parseStats(unparsedStats)

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			d.ReadsCompletedDesc, prometheus.GaugeValue,
			s.reads_completed, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.ReadsMergedDesc, prometheus.GaugeValue,
			s.reads_merged, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.SectorsReadDesc, prometheus.GaugeValue,
			s.sectors_read, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.TimeSpentReadingDesc, prometheus.GaugeValue,
			s.time_spent_reading, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.WritesCompletedDesc, prometheus.GaugeValue,
			s.writes_completed, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.WritesMergedDesc, prometheus.GaugeValue,
			s.writes_merged, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.SectorsWrittenDesc, prometheus.GaugeValue,
			s.sectors_written, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.TimeSpentWritingDesc, prometheus.GaugeValue,
			s.writes_merged, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.IOCurrentlyDesc, prometheus.GaugeValue,
			s.io_currently, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.TimeSpentDoingIODesc, prometheus.GaugeValue,
			s.time_spent_doing_io, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.WeightedTimeSpentDoingIODesc, prometheus.GaugeValue,
			s.weighted_time_spent_doing_io, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.DiscardsCompletedDesc, prometheus.GaugeValue,
			s.discards_completed, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.DiscardsMergedDesc, prometheus.GaugeValue,
			s.discards_merged, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.SectorsDiscardedDesc, prometheus.GaugeValue,
			s.sectors_discarded, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
		ch <- prometheus.MustNewConstMetric(
			d.TimeSpentDiscardingDesc, prometheus.GaugeValue,
			s.time_spent_discarding, s.labels.major_number, s.labels.minor_number, s.labels.device_name,
		)
	}
	return nil
}

func (c diskstatCollector) parseStats(unparsedStats *[]string) []diskstat {

	stats := []diskstat{}
	for _, stat := range (*unparsedStats)[1:] {

		fields := strings.Fields(stat)
		if len(fields) < 14 {
			// for now only use io-related stats
			continue
		}
		reads_completed, _ := strconv.ParseFloat(fields[3], 64)
		reads_merged, _ := strconv.ParseFloat(fields[4], 64)
		sectors_read, _ := strconv.ParseFloat(fields[5], 64)
		time_spent_reading, _ := strconv.ParseFloat(fields[6], 64)
		writes_completed, _ := strconv.ParseFloat(fields[7], 64)
		writes_merged, _ := strconv.ParseFloat(fields[8], 64)
		sectors_written, _ := strconv.ParseFloat(fields[9], 64)
		time_spent_writing, _ := strconv.ParseFloat(fields[10], 64)
		io_currently, _ := strconv.ParseFloat(fields[11], 64)
		time_spent_doing_io, _ := strconv.ParseFloat(fields[12], 64)
		weighted_time_spent_doing_io, _ := strconv.ParseFloat(fields[13], 64)
		var discards_completed, discards_merged, sectors_discarded, time_spent_discarding float64
		if len(fields) > 14 {
			// diskstats 4.18+
			discards_completed, _ = strconv.ParseFloat(fields[14], 64)
			discards_merged, _ = strconv.ParseFloat(fields[15], 64)
			sectors_discarded, _ = strconv.ParseFloat(fields[16], 64)
			time_spent_discarding, _ = strconv.ParseFloat(fields[17], 64)
		}

		stats = append(stats, diskstat{
			labels: diskLabels{
				major_number: fields[0],
				minor_number: fields[1],
				device_name:  fields[2],
			},
			reads_completed:              reads_completed,
			reads_merged:                 reads_merged,
			sectors_read:                 sectors_read,
			time_spent_reading:           time_spent_reading,
			writes_completed:             writes_completed,
			writes_merged:                writes_merged,
			sectors_written:              sectors_written,
			time_spent_writing:           time_spent_writing,
			io_currently:                 io_currently,
			time_spent_doing_io:          time_spent_doing_io,
			weighted_time_spent_doing_io: weighted_time_spent_doing_io,
			discards_completed:           discards_completed,
			discards_merged:              discards_merged,
			sectors_discarded:            sectors_discarded,
			time_spent_discarding:        time_spent_discarding,
		})
	}
	return stats
}
