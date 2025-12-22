<template>
  <div class="space-y-6 animate-fadeIn">
    <!-- Welcome & Stats Row -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Welcome Widget (Left - 3 cols) -->
      <div class="lg:col-span-3 card overflow-hidden bg-white">
        <div class="card-body p-6">
          <div class="grid grid-cols-1 md:grid-cols-[1fr,200px] lg:grid-cols-[1fr,220px] xl:grid-cols-[1fr,260px] gap-6">
            <!-- Left Column: Greeting + Countdown -->
            <div class="space-y-4 min-w-0">
              <!-- Greeting -->
              <div>
                <h2 class="text-lg md:text-2xl font-bold text-primary-800 whitespace-nowrap" style="min-height: 2rem;">
                  <span class="typing-text">{{ displayedGreeting }}</span>
                  <span v-if="!typingComplete" class="typing-cursor">|</span>
                </h2>
                <p class="text-sm text-secondary-500 mt-1">{{ currentDate }}</p>
                <p class="text-secondary-600 mt-3">
                  Selamat bekerja di Sistem Pelimpahan Dokumen Keuangan.
                </p>
                <p class="text-secondary-500 text-sm mt-1">
                  Kelola dokumen dan pelimpahan dengan mudah dan efisien.
                </p>
              </div>
              
              <!-- Countdown Box (Left) -->
              <div v-if="countdown.active" class="p-4 bg-gradient-to-r from-primary-50 to-blue-50 rounded-xl border border-primary-100">
                <p class="text-sm font-medium text-primary-700">{{ countdown.title }}</p>
                <p class="text-xs text-secondary-500 mt-1 mb-3">{{ countdown.description }}</p>
                <div class="grid grid-cols-4 gap-2">
                  <div class="bg-white px-2 py-2 rounded-lg shadow-sm text-center">
                    <p class="text-lg md:text-xl font-bold text-primary-700">{{ countdownDisplay.days }}</p>
                    <p class="text-[10px] md:text-xs text-secondary-500">Hari</p>
                  </div>
                  <div class="bg-white px-2 py-2 rounded-lg shadow-sm text-center">
                    <p class="text-lg md:text-xl font-bold text-primary-700">{{ countdownDisplay.hours }}</p>
                    <p class="text-[10px] md:text-xs text-secondary-500">Jam</p>
                  </div>
                  <div class="bg-white px-2 py-2 rounded-lg shadow-sm text-center">
                    <p class="text-lg md:text-xl font-bold text-primary-700">{{ countdownDisplay.minutes }}</p>
                    <p class="text-[10px] md:text-xs text-secondary-500">Menit</p>
                  </div>
                  <div class="bg-white px-2 py-2 rounded-lg shadow-sm text-center">
                    <p class="text-lg md:text-xl font-bold text-primary-700">{{ countdownDisplay.seconds }}</p>
                    <p class="text-[10px] md:text-xs text-secondary-500">Detik</p>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Right Column: Illustration (Responsive Width) -->
            <div class="hidden md:flex items-center justify-center">
              <img 
                :src="apiBaseUrl + '/uploads/ilustrasi-dashboard.webp'" 
                alt="Welcome Illustration" 
                class="h-48 w-full object-contain"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Stats Widget (Right - 1 col) -->
      <div class="lg:col-span-1 card flex flex-col">
        <div class="card-header py-2">
          <h3 class="font-semibold text-secondary-900 text-sm">Statistik</h3>
        </div>
        <div class="card-body p-2 flex-1 flex flex-col">
          <!-- Stats Grid -->
          <div class="grid grid-cols-2 gap-2 flex-1">
            <!-- Total Pelimpahan -->
            <div class="p-3 bg-blue-50 rounded-xl text-center flex flex-col justify-center">
              <div class="w-8 h-8 mx-auto mb-1 rounded-lg bg-blue-500 flex items-center justify-center">
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <p class="text-base md:text-xl font-bold text-blue-700">{{ stats.total_pelimpahan }}</p>
              <p class="text-[10px] md:text-xs text-blue-600">Total Dokumen</p>
            </div>

            <!-- Unit Kerja -->
            <div class="p-3 bg-emerald-50 rounded-xl text-center flex flex-col justify-center">
              <div class="w-8 h-8 mx-auto mb-1 rounded-lg bg-emerald-500 flex items-center justify-center">
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
              </div>
              <p class="text-base md:text-xl font-bold text-emerald-700">{{ stats.total_unit }}</p>
              <p class="text-[10px] md:text-xs text-emerald-600">Unit Kerja</p>
            </div>

            <!-- Jenis Pelimpahan -->
            <div class="p-3 bg-amber-50 rounded-xl text-center flex flex-col justify-center">
              <div class="w-8 h-8 mx-auto mb-1 rounded-lg bg-amber-500 flex items-center justify-center">
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
              </div>
              <p class="text-base md:text-xl font-bold text-amber-700">{{ stats.total_jenis }}</p>
              <p class="text-[10px] md:text-xs text-amber-600">Jenis Dokumen</p>
            </div>

            <!-- Total User -->
            <div class="p-3 bg-purple-50 rounded-xl text-center flex flex-col justify-center">
              <div class="w-8 h-8 mx-auto mb-1 rounded-lg bg-purple-500 flex items-center justify-center">
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              </div>
              <p class="text-base md:text-xl font-bold text-purple-700">{{ stats.total_user }}</p>
              <p class="text-[10px] md:text-xs text-purple-600">Total User</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions & Recent -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Quick Actions -->
      <div class="card h-fit">
        <div class="card-header py-3 px-4">
          <h3 class="font-semibold text-secondary-900 text-sm md:text-base">Aksi Cepat</h3>
        </div>
        <div class="card-body p-4">
          <div class="space-y-3">
            <router-link to="/dokumen/input" class="flex items-center gap-3 p-3 bg-primary-50 rounded-xl hover:bg-primary-100 transition-colors group">
              <div class="w-10 h-10 rounded-lg bg-primary-500 flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-semibold text-primary-700">Input Dokumen</p>
                <p class="text-[10px] text-primary-500">Buat pelimpahan baru</p>
              </div>
            </router-link>

            <router-link to="/dokumen" class="flex items-center gap-3 p-3 bg-emerald-50 rounded-xl hover:bg-emerald-100 transition-colors group">
              <div class="w-10 h-10 rounded-lg bg-emerald-500 flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-semibold text-emerald-700">List Dokumen</p>
                <p class="text-[10px] text-emerald-500">Lihat semua data</p>
              </div>
            </router-link>
          </div>
        </div>
      </div>

      <!-- Recent Pelimpahan -->
      <div class="lg:col-span-2 card">
        <div class="card-header flex items-center justify-between">
          <h3 class="font-semibold text-secondary-900">Dokumen Terbaru</h3>
          <router-link to="/dokumen" class="text-sm text-primary-600 hover:text-primary-700">
            Lihat Semua →
          </router-link>
        </div>
        <div class="card-body p-0">
          <div v-if="stats.recent_pelimpahan?.length" class="divide-y divide-secondary-200">
            <div
              v-for="item in stats.recent_pelimpahan"
              :key="item.id"
              class="p-4 hover:bg-secondary-50 transition-colors"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <span class="badge-info">{{ item.jenis_dokumen?.kode_jenis || 'DOC' }}</span>
                  <div>
                    <p class="font-medium text-secondary-900">Dokumen: {{ item.nomor_dokumen || item.nomor_pelimpahan }}</p>
                    <p class="text-sm text-secondary-500">{{ item.uraian || item.uraian_pelimpahan || '-' }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="font-semibold text-primary-600">{{ formatCurrency(item.jumlah || item.total_jumlah) }}</p>
                  <p class="text-xs text-secondary-400">{{ formatDate(item.tanggal || item.tanggal_pelimpahan) }}</p>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="p-8 text-center text-secondary-500">
            <svg class="w-12 h-12 mx-auto mb-4 text-secondary-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
            <p>Belum ada data dokumen</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface DashboardStats {
  total_pelimpahan: number;
  total_unit: number;
  total_jenis: number;
  total_user: number;
  recent_pelimpahan: any[];
}

const stats = ref<DashboardStats>({
  total_pelimpahan: 0,
  total_unit: 0,
  total_jenis: 0,
  total_user: 0,
  recent_pelimpahan: []
})



// API Base URL for images
const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000'

// Typing effect
const displayedGreeting = ref('')
const typingComplete = ref(false)
let typingInterval: any = null

// Countdown
const countdown = ref({
  active: false,
  title: '',
  description: '',
  target_date: null
})
const countdownDisplay = ref({ days: '00', hours: '00', minutes: '00', seconds: '00' })
let countdownInterval: any = null

const userName = computed(() => {
  return authStore.user?.name?.split(' ')[0] || 'User'
})

const fullGreeting = computed(() => {
  const hour = new Date().getHours()
  let greet = 'Selamat Malam'
  if (hour < 11) greet = 'Selamat Pagi'
  else if (hour < 15) greet = 'Selamat Siang'
  else if (hour < 18) greet = 'Selamat Sore'
  return `${greet}, ${userName.value}! 🎉`
})

const currentDate = computed(() => {
  const options: Intl.DateTimeFormatOptions = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }
  return new Date().toLocaleDateString('id-ID', options)
})

