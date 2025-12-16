<script setup lang="ts">
import { ref } from "vue";
import Sidebar from "@/components/Sidebar.vue";
import Navbar from "@/components/Navbar.vue";

const sidebarOpen = ref(true);

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value;
};
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <Sidebar :open="sidebarOpen" @toggle="toggleSidebar" />
    <div
      class="transition-all duration-300"
      :class="sidebarOpen ? 'ml-64' : 'ml-16'"
    >
      <Navbar @toggle-sidebar="toggleSidebar" />
      <main class="p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>
