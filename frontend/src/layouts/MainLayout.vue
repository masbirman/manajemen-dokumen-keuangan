<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import Sidebar from "@/components/Sidebar.vue";
import Navbar from "@/components/Navbar.vue";

const breakpoint = 768;
const isMobile = ref(window.innerWidth < breakpoint);
const sidebarOpen = ref(!isMobile.value);

const handleResize = () => {
  const mobile = window.innerWidth < breakpoint;
  if (mobile !== isMobile.value) {
    isMobile.value = mobile;
    // Auto adjust sidebar based on mode
    sidebarOpen.value = !mobile;
  }
};

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value;
};

onMounted(() => {
  window.addEventListener("resize", handleResize);
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
});
</script>

<template>
  <div class="min-h-screen bg-gray-100 relative">
    <!-- Overlay for mobile -->
    <div
      v-if="isMobile && sidebarOpen"
      class="fixed inset-0 bg-black/50 z-40 transition-opacity backdrop-blur-sm"
      @click="sidebarOpen = false"
    ></div>

    <Sidebar :open="sidebarOpen" :is-mobile="isMobile" @toggle="toggleSidebar" />

    <div
      class="transition-all duration-300 min-h-screen flex flex-col"
      :class="[isMobile ? 'ml-0' : sidebarOpen ? 'ml-64' : 'ml-16']"
    >
      <Navbar @toggle-sidebar="toggleSidebar" />
      <main class="p-4 md:p-6 flex-1 w-full overflow-x-hidden">
        <RouterView />
      </main>
    </div>
  </div>
</template>
