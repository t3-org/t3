### Features
- [ ] Auto spam
- [ ] Post-mortem
- [ ] Auto label (will be used to set on-call user) as a plugin
- [ ] Limited labels (to keep some label values in a white list)
- [ ] [Plugin system](./features/plugin_system.md) for channels(matrix,...) and sources(grafana,...).
- [ ] Reporting dashboard
- [ ] tickets graph and merge [feature.md](features/tickets_graph_and_merge_feature.md)
- [ ] AI assistant
- [ ] Auto-refresh the channels configs peridically (e.g., every 10 minutes and also provide an api endpoint to refresh it).

### Improvements
- [ ] better message rendering for channels
  - [x] set date format like this: `Feb 14, 2023, 1:42 PM` 
  - [ ] get a local in the app config to localize dates in the channels messages
- [x] Have a filter in the tickets-search page to view single ticket by id.
- [x] sync tickets filter input with the url (to be able to share and open the tickets page with specific filters)
- [x] In the edit link of matrix plugin, add filter to the url to show that ticket on the page.
- [x] Add search-help doc on the search input.
- [ ] Dispatch ticket changes to channels as async jobs.

### spam feature

- We'll detect if a ticket should be marked as spam or not by the following rule:

- if the alert is because of `no_data` or `QueryError` on grafana.

### Limited labels
We should be able to limit some labels to have specific values.
For example the `reason` label should be limited to e.g., `db,code,human,other-.*`.
In this way we prevent from having invalid label values. e.g., `db` and `DB`...
this will make our labelSets more clean and also provide better report in our report-system.


### Auto labels
How should let users to set extra labels based on some other labels.
e.g., If a new ticket has `team=orders`, add label `oncall=mehdi` to it.
We'll use an API to be able to set the extra label's value too. in this way we'll have oncall feature too.

### AI Assistant
We can use AI to interpret the comments from users per ticket. in this way we can detect reason of firing some ticket
and put it in the ticket's details.
