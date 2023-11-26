-- count of spam tickets.
select count(*)
from tickets
where is_spam = true
  and id in
      (select ticket_id from ticket_labels where (key = 'ref_id' and val = 'A'));

-- count of non-spam tickets,
-- how many low,medium,high severity tickets,
-- MTTS and MTTR
select count(*),
       SUM(CASE WHEN severity is null THEN 1 ELSE 0 END)    unspecified_severity,
       SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END)    low_Severity,
       SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) medium_severity,
       SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END)   high_severity,
       avg(tickets.seen_at - tickets.started_at)            mtts,
       avg(tickets.ended_at - tickets.started_at)           mttr
from tickets
where is_spam = false
  and id in
      (select ticket_id from ticket_labels where (key = 'ref_id' and val = 'A'));

-- how many by source.
select count(*), source
from tickets
where is_spam = false
  and id in
      (select ticket_id from ticket_labels where (key = 'ref_id' and val = 'A'))
group by source;

---- With Group filter

-- count of spam tickets.
select count(*), ticket_labels.val as group_val
from tickets
         join ticket_labels on
            tickets.id = ticket_labels.ticket_id and
            ticket_labels.key = 'ref_id' -- the group filter's value
where is_spam = true
  and id in
      (select ticket_id from ticket_labels where (key = 'ref_id' and val = 'A'))
group by val;

-- count of non-spam tickets,
-- how many low,medium,high severity tickets,
-- MTTS and MTTR
select count(*),
       SUM(CASE WHEN severity is null THEN 1 ELSE 0 END)    unspecified_severity,
       SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END)    low_Severity,
       SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) medium_severity,
       SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END)   high_severity,
       avg(tickets.seen_at - tickets.started_at)            mtts,
       avg(tickets.ended_at - tickets.started_at)           mttr,
       ticket_labels.val as                                 group_val
from tickets
         join ticket_labels on
            tickets.id = ticket_labels.ticket_id and
            ticket_labels.key = 'ref_id' -- the group filter's value
where is_spam = false
  and id in
      (select ticket_id from ticket_labels where (key = 'ref_id' and val = 'A'))
group by val;

-- how many by source.
select count(*), source, ticket_labels.val as group_val
from tickets
         join ticket_labels on
            tickets.id = ticket_labels.ticket_id and
            ticket_labels.key = 'ref_id' -- the group filter's value
where is_spam = false
  and id in
      (select ticket_id from ticket_labels where (key = 'ref_id' and val = 'A'))
group by val, source;

