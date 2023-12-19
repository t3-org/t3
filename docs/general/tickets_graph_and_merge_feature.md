#### Merging tickets:
Sometimes we get tickets repeatedly from the same panel for a short time-window. this is usually because of a the same issue
which is causing to fire single panel multiple times in a short time.
We should have a feature on T3 to let user merge these tickets automatically. 
For example we can say if in a 1-hour time-window we got multiple tickets from a panel, merge all tickets to the first one.



### Set incident tickets:
An incident may have multiple tickets in both the product which is cause of the incident and also all other 
products that are dependant to this product, 
The best way to detect the original ticket and ignore all other tickets I think is ignoring all tickets and just let
user to mark a ticket as an incident. by this we don't not have concern of having multiple tickets for single incident.
In addition to that, incident tickets are way less than non-incident tickets, so I think we can mark all tickets as 
spam and let users mark a ticket as non-spam (or doing another action which has the same result, I mean affect of do not
marking all tickets as incident when cause is just one product).

### Tickets graph
When there's an incident, we may have firing tickets on multiple products, if we could link all these tickets to the parent
ticket, it would be great, because in this way, we can check effect of an incident on other products and also behaviour of
the team in that product for that ticket.
The final result should be the `parent_id` field of each ticket which is filled-in with the ticket which is reason for this ticket.

To do this, I have a solution:
We should have a simple graph of products which shows which product is dependant to which product.
When we have an incident on a product, the oncall user on that product can `advertise` that ticket, meaning can 
notify all other products that are children of this ticket's product for a time-range (let's say 24 hours). 
In t3 we'll show the recent tickets((let's say tickets of last hour)) of children products to the user (maybe oncall of that team)
to approve that all (or subset) of these products are because of the advertised ticket. in this way we'll create the graph.


