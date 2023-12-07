<script setup>
import Ticket from "@/components/Ticket.vue";

import {onMounted, reactive, watch} from 'vue';
import T3Service from "@/service/T3Service";
import {useToast} from "primevue/usetoast";
import _ from "lodash";

const cli = new T3Service()
const toast = useToast()

// initialization
onMounted(search)


const tickets = reactive({
  items: [],
  query: "",
  ignorePaginationChanges: false,
  pagination: {
    page: 1,
    per_page: 15,
    page_count: 0,
    total_count: 0,
  }
})


async function search() {
  if (tickets.ignorePaginationChanges){
    tickets.ignorePaginationChanges=false
    return
  }
  const res = await cli.queryTickets(
      tickets.query,
      tickets.pagination.page,
      tickets.pagination.per_page,
  )

  if (res.status !== 200) {
    toast.add({severity: 'error', summary: 'can not fetch tickets', detail: res.body, life: 15000});
    return
  }

  tickets.items = res.body.data.items
  tickets.ignorePaginationChanges = true
  tickets.pagination = res.body.data.pagination
}

function updateTicketById(ticket) {
  for (let i = 0; i < tickets.items.length; i++) {
    if (tickets.items[i].id === ticket.id) {
      tickets.items[i] = ticket
    }
  }
}

const debounceQuery = _.debounce(function (e) {
  tickets.query = e.target.value;
}, 1000)

// watch(() => tickets.query, search); // uncomment this line to search on changing the query value.
watch(() => tickets.pagination, search, {deep: true});


</script>

<template>
  <form @submit.prevent="search">
    <div class="p-inputgroup mb-2">
      <InputText :value="tickets.query"
                 v-on:input="debounceQuery"
                 size="large"
                 placeholder="search in k8s label selector format (e.g., team=ordering)"/>
    </div>

  </form>

  <div>
    <Ticket v-for="ticket in tickets.items"
            :key="ticket.id"
            class="my-3"
            :value="ticket"
            @update:value="updateTicketById($event)"
    />
  </div>

  <Paginator
      @page="(e)=>tickets.pagination.page=e.page"
      :rows="tickets.pagination.per_page"
      :total-records="tickets.pagination.total_count"
      template="PrevPageLink CurrentPageReport NextPageLink"
      currentPageReportTemplate="{first} to {last} of {totalRecords} (page: {currentPage})"
  />


</template>

