<template>
  <div class="min-h-screen flex">
    <!-- Left Side: Dark Background with 3D Animation -->
    <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 relative overflow-hidden">
      <!-- Floating decorative elements with animation -->
      <div class="absolute inset-0">
        <!-- Floating shapes -->
        <div class="float-slow absolute top-[10%] left-[8%] w-4 h-4 border border-white/20 rotate-45"></div>
        <div class="float-slower absolute top-[15%] right-[12%] w-8 h-8 rounded-full border border-white/15"></div>
        <div class="float-slow absolute bottom-[25%] left-[12%] w-6 h-6 rounded-full border border-white/20"></div>
        <div class="float-slower absolute bottom-[20%] right-[15%] w-10 h-10 border border-white/10 rotate-12"></div>
        <div class="float-slow absolute top-[45%] left-[5%] w-3 h-3 bg-white/10 rounded-full"></div>
        <div class="float-slower absolute top-[55%] right-[8%] w-5 h-5 border border-white/15 rotate-45"></div>
      </div>
      
      <!-- Content with 3D Animation -->
      <div class="flex flex-col items-center justify-center w-full px-12 z-10">
        <!-- 3D Rotating Logo Animation -->
        <div class="mb-12 logo-container">
          <div class="logo-3d">
            <!-- Arrow shape 1 - Blue -->
            <div class="arrow arrow-1">
              <div class="arrow-face arrow-face-1"></div>
              <div class="arrow-face arrow-face-2"></div>
            </div>
            <!-- Arrow shape 2 - Cyan -->
            <div class="arrow arrow-2">
              <div class="arrow-face arrow-face-3"></div>
              <div class="arrow-face arrow-face-4"></div>
            </div>
            <!-- Arrow shape 3 - Yellow -->
            <div class="arrow arrow-3">
              <div class="arrow-face arrow-face-5"></div>
            </div>
          </div>
        </div>
        
        <!-- Title -->
        <div class="text-center text-white">
          <h1 class="text-3xl font-bold mb-2">{{ branding.app_name }}</h1>
          <p class="text-xl text-slate-300">{{ branding.app_subtitle }}</p>
        </div>
      </div>
    </div>

    <!-- Right Side: Login Form -->
    <div class="w-full lg:w-1/2 flex items-center justify-center p-8 bg-slate-50">
      <div class="w-full max-w-md">
        <!-- Mobile Logo -->
        <div class="lg:hidden text-center mb-8">
          <div class="inline-flex items-center justify-center w-20 h-20 mb-4">
            <img 
              v-if="logoUrl"
              :src="logoUrl" 
              alt="Logo" 
              class="w-full h-full object-contain"
            />
            <div v-else class="w-full h-full rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center">
              <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <h1 class="text-2xl font-bold text-secondary-900">{{ branding.app_name }}</h1>
          <p class="text-secondary-600">{{ branding.app_subtitle }}</p>
        </div>

        <!-- Title -->
        <div class="mb-8">
          <h2 class="text-2xl font-bold text-secondary-900">Masuk ke Sistem</h2>
          <p class="text-secondary-500 mt-1">Mohon masukkan informasi akun Anda</p>
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- Error message -->
          <div v-if="error" class="p-4 bg-red-50 border border-red-200 rounded-lg">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>

          <!-- Tahun Anggaran -->
          <div>
            <label class="block text-sm font-medium text-secondary-700 mb-1.5">Tahun Anggaran</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-secondary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <select
                v-model="form.tahun_anggaran"
                class="input pl-10 appearance-none cursor-pointer"
                required
              >
                <option value="2025">Tahun Anggaran 2025</option>
                <option value="2026">Tahun Anggaran 2026</option>
              </select>
              <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-secondary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
            <p class="text-xs text-emerald-600 mt-1 flex items-center gap-1">
              <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="5" />
              </svg>
              Tahun Anggaran Aktif
            </p>
          </div>

          <!-- Username -->
          <div>
            <label for="username" class="block text-sm font-medium text-secondary-700 mb-1.5">Username</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-secondary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <input
                id="username"
                v-model="form.username"
                type="text"
                class="input pl-10"
                placeholder="Masukkan username Anda"
                required
                autocomplete="username"
              />
            </div>
          </div>

          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-secondary-700 mb-1.5">Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-secondary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                class="input pl-10 pr-10"
                placeholder="••••••••"
                required
                autocomplete="current-password"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg v-if="!showPassword" class="w-5 h-5 text-secondary-400 hover:text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="w-5 h-5 text-secondary-400 hover:text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Cloudflare Turnstile -->
          <div>
            <div id="turnstile-container" class="flex justify-center"></div>
            <p v-if="turnstileError" class="text-xs text-red-500 mt-1 text-center">{{ turnstileError }}</p>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="authStore.loading || !turnstileToken"
            class="w-full btn-primary py-3 text-base disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg v-if="authStore.loading" class="animate-spin -ml-1 mr-2 h-5 w-5 text-white inline" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ authStore.loading ? 'Memproses...' : 'Masuk ke Sistem' }}
          </button>
        </form>

        <!-- Footer -->
        <p class="text-center text-xs text-secondary-400 mt-8">
          © {{ currentYear }} {{ branding.app_name }}. All rights reserved.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

