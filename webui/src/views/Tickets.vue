<script setup>
import Ticket from "@/components/Ticket.vue";

import {onMounted, ref} from 'vue';
import T3Service from "@/service/T3Service";

const cli = new T3Service()

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
const tickets = ref([])

// initialization
onMounted(async () => {
  const res = await cli.fetchTickets(searchTerm, 1)
  if (res.status !== 200) {
    console.log(res)
    alert('we got error on fetching tickets')
    return
  }

  tickets.value = res.body.data.items
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

  <div>
    <Ticket v-for="ticket in tickets" :key="ticket.id" class="my-3" :value="ticket"/>
  </div>

  <Paginator :rows="10" :total-records="100"
             template="PrevPageLink CurrentPageReport NextPageLink"
             currentPageReportTemplate="{first} to {last} of (page: {currentPage})"
  />


</template>

