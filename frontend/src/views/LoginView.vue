<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";
import axios from "axios";

const authStore = useAuthStore();
const router = useRouter();

const username = ref("");
const password = ref("");
const error = ref("");
const loading = ref(true);
const showPassword = ref(false);

// Login page settings
const settings = ref({
  login_logo_url: "",
  login_logo_size: "80",
  login_title: "Selamat Datang",
  login_subtitle: "Silakan login untuk mengakses Sistem Manajemen Dokumen Keuangan",
  login_info_title: "Informasi",
  login_info_content: "",
  login_bg_color: "#f3f4f6",
  login_accent_color: "#3b82f6",
  login_font_family: "Inter",
  login_title_size: "24",
  login_subtitle_size: "14",
});

const fetchLoginSettings = async () => {
  try {
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8000/api";
    const response = await axios.get(`${apiUrl}/public/login-settings`);
    const data = response.data.data || {};
    
    settings.value = {
      login_logo_url: data.login_logo_url || "",
      login_logo_size: data.login_logo_size || "80",
      login_title: data.login_title || "Selamat Datang",
      login_subtitle: data.login_subtitle || "Silakan login untuk mengakses Sistem Manajemen Dokumen Keuangan",
      login_info_title: data.login_info_title || "Informasi",
      login_info_content: data.login_info_content || "",
      login_bg_color: data.login_bg_color || "#f3f4f6",
      login_accent_color: data.login_accent_color || "#3b82f6",
      login_font_family: data.login_font_family || "Inter",
      login_title_size: data.login_title_size || "24",
      login_subtitle_size: data.login_subtitle_size || "14",
    };
  } catch {
    console.log("Using default login settings");
  } finally {
    loading.value = false;
  }
};

// Compute full logo URL
const getLogoUrl = () => {
  const url = settings.value.login_logo_url;
  if (!url) return "";
  if (url.startsWith("/")) {
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8000/api";
    // Remove /api from the end if present
    const baseUrl = apiUrl.replace(/\/api$/, "");
    return baseUrl + url;
  }
  return url;
};

const handleLogin = async () => {
  error.value = "";
  try {
    await authStore.login(username.value, password.value);
    router.push("/");
  } catch (e: unknown) {
    const err = e as { response?: { data?: { message?: string } } };
    error.value =
      err.response?.data?.message ||
      "Login gagal. Periksa username dan password.";
  }
};

// Info section types
interface InfoItem {
  text: string;
  isBold: boolean;
  linkText?: string;
  linkUrl?: string;
}

interface InfoSection {
  title: string;
  icon: string;
  color: string;
  items: InfoItem[];
}

const infoSections = ref<InfoSection[]>([]);

// Parse info content from JSON
const parseInfoContent = () => {
  const content = settings.value.login_info_content;
  if (!content) {
    infoSections.value = [];
    return;
  }
  
  try {
    infoSections.value = JSON.parse(content);
  } catch {
    // Legacy format: plain text per line
    infoSections.value = [];
  }
};

// Format item text with bold and links
const formatItemText = (item: InfoItem): string => {
  let text = item.text;
  
  // Apply bold if needed
  if (item.isBold) {
    text = `<b>${text}</b>`;
  }
  
  // Apply link if linkText exists
  if (item.linkText && text.includes(item.linkText)) {
    if (item.linkUrl) {
      text = text.replace(
        item.linkText, 
        `<a href="${item.linkUrl}" class="text-blue-600 hover:underline">${item.linkText}</a>`
      );
    } else {
      text = text.replace(
        item.linkText, 
        `<b class="text-blue-600">${item.linkText}</b>`
      );
    }
  }
  
  return text;
};

