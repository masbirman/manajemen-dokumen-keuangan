<script setup lang="ts">
import { ref, onMounted, computed, reactive } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { InputField } from "@/components/ui";

interface Settings {
  app_name: string;
  app_subtitle: string;
  logo_url: string;
  // Countdown settings
  countdown_active: boolean;
  countdown_title: string;
  countdown_description: string;
  countdown_target_date: string;
}

const toast = useToast();
const loading = ref(false);
const saving = ref(false);
const uploading = ref(false);
const activeTab = ref("branding"); // 'branding' | 'countdown'
const fileInput = ref<HTMLInputElement | null>(null);

const settings = ref<Settings>({
  app_name: "",
  app_subtitle: "",
  logo_url: "",
  countdown_active: false,
  countdown_title: "",
  countdown_description: "",
  countdown_target_date: "",
});

const countdownDate = ref("");
const countdownTime = ref("");

// Helper for image URL
const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000';
const getImageUrl = (path: string) => {
  if (!path) return '';
  if (path.startsWith('http')) return path;
  return apiBaseUrl + path;
};

const fetchSettings = async () => {
  loading.value = true;
  try {
    const response = await api.get("/settings");
    const data = response.data.map || {}; // Using the 'map' from the API response
    
    // Parse existing target date if available
    let datePart = "";
    let timePart = "";
    if (data.countdown_target_date) {
      const dateObj = new Date(data.countdown_target_date);
      if (!isNaN(dateObj.getTime())) {
        // Format to YYYY-MM-DD
        const year = dateObj.getFullYear();
        const month = String(dateObj.getMonth() + 1).padStart(2, '0');
        const day = String(dateObj.getDate()).padStart(2, '0');
        datePart = `${year}-${month}-${day}`;
        
        // Format to HH:mm
        const hours = String(dateObj.getHours()).padStart(2, '0');
        const minutes = String(dateObj.getMinutes()).padStart(2, '0');
        timePart = `${hours}:${minutes}`;
      }
    }

    // Map login_* keys to branding settings
    settings.value = {
      app_name: data.login_title || "Sistem Pelimpahan",
      app_subtitle: data.login_subtitle || "Dana UP/GU",
      logo_url: data.login_logo_url || "",
      countdown_active: data.countdown_active === "true",
      countdown_title: data.countdown_title || "",
      countdown_description: data.countdown_description || "",
      countdown_target_date: data.countdown_target_date || "",
    };
    
    // Set separate refs
    countdownDate.value = datePart;
    countdownTime.value = timePart;
    
  } catch {
    toast.error("Gagal memuat pengaturan");
  } finally {
    loading.value = false;
  }
};

const handleLogoUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files?.length) return;

  const file = target.files[0];
  if (file.size > 2 * 1024 * 1024) {
    toast.error('Ukuran file maksimal 2MB');
    return;
  }

  const formData = new FormData();
  formData.append('logo', file);

  uploading.value = true;
  try {
    const response = await api.post('/settings/upload-logo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    settings.value.logo_url = response.data.logo_url;
    toast.success('Logo berhasil diupload');
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal upload logo');
  } finally {
    uploading.value = false;
    if (fileInput.value) fileInput.value.value = '';
  }
};

const triggerFileInput = () => {
  fileInput.value?.click();
};

const saveSettings = async () => {
  saving.value = true;
  try {
    // Combine date and time matching Project-04 format
    let combinedDate = "";
    if (countdownDate.value && countdownTime.value) {
      combinedDate = `${countdownDate.value}T${countdownTime.value}:00`;
    } else if (countdownDate.value) {
      combinedDate = `${countdownDate.value}T00:00:00`;
    }

    // Explicitly construct payload with map keys expected by backend (login_* for branding)
    const payload = {
      settings: {
        login_title: settings.value.app_name,
        login_subtitle: settings.value.app_subtitle,
        login_logo_url: settings.value.logo_url,
        countdown_active: String(settings.value.countdown_active),
        countdown_title: settings.value.countdown_title,
        countdown_description: settings.value.countdown_description,
        countdown_target_date: combinedDate
      }
    };
    
    await api.put("/settings", payload);
    toast.success("Pengaturan berhasil disimpan");
    
    // Update local state to match
    settings.value.countdown_target_date = combinedDate;
    
  } catch {
    toast.error("Gagal menyimpan pengaturan");
  } finally {
    saving.value = false;
  }
};

onMounted(fetchSettings);
</script>

