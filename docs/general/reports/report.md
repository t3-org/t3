We have two types of reports:

- numeric data
- chart data

We can use dedicated dashboard for it or use grafana.

### Filters

Filters are:

- time range
- a labelSet. for example we can get reports with label set `region=teh1,team=orders`.
- An optional `group` label key. e.g., group reports by `product`. In this case we'll add a list of reports grouped
  by the group label to the report too(in addition to the base report). For example we'll return a report for
  the `orders` team and also a list of reports of the `orders` team's products.

### Numeric data

- How many spams?
- How many tickets(non-spam tickets)?
- how many low,medium,high severity tickets?
- how many by each `source`?
- p90 or mean of TTSee (time to see) of tickets.
- p90 or mean of TTResolve of tickets.

### Chart data

- Panel per each report of the numeric section. with time-series per each field of it. e.g. a panel for
  the `orders` team which shows `how many tickets, how many spams,...` per month and also the same panel
  for each product in the `orders` team.

- A pie chart to show percent of each value in the `group by` condition of our filters. for example group by `product`
  will show percent of each product in pie chart.

### Challenges of Sql Queries in Grafana :

The dynamic parts are:

- The `labelSet` to set as filters of the query like `(key='a' and val='b') and (key='c' and val='d')`
- the `group` filter to set `val=?`
- The `date` filter to set `where started_at>=? and started_at<=?`.

Solutions:

- For the label set, we'll have specific label set values. like `region` and `team`.
  Per each value, we'll set a dropdown. e.g., two dropdowns: `region` and `team`.
- For group filter, we'll use a variable. it value could be one of available keys, like `product`, `reason`,...
- For the `date` we'll use the date filter of Grafana.

### Implementation

We have two options:

- Custom implementation for monitoring part. it'll contains following fields to do filtering:
    - LabelSet filter text input.
    - The `group by` text input.
    - The `N last days` and `N days offset` numeric inputs to filter by days.
- Grafana: 
  - by setting static LabelSet to filter (as multiple dropdowns to
    allow users set values)
  - A dropdown for `group by filer`
  - Date filer will be the grafana date filer.

The selected implementation option for the first phase is Grafana, because:

- It'll allow teams customize their dashboard and queries based on their needs.
- Faster to implement with a lot of options in charts.
- It'll give companies a simple and pre-specified standard labelSets (static labelSets) to filter
  instead of needing to search by custom LabelSet by each user. actually this could
  be bad and good as well, I assume it's good and makes a standard dashboard
  for companies.

