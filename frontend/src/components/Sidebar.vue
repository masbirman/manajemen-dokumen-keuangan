<script setup lang="ts">
import { computed } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRoute } from "vue-router";

defineProps<{
  open: boolean;
}>();

defineEmits<{
  toggle: [];
}>();

const authStore = useAuthStore();
const route = useRoute();

interface MenuItem {
  name: string;
  path: string;
  icon: string;
  roles?: string[];
}

const menuItems = computed<MenuItem[]>(() => {
  const items: MenuItem[] = [
    { name: "Dashboard", path: "/", icon: "📊" },
    { name: "Input Dokumen", path: "/dokumen/input", icon: "📝" },
    { name: "List Dokumen", path: "/dokumen", icon: "📋" },
  ];

  if (authStore.isAdmin) {
    items.push(
      {
        name: "Unit Kerja",
        path: "/unit-kerja",
        icon: "🏢",
        roles: ["super_admin", "admin"],
      },
      {
        name: "PPTK",
        path: "/pptk",
        icon: "👤",
        roles: ["super_admin", "admin"],
      },
      {
        name: "Sumber Dana",
        path: "/sumber-dana",
        icon: "💰",
        roles: ["super_admin", "admin"],
      },
      {
        name: "Jenis Dokumen",
        path: "/jenis-dokumen",
        icon: "📁",
        roles: ["super_admin", "admin"],
      }
    );
  }

  if (authStore.isSuperAdmin) {
    items.push(
      {
        name: "Manajemen User",
        path: "/users",
        icon: "👥",
        roles: ["super_admin"],
      },
      {
        name: "Petunjuk",
        path: "/petunjuk",
        icon: "📖",
        roles: ["super_admin"],
      },
      {
        name: "Pengaturan",
        path: "/settings",
        icon: "⚙️",
        roles: ["super_admin"],
      }
    );
  }

  return items;
});

const isActive = (path: string) => {
  if (path === "/") return route.path === "/";
  return route.path.startsWith(path);
};
</script>

<template>
  <aside
    class="fixed left-0 top-0 h-full bg-blue-800 text-white transition-all duration-300 z-50"
    :class="open ? 'w-64' : 'w-16'"
  >
    <div class="flex items-center justify-between p-4 border-b border-blue-700">
      <h1 v-if="open" class="text-xl font-bold">Dokumen Keuangan</h1>
      <button
        @click="$emit('toggle')"
        class="p-2 rounded hover:bg-blue-700 transition-colors"
      >
        {{ open ? "◀" : "▶" }}
      </button>
    </div>

    <nav class="mt-4">
      <RouterLink
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center px-4 py-3 hover:bg-blue-700 transition-colors"
        :class="{ 'bg-blue-700': isActive(item.path) }"
      >
        <span class="text-xl">{{ item.icon }}</span>
        <span v-if="open" class="ml-3">{{ item.name }}</span>
      </RouterLink>
    </nav>
  </aside>
</template>
