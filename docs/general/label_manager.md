We should be able to limit some labels to have specific values.
For example the `reason` label should be limited to e.g., `db,code,human,other-.*`.
In this way we prevent from having invalid label values. e.g., `db` and `DB`...
this will make our labelSets more clean and also provide better report in our report-system.
