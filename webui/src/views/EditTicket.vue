<script setup>


import TicketForm from "@/components/TicketForm.vue";
import {onMounted, ref, watch} from "vue";
import T3Service from "@/service/T3Service";
import {useRoute} from "vue-router";
import {useToast} from "primevue/usetoast";
import router from "@/router";


const cli = new T3Service()
const route = useRoute()
const toast = useToast()

const ticket = ref({})
const keepOpen = ref(true)

onMounted(async () => {
  const res = await cli.getTicket(route.params.id)
  if (res.status !== 200) {
    toast.add({
      severity: 'error',
      summary: `can not fetch ticket(status: ${res.status})`,
      detail: res.body,
      life: 15000
    });
    return
  }
  console.log(res)
  ticket.value = res.body.data
})

watch(() => keepOpen.value, () => {
  router.push({name: 'dashboard', query: router.currentRoute.value.query})
});


</script>

<template>
  <TicketForm v-if="ticket.id && keepOpen" :isEdition="true" v-model:keep-open="keepOpen" :value="ticket"/>
</template>

