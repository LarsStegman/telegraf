# FIT Parser Plugin

The `fit` parser creates metrics from a Garmin FIT file. 
Only Activity files are supported. Other files are ignored.

The [FIT parser](github.com/tormoder/fit) is used to parse the FIT file.

## Configuration

```toml
[[inputs.file]]
  files = ["workout.fit"]

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ##   https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "fit"
 ```

## Metrics

One metric is created for each record with fields with values added as fields.
The type of the field is automatically determined based on the contents of 
the value.

In addition to the options above, you can use [metric filtering][] to skip over
columns and rows.

## Examples

Output:

```text
workout heart_rate=145,cadence=68,speed=7.8534 1536869008000000000
```

[metric filtering]: /docs/CONFIGURATION.md#metric-filtering