onMounted(async () => {
  await fetchLoginSettings();
  parseInfoContent();
});
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center p-4"
    :style="{ backgroundColor: settings.login_bg_color }"
  >
    <!-- Loading state -->
    <div v-if="loading" class="flex items-center justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-blue-500 border-t-transparent"></div>
    </div>

    <!-- Main Card - Contains both Info and Form -->
    <div 
      v-else 
      class="bg-white rounded-2xl shadow-2xl overflow-hidden w-full max-w-5xl flex flex-col lg:flex-row"
      style="box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);"
    >
      <!-- Left Panel - Information -->
      <div class="lg:w-1/2 p-8 lg:p-10 border-b lg:border-b-0 lg:border-r border-gray-100">
        <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center gap-2">
          <span class="text-blue-500">ℹ️</span>
          {{ settings.login_info_title }}
        </h2>
        
        <!-- Info Sections -->
        <div class="space-y-6">
          <!-- Default Info if no sections defined -->
          <div v-if="infoSections.length === 0" class="space-y-4">
            <div class="border-l-4 border-blue-500 pl-4">
              <h3 class="font-semibold text-gray-700 mb-2 flex items-center gap-2">
                <span class="text-green-500">✓</span>
                Informasi Akun
              </h3>
              <ul class="text-sm text-gray-600 space-y-2">
                <li>1. Gunakan username dan password yang telah diberikan oleh administrator.</li>
                <li>2. Jika Anda <b>belum memiliki akun</b>, silakan hubungi administrator sistem.</li>
                <li>3. Pastikan Anda logout setelah selesai menggunakan aplikasi.</li>
              </ul>
            </div>
            
            <div class="border-l-4 border-orange-500 pl-4">
              <h3 class="font-semibold text-gray-700 mb-2 flex items-center gap-2">
                <span class="text-orange-500">📋</span>
                Panduan Login
              </h3>
              <ul class="text-sm text-gray-600 space-y-2">
                <li>1. Masukkan username dan password dengan benar.</li>
                <li>2. Klik tombol <b>Login</b> untuk masuk ke sistem.</li>
                <li>3. Jika lupa password, hubungi <b class="text-blue-600">Administrator</b>.</li>
              </ul>
            </div>
          </div>

          <!-- Custom Sections from settings -->
          <div v-else class="space-y-4">
            <div
              v-for="(section, sectionIndex) in infoSections"
              :key="sectionIndex"
              class="border-l-4 pl-4"
              :style="{ borderColor: section.color }"
            >
              <h3 class="font-semibold text-gray-700 mb-2 flex items-center gap-2">
                <span :style="{ color: section.color }">{{ section.icon }}</span>
                {{ section.title }}
              </h3>
              <ul class="text-sm text-gray-600 space-y-2">
                <li 
                  v-for="(item, itemIndex) in section.items" 
                  :key="itemIndex"
                  v-html="`${itemIndex + 1}. ${formatItemText(item)}`"
                ></li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Panel - Login Form -->
      <div class="lg:w-1/2 p-8 lg:p-10 flex flex-col justify-center">
        <!-- Logo & Header -->
        <div class="text-center mb-8">
          <div v-if="settings.login_logo_url" class="mb-4">
            <img
              :src="getLogoUrl()"
              alt="Logo"
              class="w-auto mx-auto object-contain"
              :style="{ height: settings.login_logo_size + 'px' }"
            />
          </div>
          <div v-else class="mb-4">
            <div 
              class="mx-auto rounded-full flex items-center justify-center text-4xl"
              :style="{ 
                width: settings.login_logo_size + 'px', 
                height: settings.login_logo_size + 'px',
                backgroundColor: settings.login_accent_color + '15', 
                color: settings.login_accent_color 
              }"
            >
              📋
            </div>
          </div>
          
          <h1 
            class="font-bold text-gray-800"
            :style="{ 
              fontFamily: settings.login_font_family, 
              fontSize: settings.login_title_size + 'px' 
            }"
          >
            {{ settings.login_title }}
          </h1>
          <p 
            class="text-gray-500 mt-2"
            :style="{ 
              fontFamily: settings.login_font_family, 
              fontSize: settings.login_subtitle_size + 'px' 
            }"
          >
            {{ settings.login_subtitle }}
          </p>
        </div>

        <!-- Divider -->
        <div class="flex items-center gap-3 mb-6">
          <div class="flex-1 h-px bg-gray-200"></div>
          <span class="text-gray-400 text-xs uppercase tracking-wider font-medium">Login dengan Username</span>
          <div class="flex-1 h-px bg-gray-200"></div>
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-5">
          <div
            v-if="error"
            class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm flex items-center gap-2"
          >
            <span>⚠️</span>
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Username
            </label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">👤</span>
              <input
                v-model="username"
                type="text"
                required
                class="w-full pl-12 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="Username Anda"
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">🔒</span>
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="w-full pl-12 pr-12 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="••••••••"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
              >
                <span v-if="showPassword">👁️</span>
                <span v-else>👁️‍🗨️</span>
              </button>
            </div>
          </div>

          <button
            type="submit"
            :disabled="authStore.loading"
            class="w-full py-3 rounded-lg font-semibold text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:shadow-lg transform hover:-translate-y-0.5"
            :style="{ backgroundColor: settings.login_accent_color }"
          >
            <span v-if="authStore.loading" class="flex items-center justify-center gap-2">
              <span class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></span>
              Loading...
            </span>
            <span v-else>Login</span>
          </button>
        </form>

        <!-- Footer -->
        <div class="mt-8 pt-6 border-t border-gray-100 text-center">
          <p class="text-gray-400 text-xs">
            Sistem Manajemen Dokumen Keuangan © {{ new Date().getFullYear() }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