// Add missing window type definition
declare global {
  interface Window {
    turnstile: any;
  }
}

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: '',
  tahun_anggaran: '2025'
})

const showPassword = ref(false)
const error = ref('')
const turnstileToken = ref('')
const turnstileError = ref('')
const turnstileWidgetId = ref<string | null>(null)

const branding = ref({
  app_name: 'Sistem Pelimpahan',
  app_subtitle: 'Dana UP/GU',
  logo_url: ''
})

// API Base URL for images
const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000'

// Helper to resolve image URLs
function resolveImageUrl(url: string | undefined) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return apiBaseUrl + url
}

// Computed logo URL from branding settings
const logoUrl = computed(() => resolveImageUrl(branding.value.logo_url))

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Cloudflare Turnstile site key (use test key for development)
const TURNSTILE_SITE_KEY = import.meta.env.VITE_TURNSTILE_SITE_KEY || '1x00000000000000000000AA' // Test key

// Initialize Turnstile
function initTurnstile() {
  if (window.turnstile) {
    turnstileWidgetId.value = window.turnstile.render('#turnstile-container', {
      sitekey: TURNSTILE_SITE_KEY,
      callback: (token: string) => {
        turnstileToken.value = token
        turnstileError.value = ''
      },
      'error-callback': () => {
        turnstileError.value = 'Verifikasi gagal. Silakan coba lagi.'
        turnstileToken.value = ''
      },
      'expired-callback': () => {
        turnstileToken.value = ''
      },
      theme: 'light'
    })
  }
}

// Load Turnstile script
function loadTurnstileScript() {
  if (document.getElementById('turnstile-script')) {
    initTurnstile()
    return
  }
  
  const script = document.createElement('script')
  script.id = 'turnstile-script'
  script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js'
  script.async = true
  script.defer = true
  script.onload = () => {
    // Wait a bit for turnstile to be ready
    setTimeout(initTurnstile, 100)
  }
  document.head.appendChild(script)
}

onMounted(async () => {
  // Load branding
  // Load branding
  try {
    const response = await api.get('/public/login-settings')
    if (response.data.data) {
      const settings = response.data.data
      branding.value = {
        app_name: settings.login_title || 'Sistem Pelimpahan',
        app_subtitle: settings.login_subtitle || 'Dana UP/GU',
        logo_url: settings.login_logo_url || ''
      }
    }
  } catch (err) {
    console.log('Using default branding')
  }
  
  // Load Turnstile
  loadTurnstileScript()
})

