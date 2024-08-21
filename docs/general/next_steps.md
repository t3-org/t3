### Features
- [ ] Auto spam detection
- [ ] Post-mortem generator
- [ ] Auto label (will be used to set the on-call user) as a plugin
- [ ] limited label values (to keep some label values in a white list) (e.g., the `team` label should be in the predefined list)
- [ ] [Plugin system](./features/plugin_system.md) for channels(matrix,...) and sources(grafana,...).
- [ ] Reporting dashboard (currently we have the Grafana dashboard)
- [ ] tickets graph and merge [feature.md](features/tickets_graph_and_merge_feature.md)
- [ ] AI assistant
- [ ] Hot reloading feature for the channels configs (e.g., every 10 minutes read 
     from a git repo and also provide an API endpoint to refresh it).

### Improvements
- [ ] better message rendering for channels
  - [x] set date format like this: `Feb 14, 2023, 1:42 PM` 
  - [ ] Localize dates in the messages
- [x] Add an ID filter in the search page to view a ticket by its ID.
- [x] Add filters of the search page to the page's URL on change (to be able to share and open the tickets page with specific filters)
- [x] In the edit link of the matrix plugin, add the ticket's ID filter to the URL to let users see that specific ticket when they
      open the edit page
- [x] Add search-help doc to the search input.
- [ ] Dispatch ticket changes to channels as an async job.

### spam feature

- We'll detect if a ticket should be marked as spam by the following rule:

- if the alert is because of `no_data` or `QueryError` on Grafana.

### pre-defined values for labels
We should be able to limit the values of some labels.
For example the `reason` label should be limited to e.g., `db,code,human,other-.*`.
In this way, we prevent invalid label values. e.g., `db` and `DB`...
this will make our labelSets more clean and also provide better reports in our report system.


### Auto labels
How should we let users set extra labels based on some other labels?
e.g., If a new ticket has `team=orders`, add the label `oncall=mehdi` to it.
We'll use an API to set the extra label's value. in this way, we can provide the on-call feature as well.

### AI Assistant
We can use AI to interpret the comments from users per ticket. in this way, we can detect the reason for a firing ticket
and add it to the ticket's details.
