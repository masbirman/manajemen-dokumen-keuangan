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

  if (authStore.isSuperAdmin) {
    items.push(
      {
        name: "Unit Kerja",
        path: "/unit-kerja",
        icon: "🏢",
        roles: ["super_admin"],
      },
      {
        name: "PPTK",
        path: "/pptk",
        icon: "👤",
        roles: ["super_admin"],
      },
      {
        name: "Sumber Dana",
        path: "/sumber-dana",
        icon: "💰",
        roles: ["super_admin"],
      },
      {
        name: "Jenis Dokumen",
        path: "/jenis-dokumen",
        icon: "📁",
        roles: ["super_admin"],
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
        name: "Pengaturan Login",
        path: "/login-settings",
        icon: "🔐",
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
    class="fixed left-0 top-0 h-full bg-slate-900 text-slate-100 transition-all duration-300 z-50 shadow-xl"
    :class="open ? 'w-64' : 'w-16'"
  >
    <div class="flex items-center justify-between p-4 border-b border-slate-700/50">
      <h1 v-if="open" class="text-xl font-bold tracking-tight text-white">Dokumen Keuangan</h1>
      <button
        @click="$emit('toggle')"
        class="p-2 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white transition-all"
      >
        <component :is="open ? 'span' : 'span'">{{ open ? "◀" : "▶" }}</component>
      </button>
    </div>

    <nav class="mt-4 px-2 space-y-1">
      <RouterLink
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center px-3 py-2.5 rounded-lg transition-all duration-200 group"
        :class="isActive(item.path) ? 'bg-blue-600 text-white shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
      >
        <span class="text-xl group-hover:scale-110 transition-transform duration-200">{{ item.icon }}</span>
        <span v-if="open" class="ml-3 font-medium">{{ item.name }}</span>
      </RouterLink>
    </nav>
  </aside>
</template>
