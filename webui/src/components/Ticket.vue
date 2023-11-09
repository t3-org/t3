<script setup>
import {computed, reactive, ref} from 'vue';
import {formatTimeAgo, toLocalTime} from "@/lib/time_formatter";
import TicketForm from "@/components/TicketForm.vue";
import {Source} from "@/service/T3Service";

const BlockView = reactive({
  PREVIEW: 0,
  RAW: 1,
  LABELS: 2,
  ANNOTATIONS: 3,
});


const props = defineProps({
  value: { // An API ticket resource.
    type: Object,
  }
});


const isEditing = ref(true)
const selectedBlock = ref(0);
const isSpam = ref(props.value.is_spam)
const isFiring = ref(props.value.is_firing)
const selectedSeverity = ref(props.value.severity)
const severities = ref([
  {name: "Low"},
  {name: "Medium"},
  {name: "High"}
])


function activateView(event, blockViewValue) {
  selectedBlock.value = blockViewValue;
  event.preventDefault();
}

async function copyCode(event) {
  await navigator.clipboard.writeText(props.value.raw);
  event.preventDefault();
}

const sourceColor = computed(() => {
  return props.value.source === Source.GRAFANA ? "#ff5722" : "black"
})
</script>

<template>
  <div>
    <TicketForm v-model:is-editing="isEditing" :value="value" v-if="isEditing"></TicketForm>
    <div class="surface-card p-4 shadow-2 border-round" v-else>
      <div class="p-overlay-badge">
        <Badge :severity="isFiring? 'danger':'success'" :class="{'blink':isFiring}">
        </Badge>
      </div>
      <div class="flex align-items-center justify-content-between">
        <!-- Header title-->
        <span class="inline-flex align-items-center">
          <span class="text-lg font-bold">{{ value.title }}</span>
          <Tag :style="{backgroundColor:sourceColor}" class="ml-2" :value="value.source"></Tag>
      </span>

        <!-- Header Actions-->
        <div>

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
                  v-tooltip.focus.bottom="{ value: 'Copied' }"
                  @click="copyCode($event)"/>
        </div>
      </div>

      <!-- Content-->
      <div v-if="selectedBlock === BlockView.PREVIEW">
        <div>
          <p>{{ value.description }}</p>
          <p>
            <span class="text-gray-400">Values:</span>
            <ul>
              <li v-for="(val,key) in value.values">
                <span>{{ key }}</span>:
                <code>{{ val }}</code>
              </li>
            </ul>
          </p>

          <div class="flex align-items-center justify-content-between mt-6">
            <div>
          <span class="inline-flex align-items-center mr-4">
          <span class="font-bold mr-1">Severity: </span>
          <Dropdown v-model="selectedSeverity" :options="severities" optionLabel="name" placeholder="Severity"/>
        </span>

              <span class="inline-flex align-items-center mr-4">
          <span class="font-bold mr-1">Is Spam: </span>
          <InputSwitch v-model="isSpam"/>
      </span>
              <span class="inline-flex align-items-center mr-4">
          <span class="font-bold mr-1">Is Firing: </span>
          <InputSwitch v-model="isFiring"/>
      </span>
            </div>

            <div>
              <Button text label="Edit" icon="pi pi-pencil" @click="isEditing=true"/>
              <a target="_blank" :href="value.generator_url">
                <Button link label="Source" icon="pi pi-external-link"/>
              </a>
            </div>
          </div>
          <Divider/>
          <div>
          </div>
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
        <div v-if="selectedBlock === BlockView.RAW">
          <p>
            <span>id: <small class="text-gray-400">{{ value.id }}</small></span>
            <span class="ml-2">fingerprint: <small class="text-gray-400">{{ value.fingerprint }}</small></span>
          </p>
          <pre><code>{{ value.raw }}</code></pre>
        </div>
        <div v-if="selectedBlock === BlockView.LABELS">
          Labels
        </div>
        <div v-if="selectedBlock === BlockView.ANNOTATIONS">
          Annotations
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped lang="scss"></style>
