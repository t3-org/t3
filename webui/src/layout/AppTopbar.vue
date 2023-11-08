<script setup>
import {computed, onBeforeUnmount, onMounted, ref} from 'vue';
import {useLayout} from '@/layout/composables/layout';
import {useRouter} from 'vue-router';

const {layoutConfig, onMenuToggle} = useLayout();

const outsideClickListener = ref(null);
const topbarMenuActive = ref(false);
const router = useRouter();

onMounted(() => {
  bindOutsideClickListener();
});

onBeforeUnmount(() => {
  unbindOutsideClickListener();
});

const logoUrl = computed(() => {
  return `/layout/images/${layoutConfig.darkTheme.value ? 'logo-white' : 'logo-dark'}.svg`;
});

const onTopBarMenuButton = () => {
  topbarMenuActive.value = !topbarMenuActive.value;
};
const onSettingsClick = () => {
  topbarMenuActive.value = false;
  router.push('/docs');
};
const topbarMenuClasses = computed(() => {
  return {
    'layout-topbar-menu-mobile-active': topbarMenuActive.value
  };
});

const bindOutsideClickListener = () => {
  if (!outsideClickListener.value) {
    outsideClickListener.value = (event) => {
      if (isOutsideClicked(event)) {
        topbarMenuActive.value = false;
      }
    };
    document.addEventListener('click', outsideClickListener.value);
  }
};
const unbindOutsideClickListener = () => {
  if (outsideClickListener.value) {
    document.removeEventListener('click', outsideClickListener);
    outsideClickListener.value = null;
  }
};
const isOutsideClicked = (event) => {
  if (!topbarMenuActive.value) return;

  const sidebarEl = document.querySelector('.layout-topbar-menu');
  const topbarEl = document.querySelector('.layout-topbar-menu-button');

  return !(sidebarEl.isSameNode(event.target) || sidebarEl.contains(event.target) || topbarEl.isSameNode(event.target) || topbarEl.contains(event.target));
};
</script>

<template>
  <div class="layout-topbar">
    <div class="box">
      <button class="p-link layout-menu-button layout-topbar-button" @click="onMenuToggle()">
        <i class="pi pi-bars"></i>
      </button>
    </div>

    <div class="box justify-content-center">
      <router-link to="/" class="layout-topbar-logo">
        <img :src="logoUrl" alt="logo"/>
        <span>iTRACK</span>
      </router-link>
    </div>

    <div class="box">
      <div class="layout-topbar-menu" :class="topbarMenuClasses">
        <button @click="onSettingsClick()" class="p-link layout-topbar-button">
          <i class="pi pi-book"></i>
          <span>Settings</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
