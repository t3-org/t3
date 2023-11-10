<script setup>


import {computed, reactive, ref} from "vue";
import {parse, stringify} from 'yaml'
import T3Service, {SeveritiesList, ticketZeroVal} from "@/service/T3Service";
import {useToast} from "primevue/usetoast";

const availableSources = ref([{name: "grafana"}, {name: "manual"}])

const toast = useToast();

const props = defineProps({
  isEdition: {type: Boolean, default: false}, // the component is using to edit ticket or create a new one?
  keepOpen: {type: Boolean}, // keep the component open or not?
  value: { // An API ticket resource (could be null to create a new ticket.).
    type: Object,
    default: {...ticketZeroVal},
  },
})

const emit = defineEmits(["update:keepOpen", "update:value"])


// data
// If we do not do shallow copy, it'll reflect changes to the prop when we change the ticket.
// so set can not use `reactive(props.value)`
const ticket = ref({...props.value})


const ticketExtra = reactive({
  annotationsStr: ticket.value.annotations ? stringify(ticket.value.annotations) : "",
  labelsStr: ticket.value.labels ? stringify(ticket.value.labels) : "",
  valuesStr: ticket.value.values ? stringify(ticket.value.values) : "",
})

const cli = new T3Service()

async function save() {
  // update all yaml values on the ticket model:
  ticket.value.annotations = parse(ticketExtra.annotationsStr)
  ticket.value.labels = parse(ticketExtra.labelsStr)
  ticket.value.values = parse(ticketExtra.valuesStr)

  // save the data
  const res = props.isEdition ?
      await cli.patchTicket(props.value.id, ticket.value) :
      await cli.createTicket(ticket.value)

  if (res.status !== 200) {
    toast.add({severity: 'error', summary: 'failed request', detail: res.body, life: 15000});
    return
  }
  toast.add({severity: 'success', summary: 'data saved', life: 3000});
  emit("update:value", res.body.data)
  emit("update:keepOpen", false)

  // if it's creation, we should reset the ticket to create a new one again:
  if (!props.isEdition) {
    ticket.value = {...ticketZeroVal}
    ticketExtra.annotationsStr = ""
    ticketExtra.labelsStr = ""
    ticketExtra.valuesStr = ""
  }
}

const startedAt = computed({
  get() {
    if (!ticket.value.started_at) {
      return null;
    }
    return new Date(ticket.value.started_at)
  },
  set(val) {
    ticket.value.started_at = val.getTime()
  }
})

const seenAt = computed({
  get() {
    if (!ticket.value.seen_at) {
      return null
    }
    return new Date(ticket.value.seen_at)
  },
  set(val) {
    ticket.value.seen_at = val.getTime()
  }
})

const endedAt = computed({
  get() {
    if (!ticket.value.ended_at) {
      return null
    }
    return new Date(ticket.value.ended_at)
  },
  set(val) {
    ticket.value.ended_at = val.getTime()
  }
})


</script>

<template>
  <form @submit.prevent="save">
    <div class="card">
      <h4 class="text-center">{{ isEdition ? `Edit Ticket` : "New Ticket" }}</h4>
      <div class="flex  flex-column gap-4">

        <div v-if="isEdition" class="flex flex-column  gap-2">
          <label class="font-bold" for="id">Id</label>
          <InputText disabled id="id" :value="value.id"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="fingerprint">Fingerprint</label>
          <InputText :disabled="isEdition" id="fingerprint" v-model="ticket.fingerprint"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="title">Title</label>
          <InputText id="title" v-model="ticket.title"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="source">Source</label>
          <Dropdown id="source" v-model="ticket.source"
                    :options="availableSources"
                    optionLabel="name"
                    option-value="name"
                    placeholder="Select a Source"
                    class="w-full md:w-14rem"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="description">Description</label>
          <Textarea id="description" v-model="ticket.description" auto-resize/>
        </div>


        <div class="flex flex-column gap-2">
          <label class="font-bold" for="raw">Raw Source Ticket content</label>
          <Textarea id="raw" v-model="ticket.raw" auto-resize/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="generator_url">Generator url</label>
          <InputText id="generator_url" v-model="ticket.generator_url"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="annotations">Annotations(yaml format)</label>
          <Textarea id="annotations"
                    v-model="ticketExtra.annotationsStr"
          />
        </div>


        <div class="flex flex-column gap-2">
          <label class="font-bold" for="labels">Labels(yaml format)</label>
          <Textarea id="labels" v-model="ticketExtra.labelsStr" auto-resize/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="values">Values(yaml format)</label>
          <Textarea id="values" v-model="ticketExtra.valuesStr" auto-resize/>
        </div>


        <div class="flex flex-column gap-2">
          <label class="font-bold" for="started_at">Started At</label>
          <Calendar id="started_at" v-model="startedAt" showTime hourFormat="24"/>
        </div>
        <div class="flex flex-column gap-2">
          <label class="font-bold" for="seen_at">Seen At</label>
          <Calendar id="seen_at" v-model="seenAt" showTime hourFormat="24"/>
        </div>
        <div class="flex flex-column gap-2">
          <label class="font-bold" for="ended_at">Ended At</label>
          <Calendar id="ended_at" v-model="endedAt" showTime hourFormat="24"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="severity">Severity</label>
          <Dropdown id="severity" v-model="ticket.severity"
                    :options="SeveritiesList"
                    optionLabel="name"
                    optionValue="value"
                    placeholder="Severity"
                    class="w-full md:w-14rem"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="is_firing">Is Firing</label>
          <InputSwitch id="is_firing" v-model="ticket.is_firing"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="ended_at">Is Spam</label>
          <InputSwitch id="is_spam" v-model="ticket.is_spam"/>
        </div>


      </div>
      <div class="flex align-items-center justify-content-end mt-4 w-full">
        <Button v-if="isEdition" label="Close" type="button" class="mx-2"
                @click.prevent="$emit('update:keepOpen',false)"/>

        <Button :label="isEdition?'Save':'save and create a new one'" type="submit"/>
      </div>
    </div>
  </form>

</template>

<style scoped lang="scss">

</style>
