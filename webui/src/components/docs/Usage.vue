<script setup>

</script>

<template>
  <div class="card">
    <h3>Usage</h3>
    <ul class="line-height-3">
      <li>Install T3</li>
      <li>Setup a webhook on Grafana to call T3 for alerts</li>
      <li>Enable matrix channel on T3 (to send alerts to matrix and inform you)</li>
      <li>you will talk to T3 about each alert through commands on matrix or the T3 dashboard</li>
      <li>
        T3 will capture data about the ticket:
        <ul>
          <li>The time ticket is created</li>
          <li>The time you've seen the ticket</li>
          <li>The time it's taken to solve the issue and set status of ticket as resolved</li>
          <li>The labels about the ticket (its for which team, product,...)</li>
          <li>...</li>
        </ul>
      </li>
      <li>
        you can finally check out the T3's Grafana dashboard or even
        run custom queries on its DB (e.g., through metabase) to have better insight
        about your apps incidents
      </li>
    </ul>

    <h5>Notes</h5>
    <ul class="line-height-3">
      <li>Currently the sources we support are: manual(create tickets manually) and grafana alerts</li>
      <li>Currently the communication channels that we support are T3 dashboard and matrix</li>
    </ul>

    <h3>Channels</h3>
    <h5>Matrix</h5>
    You can set matrix server config on the `chanenls.yaml` config file before

    running T3 (it will be a hot-reload config). <br>

    T3 will send a message per ticket in the specified matrix room. <br>

    You can run commands on a ticket by replying to the message's thread in the room. <br>

    For example to set a ticket as `resolved`, replay to the message of the

    ticket in your matrix room with message `!!resolved` <br>

    Available commands are:
    <pre><code>
      !!help: shows the help message.

      !!dash: returns link to the dashboard

      !!new: Creaet a new ticket by providing its yaml data

      !!new_link: get a link to the ticket creation form in the UI.

      !!patch: Patch the ticket using yaml value passed as param of the command

      !!edit_link: Get a link to the ticket edition form in the UI.

      !!ticket: Get the ticket in different formats(yaml, short, verbose). e.g., !!ticket, !!ticket short, !!ticket yaml

      !!ticket_yaml: Get the ticket in yaml format.

      !!seen: mark ticket as seen. e.g., !!seen {minutes(default: 0)

      !!spam: set the spam flag on the ticket. e.g., !!spam, !!spam false, !!spam true(default)

      !!resolved: set a ticket as resolved. e.g., !!resolved {minutes(default: 0)}

      !!firing: set a ticket as firing.

      !!level: set level of a ticket. e.g., !!level {level: low,medium or high}

      !!description: set description on a ticket. It’ll remove the previous description content if it’s not empty. e.g., !!description {msg}
    </code></pre>
    <p></p>
  </div>
</template>

<style scoped lang="scss">

</style>
