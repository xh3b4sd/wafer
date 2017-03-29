# 001

The chart information provided here describes the ether price history provided
by https://etherchain.org. The CSV file is prepared using the following tools.

- [jq](https://github.com/stedolan/jq)
- [pr](https://www.gnu.org/software/coreutils/manual/coreutils.html#pr-invocation)

When the prerequisites are in place the CSV file can be created as follows.

```
curl https://etherchain.org/api/statistics/price > chart.json
cat chart.json | jq '.data | .[].time | strptime("%Y-%m-%dT%H:%M:%S.000Z")|mktime' > time.txt
cat chart.json | jq '.data | .[].usd' > usd.txt
pr -mts, time.txt usd.txt usd.txt > chart.csv
rm chart.json time.txt usd.txt
```
