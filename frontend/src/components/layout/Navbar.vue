<template>
  <header class="h-16 bg-white border-b border-secondary-200 flex items-center justify-between px-4 md:px-6">
    <!-- Left side -->
    <div class="flex items-center gap-4">
      <!-- Mobile hamburger menu -->
      <button
        @click="$emit('toggle-mobile-menu')"
        class="lg:hidden p-2 rounded-lg hover:bg-secondary-100 transition-colors"
      >
        <svg class="w-5 h-5 text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <!-- Desktop sidebar toggle -->
      <button
        @click="$emit('toggle-sidebar')"
        class="hidden lg:block p-2 rounded-lg hover:bg-secondary-100 transition-colors"
      >
        <svg class="w-5 h-5 text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
    </div>

    <!-- Right side -->
    <div class="flex items-center gap-2 md:gap-4">
      <!-- Tahun Anggaran badge -->
      <span class="px-3 py-1 rounded-full text-xs font-medium bg-amber-100 text-amber-700 hidden md:inline-flex items-center gap-1">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        Tahun Anggaran {{ authStore.selectedYear }}
      </span>
      
      <!-- Role badge - hidden on mobile -->
      <span :class="[roleBadgeClass, 'hidden md:inline-flex']">
        {{ roleLabel }}
      </span>

      <!-- User dropdown -->
      <div class="relative" ref="dropdownRef">
        <button
          @click="showDropdown = !showDropdown"
          class="flex items-center gap-2 p-2 rounded-lg hover:bg-secondary-100 transition-colors"
        >
          <!-- Avatar with image or initials -->
          <div v-if="avatarUrl" class="w-8 h-8 rounded-full overflow-hidden">
            <img :src="avatarUrl" alt="Avatar" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-8 h-8 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white font-semibold text-sm">
            {{ authStore.user?.name?.charAt(0).toUpperCase() || 'U' }}
          </div>
          <svg class="w-4 h-4 text-secondary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        <!-- Dropdown menu -->
        <transition name="dropdown">
          <div
            v-if="showDropdown"
            class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-secondary-200 py-1 z-50"
          >
            <div class="px-4 py-2 border-b border-secondary-100">
              <p class="text-sm font-medium text-secondary-900">{{ authStore.user?.name }}</p>
              <p class="text-xs text-secondary-500">{{ (authStore.user as any)?.email || '' }}</p>
            </div>
            <router-link
              to="/profile"
              @click="showDropdown = false"
              class="w-full text-left px-4 py-2 text-sm text-secondary-700 hover:bg-secondary-50 transition-colors flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              Edit Profil
            </router-link>
            <button
              @click="handleLogout"
              class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Logout
            </button>
          </div>
        </transition>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { onClickOutside } from '@vueuse/core'

const router = useRouter()
const authStore = useAuthStore()

defineEmits(['toggle-sidebar', 'toggle-mobile-menu'])

const showDropdown = ref(false)
const dropdownRef = ref(null)

// API Base URL for images
const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000'

// Helper to resolve image URLs (handles relative and absolute)
function resolveImageUrl(url: string | undefined) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return apiBaseUrl + url
}

// Computed avatar URL
const avatarUrl = computed(() => resolveImageUrl(authStore.user?.avatar_path))

// Close dropdown when clicking outside
onClickOutside(dropdownRef, () => {
  showDropdown.value = false
})

const roleLabel = computed(() => {
  const roles: Record<string, string> = {
    super_admin: 'Super Admin',
    bendahara: 'Bendahara',
    operator: 'Operator',
    admin: 'Admin'
  }
  return (authStore.user?.role && roles[authStore.user.role]) || 'User'
})

const roleBadgeClass = computed(() => {
  const baseClass = 'px-3 py-1 rounded-full text-xs font-medium'
  const roleClasses: Record<string, string> = {
    super_admin: 'bg-purple-100 text-purple-700',
    bendahara: 'bg-emerald-100 text-emerald-700',
    operator: 'bg-blue-100 text-blue-700',
    admin: 'bg-blue-100 text-blue-700'
  }
  return `${baseClass} ${authStore.user?.role ? roleClasses[authStore.user.role] : 'bg-secondary-100 text-secondary-700'}`
})

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.dropdown-enter-active {
  transition: all 0.2s ease-out;
}

.dropdown-leave-active {
  transition: all 0.15s ease-in;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
