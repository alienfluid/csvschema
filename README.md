# csvschema
Given a CSV file, generates a schema by sampling the data within. Sampling is performed using the reservoir sampling algorithm. For schema determination, the following order of precedence is used -

* Int32
* Int64
* Float32
* Float64
* Timestamp (various formats)
* String

## Usage
```
λ ./csvschema --help
Usage of ./csvschema:
  -delimiter string
    	Column delimiter (default ",")
  -lines int
    	Minimum number of lines to sample (unless file is smaller) (default 1000)
  -noheader
    	Don't consider the first line to be the header
```

## Examples

```
λ ./csvschema sample1.csv
Sampling 1000 records from the file
Sampled 1000 records from the file (out of total 196418)
Column Incident ID: int32
Column CR Number: int32
Column Dispatch Date / Time: string
Column Class: int32
Column Class Description: string
Column Police District Name: string
Column Block Address: string
Column City: string
Column State: string
Column Zip Code: int32
Column Agency: string
Column Place: string
Column Sector: string
Column Beat: unknown
Column PRA: int32
Column Start Date / Time: string
Column End Date / Time: string
Column Police District Number: string
Column Location: string
Column Address Number: int32

λ ./csvschema sample2.csv
Sampling 1000 records from the file
Sampled 22 records from the file (out of total 22)
Column Date: timestamp
Column Open: float32
Column High: float32
Column Low: float32
Column Close: float32
Column Adj Close: float32
Column Volume: int32
```