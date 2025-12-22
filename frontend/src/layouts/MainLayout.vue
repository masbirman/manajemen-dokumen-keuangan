<template>
  <div class="flex h-screen bg-secondary-100 overflow-hidden">
    <!-- Desktop Sidebar (Hidden on Mobile) -->
    <aside class="hidden lg:block h-full flex-shrink-0 z-30">
      <Sidebar :collapsed="sidebarCollapsed" />
    </aside>

    <!-- Mobile Sidebar (Off-canvas) -->
    <div 
      v-if="mobileMenuOpen" 
      class="fixed inset-0 z-50 lg:hidden"
    >
      <!-- Backdrop -->
      <div 
        class="absolute inset-0 bg-black/50 transition-opacity"
        @click="mobileMenuOpen = false"
      ></div>
      
      <!-- Drawer -->
      <aside class="absolute inset-y-0 left-0 w-64 max-w-[80%] bg-white shadow-xl transform transition-transform duration-300">
        <Sidebar :collapsed="false" @close="mobileMenuOpen = false" />
      </aside>
    </div>
    
    <!-- Main Content Wrapper -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Navbar -->
      <Navbar 
        @toggle-sidebar="toggleSidebar" 
        @toggle-mobile-menu="mobileMenuOpen = !mobileMenuOpen"
      />
      
      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto p-4 md:p-6 pb-20 md:pb-6 relative">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- Notifications -->
    <div class="fixed top-4 right-4 z-50 space-y-2 pointer-events-none">
      <TransitionGroup name="notification">
        <div
          v-for="notification in notificationStore.notifications"
          :key="notification.id"
          :class="[
            'pointer-events-auto p-4 rounded-lg shadow-lg max-w-[90vw] md:max-w-sm animate-fadeIn',
            {
              'bg-emerald-500 text-white': notification.type === 'success',
              'bg-red-500 text-white': notification.type === 'error',
              'bg-amber-500 text-white': notification.type === 'warning',
              'bg-primary-500 text-white': notification.type === 'info'
            }
          ]"
        >
          <div class="flex items-center justify-between gap-3 md:gap-4">
            <p class="text-xs md:text-sm font-medium break-words">{{ notification.message }}</p>
            <button 
              @click="notificationStore.remove(notification.id)"
              class="hover:opacity-75 flex-shrink-0"
            >
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import Navbar from '@/components/layout/Navbar.vue'
import { useNotificationStore } from '@/stores/notification'

const sidebarCollapsed = ref(false)
const mobileMenuOpen = ref(false)
const notificationStore = useNotificationStore()

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.notification-enter-active {
  transition: all 0.3s ease-out;
}

.notification-leave-active {
  transition: all 0.2s ease-in;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