// Looping typing effect with multiple texts
function startTypingLoop() {
  const texts = [
    fullGreeting.value,
    'Sistem Pelimpahan Dokumen Keuangan'
  ]
  let textIndex = 0
  let i = 0
  let isDeleting = false
  
  const type = () => {
    const currentText = texts[textIndex]
    
    if (!isDeleting) {
      if (i < currentText.length) {
        displayedGreeting.value = currentText.substring(0, i + 1)
        i++
        typingComplete.value = false
        typingInterval = setTimeout(type, 60)
      } else {
        typingComplete.value = true
        typingInterval = setTimeout(() => {
          isDeleting = true
          type()
        }, 3000) // Wait 3 seconds before deleting
      }
    } else {
      if (i > 0) {
        displayedGreeting.value = currentText.substring(0, i - 1)
        i--
        typingComplete.value = false
        typingInterval = setTimeout(type, 30)
      } else {
        isDeleting = false
        textIndex = (textIndex + 1) % texts.length // Move to next text
        typingInterval = setTimeout(type, 500) // Wait before typing again
      }
    }
  }
  
  setTimeout(type, 500)
}

// Countdown timer calculation
function updateCountdown() {
  if (!countdown.value.target_date) return
  
  const now = new Date().getTime()
  const target = new Date(countdown.value.target_date).getTime()
  const diff = target - now
  
  if (diff <= 0) {
    countdownDisplay.value = { days: '00', hours: '00', minutes: '00', seconds: '00' }
    countdown.value.active = false
    return
  }
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)
  
  countdownDisplay.value = {
    days: String(days).padStart(2, '0'),
    hours: String(hours).padStart(2, '0'),
    minutes: String(minutes).padStart(2, '0'),
    seconds: String(seconds).padStart(2, '0')
  }
}

onMounted(async () => {
  // Load stats
  try {
    const response = await api.get('/dashboard/stats')
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load dashboard stats:', error)
  }



  // Load countdown settings
  try {
    const response = await api.get('/settings/countdown')
    if (response.data.success && response.data.data) {
      countdown.value = response.data.data
      if (countdown.value.active) {
        updateCountdown()
        countdownInterval = setInterval(updateCountdown, 1000)
      }
    }
  } catch (error) {
    console.log('No countdown settings found')
  }

  // Start typing effect
  startTypingLoop()
})

onUnmounted(() => {
  if (typingInterval) clearTimeout(typingInterval)
  if (countdownInterval) clearInterval(countdownInterval)
})

function formatCurrency(value: number | undefined) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}

function formatDate(date: string | undefined) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}
</script>

<style scoped>
.typing-cursor {
  animation: blink 1s step-end infinite;
  color: #3b82f6;
  font-weight: normal;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
</style>
