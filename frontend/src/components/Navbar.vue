<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";

defineEmits<{
  toggleSidebar: [];
}>();

const authStore = useAuthStore();
const router = useRouter();
const dropdownOpen = ref(false);

const roleLabels: Record<string, string> = {
  super_admin: "Super Admin",
  admin: "Admin",
  operator: "Operator",
};

const avatarUrl = computed(() => {
  const path = authStore.user?.avatar_path;
  if (!path) return "";
  return `${
    import.meta.env.VITE_API_URL || "http://localhost:8000/api"
  }/files/${path}`;
});

const logout = async () => {
  await authStore.logout();
  router.push("/login");
};
</script>

<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="flex items-center justify-between px-6 py-4">
      <div class="flex items-center">
        <button
          @click="$emit('toggleSidebar')"
          class="p-2 rounded-lg hover:bg-gray-100 md:hidden"
        >
          ☰
        </button>
        <h2 class="text-lg font-semibold text-gray-800 ml-2">
          Sistem Manajemen Dokumen Keuangan
        </h2>
      </div>

      <div class="relative">
        <button
          @click="dropdownOpen = !dropdownOpen"
          class="flex items-center space-x-3 p-2 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center overflow-hidden"
            :class="avatarUrl ? '' : 'bg-blue-600 text-white font-semibold'"
          >
            <img
              v-if="avatarUrl"
              :src="avatarUrl"
              class="w-full h-full object-cover"
            />
            <span v-else>{{
              authStore.user?.name?.charAt(0)?.toUpperCase() || "U"
            }}</span>
          </div>
          <div class="text-left hidden sm:block">
            <p class="text-sm font-medium text-gray-700">
              {{ authStore.user?.name }}
            </p>
            <p class="text-xs text-gray-500">
              {{ roleLabels[authStore.user?.role || ""] }}
            </p>
          </div>
          <span class="text-gray-400">▼</span>
        </button>

        <div
          v-if="dropdownOpen"
          class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-50"
          @click="dropdownOpen = false"
        >
          <div
            class="px-4 py-2 border-b border-gray-100 flex items-center gap-3"
          >
            <div
              class="w-10 h-10 rounded-full flex items-center justify-center overflow-hidden"
              :class="avatarUrl ? '' : 'bg-blue-600 text-white font-semibold'"
            >
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                class="w-full h-full object-cover"
              />
              <span v-else>{{
                authStore.user?.name?.charAt(0)?.toUpperCase() || "U"
              }}</span>
            </div>
            <div>
              <p class="text-sm font-medium text-gray-700">
                {{ authStore.user?.name }}
              </p>
              <p class="text-xs text-gray-500">
                {{ authStore.user?.username }}
              </p>
            </div>
          </div>
          <button
            @click="logout"
            class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
          >
            🚪 Logout
          </button>
        </div>
      </div>
    </div>
  </header>
</template>
