<script setup>

import {ref} from "vue";
import {formatTimeAgo} from "../lib/time_formatter";

const props = defineProps({
  value: { // An API ticket resource.
    type: Object,
  }
});

const isSpam = ref(props.value.is_spam)
const isFiring = ref(props.value.is_firing)
const selectedSeverity = ref(props.value.severity)
const severities = ref([
  {name: "Low"},
  {name: "Medium"},
  {name: "High"}
])

function toLocalTime(val) {
  if (!val) {
    return "-"
  }

  const now = new Date()
  const d = new Date(val)
  if (now.toLocaleDateString() === d.toLocaleDateString()) {
    return d.toLocaleTimeString()
  }

  return d.toLocaleString()
}


</script>

<template>
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

      <a target="_blank" :href="value.generator_url">
        <Button link label="Source" icon="pi pi-external-link"/>
      </a>
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
</template>

<style scoped lang="scss">

</style>