onUnmounted(() => {
  // Cleanup Turnstile widget
  if (turnstileWidgetId.value && window.turnstile) {
    window.turnstile.remove(turnstileWidgetId.value)
  }
})

async function handleLogin() {
  error.value = ''
  
  if (!turnstileToken.value) {
    error.value = 'Silakan selesaikan verifikasi terlebih dahulu'
    return
  }
  
  const result = await authStore.login(
    form.username, 
    form.password, 
    form.tahun_anggaran,
    turnstileToken.value
  )
  
  if (result.success) {
    router.push('/')
  } else {
    error.value = result.message || 'Login Failed'
    // Reset Turnstile on error
    if (window.turnstile && turnstileWidgetId.value) {
      window.turnstile.reset(turnstileWidgetId.value)
      turnstileToken.value = ''
    }
  }
}
</script>

<style scoped>
/* Logo container with perspective */
.logo-container {
  perspective: 1000px;
  width: 250px;
  height: 250px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 3D rotating logo */
.logo-3d {
  width: 200px;
  height: 200px;
  position: relative;
  transform-style: preserve-3d;
  animation: rotate3d 8s linear infinite;
}

/* Arrow base styles */
.arrow {
  position: absolute;
  transform-style: preserve-3d;
}

/* Arrow 1 - Blue pointing right */
.arrow-1 {
  width: 100px;
  height: 80px;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -80%) rotateY(0deg);
}

.arrow-face-1 {
  position: absolute;
  width: 0;
  height: 0;
  border-left: 50px solid transparent;
  border-right: 50px solid transparent;
  border-bottom: 80px solid #2563eb;
  filter: brightness(1.1);
}

.arrow-face-2 {
  position: absolute;
  width: 0;
  height: 0;
  border-left: 50px solid transparent;
  border-right: 50px solid transparent;
  border-bottom: 80px solid #1d4ed8;
  transform: rotateY(180deg);
}

/* Arrow 2 - Cyan pointing bottom-left */
.arrow-2 {
  width: 100px;
  height: 80px;
  top: 50%;
  left: 50%;
  transform: translate(-100%, 10%) rotateZ(120deg);
}

.arrow-face-3 {
  position: absolute;
  width: 0;
  height: 0;
  border-left: 50px solid transparent;
  border-right: 50px solid transparent;
  border-bottom: 80px solid #06b6d4;
  filter: brightness(1.1);
}

.arrow-face-4 {
  position: absolute;
  width: 0;
  height: 0;
  border-left: 50px solid transparent;
  border-right: 50px solid transparent;
  border-bottom: 80px solid #0891b2;
  transform: rotateY(180deg);
}

/* Arrow 3 - Yellow pointing bottom-right */
.arrow-3 {
  width: 100px;
  height: 80px;
  top: 50%;
  left: 50%;
  transform: translate(0%, 10%) rotateZ(-120deg);
}

.arrow-face-5 {
  position: absolute;
  width: 0;
  height: 0;
  border-left: 50px solid transparent;
  border-right: 50px solid transparent;
  border-bottom: 80px solid #fbbf24;
  filter: brightness(1.1);
}

/* Floating decorative elements */
.float-slow {
  animation: floatSlow 6s ease-in-out infinite;
}

.float-slower {
  animation: floatSlower 8s ease-in-out infinite;
}

/* Keyframes */
@keyframes rotate3d {
  0% {
    transform: rotateY(0deg) rotateX(10deg);
  }
  100% {
    transform: rotateY(360deg) rotateX(10deg);
  }
}

@keyframes floatSlow {
  0%, 100% {
    transform: translateY(0) rotate(45deg);
    opacity: 0.2;
  }
  50% {
    transform: translateY(-15px) rotate(45deg);
    opacity: 0.4;
  }
}

@keyframes floatSlower {
  0%, 100% {
    transform: translateY(0) rotate(12deg);
    opacity: 0.15;
  }
  50% {
    transform: translateY(-20px) rotate(12deg);
    opacity: 0.3;
  }
}
</style>
