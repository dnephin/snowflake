# snowid

[![Go Reference](https://pkg.go.dev/badge/github.com/dnephin/snowid.svg)](https://pkg.go.dev/github.com/dnephin/snowid)

snowid is a [Go](https://golang.org/) package that provides
* A very simple Twitter snowflake generator.
* Methods to parse existing snowflake IDs.
* Methods to convert a snowflake ID into several other data types and back.
* JSON Marshal/Unmarshal functions to easily use snowflake IDs within a JSON API.
* Monotonic Clock calculations protect from clock drift.

## ID Format

By default, the ID format follows the original Twitter snowflake format.
* The ID as a whole is a 63 bit integer stored in an int64
* 41 bits are used to store a timestamp with millisecond precision, using a custom epoch.
* 10 bits are used to store a node id - a range from 0 through 1023.
* 12 bits are used to store a sequence number - a range from 0 through 4095.
