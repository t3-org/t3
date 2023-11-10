<script setup>
import {computed, reactive, ref, watch} from 'vue';
import {formatTimeAgo, toLocalTime} from "@/lib/time_formatter";
import TicketForm from "@/components/TicketForm.vue";
import T3Service, {SeveritiesList, Source} from "@/service/T3Service";
import {useToast} from "primevue/usetoast";
import _ from "lodash"

const BlockView = reactive({
  PREVIEW: 0,
  RAW: 1,
  LABELS: 2,
  ANNOTATIONS: 3,
});

const cli = new T3Service()
const toast = useToast()

const props = defineProps({
  value: { // An API ticket resource.
    type: Object,
    required: true,
  }
});
const emit = defineEmits(["update:value"]);

// data
const isEditing = ref(false)
const selectedBlock = ref(0);
const ticket = ref({...props.value})


// methods
function activateView(event, blockViewValue) {
  selectedBlock.value = blockViewValue;
  event.preventDefault();
}

async function copyCode(event) {
  await navigator.clipboard.writeText(props.value.raw);
  event.preventDefault();
}

async function saveTicket() {
  let res = await cli.patchTicket(props.value.id, ticket.value)

  if (res.status !== 200) {
    toast.add({severity: 'error', summary: 'failed request', detail: res.body, life: 15000});
    return
  }
  toast.add({severity: 'success', summary: 'Ticket updated', life: 3000});
  emit("update:value", res.body.data)
}

// watch
watch(() => props.value, (newVal, oldVal) => {
  ticket.value = {...props.value}
}, {deep: true});


const sourceColor = computed(() => {
  return props.value.source === Source.GRAFANA ? "#ff5722" : "black"
})

const isTicketChanged = computed(() => {
  return !_.isEqual(props.value, ticket.value)
})


</script>

<template>
  <div>
    <TicketForm v-if="isEditing"
                :isEdition="true"
                v-model:keep-open="isEditing"
                :value="ticket"
                @update:value="$emit('update:value',$event)"
    />

    <div class="surface-card p-4 shadow-2 border-round" v-else>
      <!-- The header badge-->
      <div class="p-overlay-badge">
        <Badge :severity="ticket.is_firing? 'danger':'success'" :class="{'blink':ticket.is_firing}">
        </Badge>
      </div>

      <div class="header flex align-items-center justify-content-between">
        <!-- Header title-->
        <span class="inline-flex align-items-center">
          <span class="text-lg font-bold">{{ value.title }}</span>
          <Tag :style="{backgroundColor:sourceColor}" class="ml-2" :value="value.source"></Tag>
        </span>

        <div class="actions">

          <Button class="mx-1"
                  size="small"
                  :outlined="selectedBlock !== BlockView.PREVIEW"
                  @click="activateView($event, BlockView.PREVIEW)">
            Ticket
          </Button>

          <Button class="mx-1"
                  :text="selectedBlock!==BlockView.LABELS"
                  size="small"
                  @click="activateView($event, BlockView.LABELS)">
            Labels
          </Button>
          <Button class="mx-1"
                  size="small"
                  :text="selectedBlock!==BlockView.ANNOTATIONS"
                  @click="activateView($event, BlockView.ANNOTATIONS)">
            Annotations
          </Button>

          <Button class="mx-1"
                  rounded
                  icon="pi pi-code"
                  :text="selectedBlock!==BlockView.RAW"
                  @click="activateView($event, BlockView.RAW)"/>

          <Button class="mx-1"
                  text
                  rounded
                  icon="pi pi-copy"
                  aria-label="Submit"
                  v-tooltip.focus.bottom="{ value: 'Copied!' }"
                  @click="copyCode($event)"/>
        </div>
      </div>

      <div class="content">
        <div class="preview" v-if="selectedBlock === BlockView.PREVIEW">
          <p>{{ value.description }}</p>
          <div>
            <span class="text-gray-400">Values:</span>
            <ul>
              <li v-for="(val,key) in value.values">
                <span>{{ key }}</span>:
                <code>{{ val }}</code>
              </li>
            </ul>
          </div>


        </div>
        <div class="raw" v-if="selectedBlock === BlockView.RAW">
          <p>
            <span>id: <small class="text-gray-400">{{ value.id }}</small></span>
            <span class="ml-2">fingerprint: <small class="text-gray-400">{{ value.fingerprint }}</small></span>
          </p>
          <pre><code>{{ value.raw }}</code></pre>
        </div>
        <div class="labels my-6" v-if="selectedBlock === BlockView.LABELS">
          <ul>
            <li v-for="(val,key) in value.labels" class="mb-2">
              <span class="text-gray-400">{{ key }}</span>:
              <Chip><code class="p-1">{{ val }}</code></Chip>
            </li>
          </ul>
        </div>
        <div class="annotations mb-6" v-if="selectedBlock === BlockView.ANNOTATIONS">
          <ul>
            <li v-for="(val,key) in value.annotations" class="mb-2">
              <span class="text-gray-400">{{ key }}</span>:
              <span>{{ val }}</span>
            </li>
          </ul>
        </div>
      </div>
      <div class="footer">
        <div class="flex align-items-center justify-content-between mt-6">
          <div>
            <span class="inline-flex align-items-center mr-4">
              <span class="font-bold mr-1">Severity: </span>
              <Dropdown v-model="ticket.severity" :options="SeveritiesList" optionLabel="name" placeholder="Severity"/>
            </span>
            <span class="inline-flex align-items-center mr-4">
              <span class="font-bold mr-1">Is Firing: </span>
              <InputSwitch v-model="ticket.is_firing"/>
            </span>
            <span class="inline-flex align-items-center mr-4">
              <span class="font-bold mr-1">Is Spam: </span>
              <InputSwitch v-model="ticket.is_spam"/>
            </span>
          </div>

          <div>
            <Button v-if="isTicketChanged" label="Save changes" icon="pi pi-pencil" @click="saveTicket"/>
            <Button text label="Edit" icon="pi pi-pencil" @click="isEditing=true"/>
            <a target="_blank" :href="value.generator_url">
              <Button link label="Generator" icon="pi pi-external-link"/>
            </a>
          </div>
        </div>
        <Divider/>

        <div class="flex align-items-center justify-content-between">
          <div>
            <span class="text-gray-400">Fired at: </span>
            <span>{{ toLocalTime(value.started_at) }}</span>
            <span class="text-gray-500">({{ formatTimeAgo(value.started_at) }})</span>
          </div>
          <div>
            <span class="text-gray-400">Seen at: </span>
            <span>{{ toLocalTime(value.seen_at) }}</span>
            <span v-if="value.seen_at" class="text-gray-500">({{ formatTimeAgo(value.seen_at) }})</span>
          </div>
          <div>
            <span class="text-gray-400">Resolved at:</span>
            <span>{{ toLocalTime(value.ended_at) }}</span>
            <span v-if="value.ended_at" class="text-gray-500">({{ formatTimeAgo(value.ended_at) }})</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
