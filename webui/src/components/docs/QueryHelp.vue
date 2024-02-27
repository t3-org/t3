<script setup>

</script>

<template>
  <div class="card">
    <h3>Query Documentation</h3>
    <p>
      The query could be written in k8s selectors format. e.g., label_a=val_b,... <br>
      We have two category of fields that can be used in the query.
    </p>
    <ul class="line-height-3">
      <li>the ticket fields, prefixed with
        <Tag>f.</Tag>
        . e.g., <i class="code">f.id=abdcd</i></li>
      <li>the ticket labels. e.g., <i class="code">abc=def</i></li>
    </ul>
    <Divider/>

    <h5>Ticket avaiable fields to query</h5>
    <ul class="line-height-4">
      <li><span class="text-primary font-medium">f.id (text)</span></li>
      <li><span class="text-primary font-medium">f.fingerpint (text)</span></li>
      <li><span class="text-primary font-medium">f.source (text = grafana|manual|...)</span></li>
      <li><span class="text-primary font-medium">f.is_firing (bool = true|false)</span></li>
      <li><span class="text-primary font-medium">f.started_at (date, format: )</span></li>
      <li><span class="text-primary font-medium">f.ended_at (date, format: 2006-01-02 15:04:05 )</span></li>
      <li><span class="text-primary font-medium">f.is_spam(bool= true|false)</span></li>
      <li><span class="text-primary font-medium">f.severity(enum = low, medium, high)</span></li>
      <li><span class="text-primary font-medium">f.seen_at (date, format: 2006-01-02 15:04:05)</span></li>
      <li><span class="text-primary font-medium">f.seen_at (date, format: 2006-01-02 15:04:05)</span></li>
    </ul>

    In addition to the ticket fields, you can use any label on the tickets too.

    <Divider/>
    <h5>Examples</h5>
    <ul class="line-height-3">
      <li class="mb-1">select ticket with id=abc: <i class="code">f.id=abc</i></li>
      <li>Select spam tickets: <i class="code">f.is_spam=true</i></li>
      <li>Select high severity tickets that are not spam: <i class="code">f.is_spam=false,severity=high</i></li>
      <li>
        Select high severity tickets that are not spam for the "orders" product:
        <i class="code">f.is_spam=false,severity=high,product=orders</i>
      </li>
    </ul>
  </div>
</template>




