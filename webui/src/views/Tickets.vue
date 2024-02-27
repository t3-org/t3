<script setup>
import Ticket from "@/components/Ticket.vue";

import {computed, onMounted, reactive, ref, watch} from 'vue';
import T3Service from "@/service/T3Service";
import {useToast} from "primevue/usetoast";
import {useRouter} from "vue-router";
import QueryHelp from "../components/docs/QueryHelp.vue";

const cli = new T3Service()
const toast = useToast()
const router = useRouter();

// initialization
onMounted(search)

const showHelp = ref(false)
const tickets = reactive({
  items: [],
  ignorePaginationChanges: false,
  pagination: {
    page: 1,
    per_page: 15,
    page_count: 0,
    total_count: 0,
  }
})


async function search() {
  if (tickets.ignorePaginationChanges) {
    tickets.ignorePaginationChanges = false
    return
  }
  const res = await cli.queryTickets(
      query.value,
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

/* uncomment this line to search on changing the query value
// and also use the specified debounce function for the search input.
// watch(() => query, search);
// e.g.,

import _ from "lodash";

const debounceQuery = _.debounce(function (e) {
  query.value = e.target.value;
}, 1000)

      <InputText :value="query"
                 v-on:input="debounceQuery"
                 size="large"
                 placeholder="search in k8s label selector format (e.g., team=ordering)"/>
**/
watch(() => tickets.pagination, search, {deep: true});

// Bind query param 'q' to the 'query' variable.
const query = computed({
  get: function () {
    // set to empty value if it's undefined
    if (router.currentRoute.value.query.q === undefined) {
      return ""
    }
    return router.currentRoute.value.query.q // e.g. abc.com?q=...
  },
  set(val) {
    // using replace to do not pushing a new page to the history stack.
    router.replace({
      query: {
        ...router.currentRoute.value.query,
        q: val
      }
    })
  }
})


</script>

<template>

  <QueryHelp v-show="showHelp"></QueryHelp>
  <form @submit.prevent="search">
    <div class="p-inputgroup p-overlay-badge mb-2">
      <Badge value="?" class="border-circle cursor-pointer" @click="showHelp=!showHelp"></Badge>
      <InputText v-model="query" class="shadow-none"
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

