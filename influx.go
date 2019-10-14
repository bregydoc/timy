package timy

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go"
)

func (bag *InfluxBag) newInfluxMetricRow(c context.Context, influx *influxdb.Client, entry *Entry, sync bool) error {
	t, err := bag.getEventType(entry.TypeID)
	if err != nil {
		return err
	}

	r, err := bag.getEventRoot(t.RootID)
	if err != nil {
		return err
	}

	metricName := fmt.Sprintf("%s.%s", r.Name, t.Name)

	fields := map[string]interface{}{"value": entry.Value}

	for k, v := range entry.Modifiers {
		fields[k] = v
	}

	at := time.Now()

	metric := influxdb.NewRowMetric(
		fields,
		metricName,
		map[string]string{
			"id":   entry.ID,
			"root": r.ID,
			"type": t.ID,
		},
		at,
	)

	if sync {
		// entry.ID = generateNewID()
		// entry.At = at
		occ := t.Occurrences + 1
		if err = bag.UpdateEventType(c, t.ID, &EventTypeUpdate{
			Occurrences: &occ,
		}); err != nil {
			return err
		}
	} else {
		go func() {
			occ := t.Occurrences + 1
			if err = bag.UpdateEventType(c, t.ID, &EventTypeUpdate{
				Occurrences: &occ,
			}); err != nil {
				fmt.Println("error:", err.Error())
			}
		}()
	}
	_, err = influx.Write(c, timyBucketName, defaultOrgName, metric)
	return err
}
