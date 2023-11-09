<script setup>
import Ticket from "@/components/Ticket.vue";

import {reactive, ref} from 'vue';

const raw = `{
      "status": "firing",
      "labels": {
        "alertname": "High memory usage",
        "team": "blue",
        "zone": "us-1"
      },
      "annotations": {
        "description": "The system has high memory usage",
        "runbook_url": "https://myrunbook.com/runbook/1234",
        "summary": "This alert was triggered for zone us-1"
      },
      "startsAt": "2021-10-12T09:51:03.157076+02:00",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "https://play.grafana.org/alerting/1afz29v7z/edit",
      "fingerprint": "c6eadffa33fcdf37",
      "silenceURL": "https://play.grafana.org/alerting/silence/new?alertmanager=grafana&matchers=alertname%3DT2%2Cteam%3Dblue%2Czone%3Dus-1",
      "dashboardURL": "",
      "panelURL": "",
      "values": {
        "B": 44.23943737541908,
        "C": 1
      }
    }`

const ticket = reactive({
  raw,
  source: "grafana",
  title: "high memory usage",
  description: "the memory usage is high",
  id: "12343247823",
  fingerprint: "2e33938jfehbfubhy3fsda",
  values: {
    "B": "44.23943737541908",
    "C": "1"
  },
  is_spam: false,
  is_firing: true,
  generator_url: "https://play.grafana.org/alerting/1afz29v7z/edit",
  severity: "low",
  started_at: 1699517416399,
  seen_at: 1699517416399,
  ended_at: 1699517416399,
})

const searchTerm = ref()

async function search(event, val) {
  console.log("searching...", val)
}

</script>

<template>
  <form @submit.prevent="search($event,searchTerm)">
    <div class="p-inputgroup mb-2">
      <InputText v-model="searchTerm"
                 size="large"
                 placeholder="search in k8s label selector format (e.g., team=ordering)"/>
    </div>

  </form>

  <Ticket :value="ticket">
    <p>Content</p>
  </Ticket>
</template>

