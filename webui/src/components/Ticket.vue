<script setup>
import {computed, reactive, ref} from 'vue';
import TicketContent from "@/components/TicketContent.vue";

const props = defineProps({
  value: { // An API ticket resource.
    type: Object,
  }
});

const isFiring = ref(props.value.is_firing)

const BlockView = reactive({
  PREVIEW: 0,
  RAW: 1,
  LABELS: 2,
  ANNOTATIONS: 3,
});

const blockView = ref(0);

function activateView(event, blockViewValue) {
  blockView.value = blockViewValue;
  event.preventDefault();
}

async function copyCode(event) {
  await navigator.clipboard.writeText(props.value.raw);
  event.preventDefault();
}


const SourceGrafana = "grafana"
const sourceColor = computed(() => {
  return props.value.source === SourceGrafana ? "#ff5722" : "black"
})
</script>

<template>
  <div class="surface-card p-4 shadow-2 border-round">
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
                :outlined="blockView !== BlockView.PREVIEW"
                @click="activateView($event, BlockView.PREVIEW)">
          Ticket
        </Button>

        <Button class="mx-1"
                :text="blockView!==BlockView.LABELS"
                size="small"
                @click="activateView($event, BlockView.LABELS)">
          Labels
        </Button>
        <Button class="mx-1"
                size="small"
                :text="blockView!==BlockView.ANNOTATIONS"
                @click="activateView($event, BlockView.ANNOTATIONS)">
          Annotations
        </Button>

        <Button class="mx-1"
                rounded
                icon="pi pi-code"
                :text="blockView!==BlockView.RAW"
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
    <div>
      <TicketContent :value="value" v-if="blockView === BlockView.PREVIEW"></TicketContent>
      <div v-if="blockView === BlockView.RAW">
        <p>
          <span>id: <small class="text-gray-400">{{ value.id }}</small></span>
          <span class="ml-2">fingerprint: <small class="text-gray-400">{{ value.fingerprint }}</small></span>
        </p>
        <pre><code>{{ value.raw }}</code></pre>
      </div>
      <div v-if="blockView === BlockView.LABELS">
        Labels
      </div>
      <div v-if="blockView === BlockView.ANNOTATIONS">
        Annotations
      </div>
    </div>

  </div>
</template>

<style scoped lang="scss"></style>
