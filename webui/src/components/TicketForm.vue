<script setup>


import {computed, reactive, ref} from "vue";
import {parse, stringify} from 'yaml'
import T3Service, {SeveritiesList} from "@/service/T3Service";

const availableSources = ref([{name: "grafana"}, {name: "manual"}])

const props = defineProps({
  isEditing: {type: Boolean},
  value: { // An API ticket resource (could be null to create a new ticket.).
    type: Object,
  },

});

const emits = defineEmits(["update:isEditing"])


// data
const ticket = reactive({
  fingerprint: props.value.fingerprint,

  source: props.value.source,
  raw: props.value.raw,
  annotations: props.value.annotations,
  is_firing: props.value.is_firing,
  started_at: props.value.started_at,
  ended_at: props.value.ended_at,
  values: props.value.values,
  generator_url: props.value.generator_url,
  is_spam: props.value.is_spam,
  severity: props.value.sevenrity,
  title: props.value.title,
  description: props.value.description,
  seen_at: props.value.seen_at,
  labels: props.value.labels,
})

const ticketOtherData = reactive({
  annotationsStr: ticket.annotations ? stringify(ticket.annotations) : "",
  labelsStr: ticket.labels ? stringify(ticket.labels) : "",
  valuesStr: ticket.values ? stringify(ticket.values) : "",
})

const startedAt = computed({
  get() {
    if (!ticket.started_at) {
      return null;
    }
    return new Date(ticket.started_at)
  },
  set(val) {
    ticket.started_at = val.getTime()
  }
})

const seenAt = computed({
  get() {
    if (!ticket.seen_at) {
      return null
    }
    return new Date(ticket.seen_at)
  },
  set(val) {
    ticket.seen_at = val.getTime()
  }
})

const endedAt = computed({
  get() {
    if (!ticket.ended_at) {
      return null
    }
    return new Date(ticket.ended_at)
  },
  set(val) {
    ticket.ended_at = val.getTime()
  }
})

const tcli = new T3Service()

async function save() {
  // update all yaml values on the ticket model:
  ticket.annotations = parse(ticketOtherData.annotationsStr)
  ticket.labels = parse(ticketOtherData.labelsStr)
  ticket.values = parse(ticketOtherData.valuesStr)

  // save the data
  if (props.value.id) {
    const res = await tcli.patchTicket(props.value.id, ticket)
    console.log("res", res)
  }

}

</script>

<template>
  <form @submit.prevent="save">
    <div class="card">
      <div class="flex  flex-column gap-4">

        <div class="flex flex-column  gap-2">
          <label class="font-bold" for="id">Id</label>
          <InputText disabled id="id" :value="value.id"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="fingerprint">Fingerprint</label>
          <InputText disabled id="fingerprint" v-model="ticket.fingerprint"/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="title">Title</label>
          <InputText disabled id="title" v-model="ticket.title"/>
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
                    :value="stringify(ticket.annotations)"
                    @update:model-value="ticketOtherData.annotationsStr=$event"
          />
        </div>


        <div class="flex flex-column gap-2">
          <label class="font-bold" for="labels">Labels(yaml format)</label>
          <Textarea id="labels" :value="stringify(ticket.labels)" auto-resize/>
        </div>

        <div class="flex flex-column gap-2">
          <label class="font-bold" for="values">Values(yaml format)</label>
          <Textarea id="values" :value="stringify(ticket.values)" auto-resize/>
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
        <Button label="Cancel" type="button" class="mx-2" @click.prevent="$emit('update:isEditing',false)"/>
        <Button label="Save" type="submit"/>
      </div>
    </div>
  </form>

</template>

<style scoped lang="scss">

</style>
