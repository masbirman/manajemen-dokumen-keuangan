<template>
  <aside
    :class="[
      'bg-white border-r border-secondary-200 transition-all duration-300 flex flex-col h-screen min-h-screen',
      collapsed ? 'w-20' : 'w-64'
    ]"
  >
    <!-- Logo -->
    <div class="h-16 flex items-center justify-center border-b border-secondary-200 px-4">
      <div v-if="!collapsed" class="flex items-center gap-3">
        <div v-if="logoUrl" class="w-10 h-10 rounded-xl overflow-hidden">
          <img :src="logoUrl" alt="Logo" class="w-full h-full object-cover" />
        </div>
        <div v-else class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div>
          <h1 class="font-bold text-secondary-900 text-sm">{{ branding.app_name }}</h1>
          <p class="text-xs text-secondary-500">{{ branding.app_subtitle }}</p>
        </div>
      </div>
      <div v-else>
        <div v-if="logoUrl" class="w-10 h-10 rounded-xl overflow-hidden">
          <img :src="logoUrl" alt="Logo" class="w-full h-full object-cover" />
        </div>
        <div v-else class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 py-4 px-3 space-y-1 overflow-y-auto">
      <router-link
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        v-slot="{ isActive }"
      >
        <div
          v-if="!item.roles || authStore.hasAnyRole(item.roles)"
          :class="[
            'sidebar-link',
            { 'active': isActive }
          ]"
          :title="collapsed ? item.label : ''"
        >
          <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
          <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
        </div>
      </router-link>
    </nav>

    <!-- User Info -->
    <div class="border-t border-secondary-200 p-4">
      <div :class="['flex items-center gap-3', collapsed ? 'justify-center' : '']">
        <div v-if="avatarUrl" class="w-10 h-10 rounded-full overflow-hidden">
          <img :src="avatarUrl" alt="Avatar" class="w-full h-full object-cover" />
        </div>
        <div v-else class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white font-semibold">
          {{ authStore.user?.name?.charAt(0).toUpperCase() || 'U' }}
        </div>
        <div v-if="!collapsed" class="flex-1 min-w-0">
          <p class="text-sm font-medium text-secondary-900 truncate">{{ authStore.user?.name }}</p>
          <p class="text-xs text-secondary-500 truncate">{{ roleLabel }}</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

const authStore = useAuthStore()

const branding = ref({
  app_name: 'Pelimpahan',
  app_subtitle: 'Dana UP/GU',
  logo_url: ''
})

// API Base URL for images
const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000'

// Helper to resolve image URLs (handles relative and absolute)
function resolveImageUrl(url: string | undefined) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return apiBaseUrl + url
}

// Computed logo URL
const logoUrl = computed(() => resolveImageUrl(branding.value.logo_url))

// Computed avatar URL
const avatarUrl = computed(() => resolveImageUrl(authStore.user?.avatar_path))

onMounted(async () => {
  try {
    const response = await api.get('/settings/branding')
    if (response.data.success && response.data.data) {
      branding.value = {
        app_name: response.data.data.app_name || 'Pelimpahan',
        app_subtitle: response.data.data.app_subtitle || 'Dana UP/GU',
        logo_url: response.data.data.logo_url || ''
      }
    }
  } catch (error) {
    console.log('Using default branding')
  }
})

const roleLabel = computed(() => {
  const roles: Record<string, string> = {
    super_admin: 'Super Admin',
    bendahara: 'Bendahara',
    operator: 'Operator'
  }
  return (authStore.user?.role && roles[authStore.user.role]) || 'User'
})

// Icon components
const DashboardIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' })
    ])
  }
}

const UnitIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4' })
    ])
  }
}

const JenisIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z' })
    ])
  }
}

const PelimpahanIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z' })
    ])
  }
}

const InputIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 6v6m0 0v6m0-6h6m-6 0H6' })
    ])
  }
}

const UserIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' })
    ])
  }
}

const SaldoIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z' })
    ])
  }
}

const LaporanIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' })
    ])
  }
}

const BackupIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4' })
    ])
  }
}

const SettingsIcon = {
  render() {
    return h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' }),
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z' })
    ])
  }
}

interface MenuItem {
    path: string;
    label: string;
    icon: any;
    roles?: string[];
}

const menuItems: MenuItem[] = [
  { path: '/', label: 'Dashboard', icon: DashboardIcon },
  { path: '/dokumen/input', label: 'Input Dokumen', icon: InputIcon, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/dokumen', label: 'List Dokumen', icon: PelimpahanIcon },
  { path: '/unit-kerja', label: 'Unit Kerja', icon: UnitIcon, roles: ['super_admin', 'admin'] },
  { path: '/jenis-dokumen', label: 'Jenis Dokumen', icon: JenisIcon, roles: ['super_admin', 'admin'] },
  { path: '/sumber-dana', label: 'Sumber Dana', icon: SaldoIcon, roles: ['super_admin', 'admin'] },
  { path: '/pptk', label: 'PPTK', icon: UserIcon, roles: ['super_admin', 'admin'] },
  { path: '/users', label: 'Manajemen User', icon: UserIcon, roles: ['super_admin'] },
  { path: '/settings', label: 'Pengaturan', icon: SettingsIcon, roles: ['super_admin'] },
  { path: '/login-settings', label: 'Tampilan Login', icon: SettingsIcon, roles: ['super_admin'] },
  { path: '/login-content', label: 'Konten Login', icon: LaporanIcon, roles: ['super_admin'] }
]
</script>
