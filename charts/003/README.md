# 002

The chart information provided here describes the bitcoin price history provided
by http://api.bitcoincharts.com. The CSV file is prepared using the following
tools.

- [wget](https://www.gnu.org/software/wget/manual/wget.html)
- [gzip](https://linux.die.net/man/1/gzip)

When the prerequisites are in place the CSV file can be created as follows. Note
that the downloaded file is super big (~500M), so we only take the last 15000
lines of it.

```
wget http://api.bitcoincharts.com/v1/csv/coinbaseUSD.csv.gz
gzip -d coinbaseUSD.csv.gz
tail -n 15000 coinbaseUSD.csv > chart.csv
rm coinbaseUSD.csv
```