<template>
  <div class="max-w-4xl">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Pengaturan</h1>
        <p class="text-gray-500">Kelola pengaturan aplikasi</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex border-b border-gray-200 mb-6">
      <button
        @click="activeTab = 'branding'"
        :class="[
          'px-4 py-2 text-sm font-medium border-b-2 transition-colors',
          activeTab === 'branding'
            ? 'border-blue-600 text-blue-600'
            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
        ]"
      >
        Branding Aplikasi
      </button>
      <button
        @click="activeTab = 'countdown'"
        :class="[
          'px-4 py-2 text-sm font-medium border-b-2 transition-colors',
          activeTab === 'countdown'
            ? 'border-blue-600 text-blue-600'
            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
        ]"
      >
        Countdown Dashboard
      </button>
    </div>

    <div
      v-if="loading"
      class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center"
    >
      <p class="text-gray-500">Loading...</p>
    </div>

    <div
      v-else
      class="bg-white rounded-xl shadow-sm p-6 border border-gray-100"
    >
      <form @submit.prevent="saveSettings">
        <!-- Branding Tab -->
        <div v-if="activeTab === 'branding'" class="space-y-6">
            <div class="bg-blue-50 p-4 rounded-lg border border-blue-100 mb-6">
                <h3 class="text-sm font-medium text-blue-800 mb-1">Branding Aplikasi</h3>
                <p class="text-xs text-blue-600">
                  Ganti logo dan nama aplikasi yang tampil di sidebar.
                </p>
            </div>

            <InputField v-model="settings.app_name" label="Nama Aplikasi" />
            <InputField
                v-model="settings.app_subtitle"
                label="Subtitle Aplikasi"
            />

            <!-- Logo Upload -->
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Logo Aplikasi</label>
                <div class="flex items-start gap-6">
                    <div class="w-20 h-20 rounded-xl bg-gray-100 border border-gray-200 overflow-hidden flex items-center justify-center">
                        <img 
                            v-if="settings.logo_url" 
                            :src="getImageUrl(settings.logo_url)" 
                            alt="Logo" 
                            class="w-full h-full object-cover"
                        />
                        <svg v-else class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                    </div>
                    <div>
                        <button 
                            type="button"
                            @click="triggerFileInput"
                            :disabled="uploading"
                            class="px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 flex items-center gap-2"
                        >
                             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                            </svg>
                            {{ uploading ? 'Mengupload...' : 'Upload Logo' }}
                        </button>
                        <p class="text-xs text-gray-500 mt-2">PNG, JPG, WEBP. Max 2MB. Ukuran ideal: 100x100px</p>
                        <input 
                            type="file" 
                            ref="fileInput" 
                            class="hidden" 
                            accept="image/*"
                            @change="handleLogoUpload"
                        />
                    </div>
                </div>
            </div>

            <!-- Preview Sidebar -->
            <div class="mt-8 pt-6 border-t border-gray-100">
                <p class="text-sm font-medium text-gray-500 mb-4">Preview Sidebar:</p>
                <div class="w-64 bg-white border border-gray-200 rounded-lg shadow-sm">
                    <div class="h-16 flex items-center px-4 border-b border-gray-200">
                        <div class="flex items-center gap-3">
                             <div class="w-10 h-10 rounded-xl overflow-hidden bg-gradient-to-br from-blue-500 to-blue-700 flex-shrink-0">
                                <img 
                                    v-if="settings.logo_url" 
                                    :src="getImageUrl(settings.logo_url)" 
                                    alt="Logo" 
                                    class="w-full h-full object-cover"
                                />
                                <div v-else class="w-full h-full flex items-center justify-center">
                                     <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                                    </svg>
                                </div>
                            </div>
                            <div>
                                <h1 class="font-bold text-gray-900 text-sm leading-tight">{{ settings.app_name || 'Nama Aplikasi' }}</h1>
                                <p class="text-xs text-gray-500 leading-tight mt-0.5">{{ settings.app_subtitle || 'Subtitle' }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Countdown Tab -->
        <div v-if="activeTab === 'countdown'" class="space-y-6">
          <div class="bg-blue-50 p-4 rounded-lg border border-blue-100">
            <h3 class="text-sm font-medium text-blue-800 mb-1">Countdown Dashboard</h3>
            <p class="text-xs text-blue-600">
              Tampilkan countdown di halaman dashboard untuk event atau deadline tertentu.
            </p>
          </div>

          <!-- Active Switch -->
          <div class="flex items-center justify-between p-4 border rounded-lg">
            <div>
              <p class="font-medium text-gray-900">Aktifkan Countdown</p>
              <p class="text-sm text-gray-500">Tampilkan countdown di dashboard</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input 
                type="checkbox" 
                v-model="settings.countdown_active" 
                class="sr-only peer"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>

          <InputField
            v-model="settings.countdown_title"
            label="Judul Countdown"
            placeholder="Contoh: Jadwal Penatausahaan Tahun 2025"
          />
          
          <InputField
            v-model="settings.countdown_description"
            label="Deskripsi (opsional)"
            placeholder="Contoh: Penatausahaan Keuangan akan berakhir tanggal 31 Desember 2025"
          />

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Target</label>
              <input
                v-model="countdownDate"
                type="date"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Jam Target</label>
              <input
                v-model="countdownTime"
                type="time"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
          </div>

          <!-- Preview -->
          <div v-if="settings.countdown_active" class="mt-8">
             <p class="text-sm font-medium text-gray-500 mb-2">Preview:</p>
             <div class="p-6 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl border border-blue-100">
               <h4 class="text-lg font-bold text-blue-900">{{ settings.countdown_title || 'Judul Countdown' }}</h4>
               <p class="text-sm text-blue-600 mt-1 mb-4">{{ settings.countdown_description || 'Deskripsi countdown' }}</p>
               
               <div class="flex gap-4">
                 <div class="text-center">
                   <span class="block text-2xl font-bold text-blue-600">09</span>
                   <span class="text-xs text-blue-400">Hari</span>
                 </div>
                 <div class="text-center">
                   <span class="block text-2xl font-bold text-blue-600">07</span>
                   <span class="text-xs text-blue-400">Jam</span>
                 </div>
                 <div class="text-center">
                   <span class="block text-2xl font-bold text-blue-600">03</span>
                   <span class="text-xs text-blue-400">Menit</span>
                 </div>
                 <div class="text-center">
                   <span class="block text-2xl font-bold text-blue-600">43</span>
                   <span class="text-xs text-blue-400">Detik</span>
                 </div>
               </div>
             </div>
          </div>
        </div>

        <div class="flex justify-end pt-6 border-t border-gray-100 mt-6 gap-3">
          <button 
             type="button" 
             class="px-6 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium"
          >
             Batal
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 flex items-center gap-2 font-medium shadow-sm transition-all"
          >
            <svg v-if="saving" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ saving ? "Menyimpan..." : "Simpan Branding" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
