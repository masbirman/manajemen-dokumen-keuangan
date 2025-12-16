<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { InputField } from "@/components/ui";

interface LoginSettings {
  login_logo_url: string;
  login_logo_size: string;
  login_title: string;
  login_subtitle: string;
  login_info_title: string;
  login_info_content: string;
  login_bg_color: string;
  login_accent_color: string;
  login_font_family: string;
  login_title_size: string;
  login_subtitle_size: string;
}

const toast = useToast();
const loading = ref(false);
const saving = ref(false);
const settings = ref<LoginSettings>({
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

// Font options
const fontOptions = [
  { value: "Inter", label: "Inter" },
  { value: "Roboto", label: "Roboto" },
  { value: "Poppins", label: "Poppins" },
  { value: "Open Sans", label: "Open Sans" },
  { value: "Lato", label: "Lato" },
  { value: "Montserrat", label: "Montserrat" },
  { value: "Nunito", label: "Nunito" },
  { value: "Source Sans Pro", label: "Source Sans Pro" },
];

const fetchSettings = async () => {
  loading.value = true;
  try {
    const response = await api.get("/settings");
    const data = response.data.map || {};
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
    toast.error("Gagal memuat pengaturan halaman login");
  } finally {
    loading.value = false;
  }
};

const saveSettings = async () => {
  saving.value = true;
  try {
    await api.put("/settings", {
      settings: {
        login_logo_url: settings.value.login_logo_url,
        login_logo_size: settings.value.login_logo_size,
        login_title: settings.value.login_title,
        login_subtitle: settings.value.login_subtitle,
        login_info_title: settings.value.login_info_title,
        login_info_content: settings.value.login_info_content,
        login_bg_color: settings.value.login_bg_color,
        login_accent_color: settings.value.login_accent_color,
        login_font_family: settings.value.login_font_family,
        login_title_size: settings.value.login_title_size,
        login_subtitle_size: settings.value.login_subtitle_size,
      },
    });
    toast.success("Pengaturan halaman login berhasil disimpan");
  } catch {
    toast.error("Gagal menyimpan pengaturan");
  } finally {
    saving.value = false;
  }
};

const previewLogin = () => {
  window.open("/login", "_blank");
};

// Logo upload
const logoInput = ref<HTMLInputElement | null>(null);
const uploadingLogo = ref(false);

const triggerLogoUpload = () => {
  logoInput.value?.click();
};

const handleLogoUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/svg+xml', 'image/webp'];
  if (!allowedTypes.includes(file.type)) {
    toast.error("Format file tidak didukung. Gunakan JPG, PNG, GIF, SVG, atau WebP.");
    return;
  }

  // Validate file size (max 2MB)
  if (file.size > 2 * 1024 * 1024) {
    toast.error("Ukuran file maksimal 2MB");
    return;
  }

  uploadingLogo.value = true;
  try {
    const formData = new FormData();
    formData.append("logo", file);

    const response = await api.post("/settings/upload-logo", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    settings.value.login_logo_url = response.data.logo_url;
    toast.success("Logo berhasil diupload");
  } catch {
    toast.error("Gagal mengupload logo");
  } finally {
    uploadingLogo.value = false;
    // Reset input
    if (target) target.value = "";
  }
};

const removeLogo = async () => {
  settings.value.login_logo_url = "";
  try {
    await api.put("/settings", {
      settings: {
        login_logo_url: "",
      },
    });
    toast.success("Logo berhasil dihapus");
  } catch {
    toast.error("Gagal menghapus logo");
  }
};

// ============================================
// Structured Info Panel Content Editor
// ============================================

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

// Color options for section border
const sectionColors = [
  { value: "#3b82f6", label: "Biru", bg: "bg-blue-500" },
  { value: "#f97316", label: "Orange", bg: "bg-orange-500" },
  { value: "#22c55e", label: "Hijau", bg: "bg-green-500" },
  { value: "#a855f7", label: "Ungu", bg: "bg-purple-500" },
  { value: "#ef4444", label: "Merah", bg: "bg-red-500" },
  { value: "#6b7280", label: "Abu-abu", bg: "bg-gray-500" },
];

// Icon options for section
const sectionIcons = [
  { value: "✓", label: "Centang" },
  { value: "📋", label: "Clipboard" },
  { value: "ℹ️", label: "Info" },
  { value: "⚠️", label: "Warning" },
  { value: "💡", label: "Lampu" },
  { value: "📌", label: "Pin" },
  { value: "🔑", label: "Kunci" },
  { value: "👤", label: "User" },
];

// Parse JSON content to sections
const parseInfoContent = () => {
  const content = settings.value.login_info_content;
  if (!content) {
    // Default sections
    infoSections.value = [
      {
        title: "Informasi Akun",
        icon: "✓",
        color: "#3b82f6",
        items: [
          { text: "Gunakan username dan password yang telah diberikan oleh administrator.", isBold: false },
          { text: "Jika Anda belum memiliki akun, silakan hubungi administrator sistem.", isBold: false, linkText: "administrator", linkUrl: "" },
          { text: "Pastikan Anda logout setelah selesai menggunakan aplikasi.", isBold: false },
        ]
      },
      {
        title: "Panduan Login",
        icon: "📋",
        color: "#f97316",
        items: [
          { text: "Masukkan username dan password dengan benar.", isBold: false },
          { text: "Klik tombol Login untuk masuk ke sistem.", isBold: false, linkText: "Login", linkUrl: "" },
          { text: "Jika lupa password, hubungi Administrator.", isBold: false, linkText: "Administrator", linkUrl: "" },
        ]
      }
    ];
    return;
  }
  
  try {
    infoSections.value = JSON.parse(content);
  } catch {
    // If not JSON, create default
    infoSections.value = [];
  }
};

// Convert sections to JSON string for saving
const sectionsToJson = () => {
  return JSON.stringify(infoSections.value);
};

// Add new section
const addSection = () => {
  infoSections.value.push({
    title: "Seksi Baru",
    icon: "ℹ️",
    color: "#3b82f6",
    items: []
  });
};

// Remove section
const removeSection = (index: number) => {
  infoSections.value.splice(index, 1);
};

// Add item to section
const addItem = (sectionIndex: number) => {
  infoSections.value[sectionIndex].items.push({
    text: "",
    isBold: false,
  });
};

// Remove item from section
const removeItem = (sectionIndex: number, itemIndex: number) => {
  infoSections.value[sectionIndex].items.splice(itemIndex, 1);
};

// Move section up
const moveSectionUp = (index: number) => {
  if (index > 0) {
    const temp = infoSections.value[index];
    infoSections.value[index] = infoSections.value[index - 1];
    infoSections.value[index - 1] = temp;
  }
};

// Move section down
const moveSectionDown = (index: number) => {
  if (index < infoSections.value.length - 1) {
    const temp = infoSections.value[index];
    infoSections.value[index] = infoSections.value[index + 1];
    infoSections.value[index + 1] = temp;
  }
};

// Update saveSettings to include sections
const saveSettingsWithSections = async () => {
  saving.value = true;
  try {
    // Convert sections to JSON
    settings.value.login_info_content = sectionsToJson();
    
    await api.put("/settings", {
      settings: {
        login_logo_url: settings.value.login_logo_url,
        login_logo_size: settings.value.login_logo_size,
        login_title: settings.value.login_title,
        login_subtitle: settings.value.login_subtitle,
        login_info_title: settings.value.login_info_title,
        login_info_content: settings.value.login_info_content,
        login_bg_color: settings.value.login_bg_color,
        login_accent_color: settings.value.login_accent_color,
        login_font_family: settings.value.login_font_family,
        login_title_size: settings.value.login_title_size,
        login_subtitle_size: settings.value.login_subtitle_size,
      },
    });
    toast.success("Pengaturan halaman login berhasil disimpan");
  } catch {
    toast.error("Gagal menyimpan pengaturan");
  } finally {
    saving.value = false;
  }
};

onMounted(async () => {
  await fetchSettings();
  parseInfoContent();
});
</script>

<template>
  <div class="max-w-4xl">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Pengaturan Halaman Login</h1>
        <p class="text-gray-500 text-sm mt-1">Konfigurasi tampilan halaman login untuk pengguna</p>
      </div>
      <button
        @click="previewLogin"
        class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 flex items-center gap-2"
      >
        <span>👁️</span>
        Preview
      </button>
    </div>

    <div
      v-if="loading"
      class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center"
    >
      <div class="animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent mx-auto"></div>
      <p class="text-gray-500 mt-4">Loading...</p>
    </div>

    <div v-else class="space-y-6">
      <!-- Preview Card -->
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
          <span>🎨</span>
          Preview Halaman Login
        </h2>
        <div
          class="rounded-lg overflow-hidden h-48 flex items-center justify-center p-4"
          :style="{ backgroundColor: settings.login_bg_color }"
        >
          <!-- Single Card Preview -->
          <div class="bg-white rounded-lg shadow-xl flex overflow-hidden w-full max-w-md h-40">
            <!-- Left - Info Panel -->
            <div class="w-1/2 p-3 border-r border-gray-100">
              <div class="text-xs font-semibold text-gray-700 mb-2 flex items-center gap-1">
                <span class="text-blue-500 text-xs">ℹ️</span>
                {{ settings.login_info_title }}
              </div>
              <div class="space-y-1">
                <div class="h-2 bg-gray-200 rounded w-full"></div>
                <div class="h-2 bg-gray-200 rounded w-3/4"></div>
                <div class="h-2 bg-gray-200 rounded w-5/6"></div>
                <div class="h-2 bg-gray-200 rounded w-2/3"></div>
              </div>
            </div>
            <!-- Right - Form Panel -->
            <div class="w-1/2 p-3 flex flex-col items-center justify-center">
              <div v-if="settings.login_logo_url" class="h-6 mb-1 flex items-center justify-center">
                <img :src="settings.login_logo_url" class="h-full object-contain" />
              </div>
              <div v-else class="h-6 w-6 rounded-full mb-1 flex items-center justify-center text-xs"
                   :style="{ backgroundColor: settings.login_accent_color + '20' }">
                📋
              </div>
              <div 
                class="font-semibold text-gray-700 truncate mb-2"
                :style="{ fontFamily: settings.login_font_family, fontSize: '10px' }"
              >
                {{ settings.login_title }}
              </div>
              <div class="w-full space-y-1">
                <div class="h-4 bg-gray-100 rounded w-full"></div>
                <div class="h-4 bg-gray-100 rounded w-full"></div>
                <div class="h-4 rounded w-full" :style="{ backgroundColor: settings.login_accent_color }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Settings Form -->
      <form @submit.prevent="saveSettings" class="space-y-6">
        <!-- Appearance Settings -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
            <span>🎨</span>
            Tampilan
          </h2>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Warna Background</label>
              <div class="flex items-center gap-3">
                <input
                  v-model="settings.login_bg_color"
                  type="color"
                  class="w-12 h-12 rounded-lg border border-gray-300 cursor-pointer"
                />
                <input
                  v-model="settings.login_bg_color"
                  type="text"
                  class="flex-1 px-4 py-2 border border-gray-300 rounded-lg"
                  placeholder="#f3f4f6"
                />
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Warna Aksen (Tombol)</label>
              <div class="flex items-center gap-3">
                <input
                  v-model="settings.login_accent_color"
                  type="color"
                  class="w-12 h-12 rounded-lg border border-gray-300 cursor-pointer"
                />
                <input
                  v-model="settings.login_accent_color"
                  type="text"
                  class="flex-1 px-4 py-2 border border-gray-300 rounded-lg"
                  placeholder="#3b82f6"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Font Settings -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
            <span>🔤</span>
            Pengaturan Font
          </h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Jenis Font</label>
              <select
                v-model="settings.login_font_family"
                class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option v-for="font in fontOptions" :key="font.value" :value="font.value">
                  {{ font.label }}
                </option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Ukuran Judul (px)</label>
              <input
                v-model="settings.login_title_size"
                type="number"
                min="16"
                max="48"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="24"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Ukuran Subtitle (px)</label>
              <input
                v-model="settings.login_subtitle_size"
                type="number"
                min="10"
                max="24"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="14"
              />
            </div>
          </div>

          <!-- Font Preview -->
          <div class="mt-4 p-4 bg-gray-50 rounded-lg">
            <p class="text-xs text-gray-500 mb-2">Preview Font:</p>
            <p 
              class="text-gray-800 font-bold"
              :style="{ fontFamily: settings.login_font_family, fontSize: settings.login_title_size + 'px' }"
            >
              {{ settings.login_title }}
            </p>
            <p 
              class="text-gray-500 mt-1"
              :style="{ fontFamily: settings.login_font_family, fontSize: settings.login_subtitle_size + 'px' }"
            >
              {{ settings.login_subtitle }}
            </p>
          </div>
        </div>

        <!-- Logo & Header Settings -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
            <span>🖼️</span>
            Logo & Header
          </h2>
          
          <div class="space-y-4">
            <!-- Logo Upload Section -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Logo</label>
              
              <!-- Hidden file input -->
              <input
                ref="logoInput"
                type="file"
                accept="image/jpeg,image/png,image/gif,image/svg+xml,image/webp"
                class="hidden"
                @change="handleLogoUpload"
              />
              
              <!-- Logo Preview or Upload Button -->
              <div class="flex items-start gap-4">
                <div 
                  v-if="settings.login_logo_url"
                  class="relative group"
                >
                  <div class="p-4 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
                    <img 
                      :src="settings.login_logo_url.startsWith('/') ? 'http://localhost:8000' + settings.login_logo_url : settings.login_logo_url" 
                      alt="Logo Preview"
                      class="object-contain"
                      :style="{ height: settings.login_logo_size + 'px', maxWidth: '200px' }"
                    />
                  </div>
                  <div class="mt-2 flex gap-2">
                    <button
                      type="button"
                      @click="triggerLogoUpload"
                      :disabled="uploadingLogo"
                      class="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition-colors"
                    >
                      Ganti Logo
                    </button>
                    <button
                      type="button"
                      @click="removeLogo"
                      class="px-3 py-1.5 text-sm bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition-colors"
                    >
                      Hapus
                    </button>
                  </div>
                </div>
                
                <!-- Upload Button (when no logo) -->
                <div v-else>
                  <button
                    type="button"
                    @click="triggerLogoUpload"
                    :disabled="uploadingLogo"
                    class="flex flex-col items-center justify-center w-48 h-32 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all cursor-pointer"
                  >
                    <span v-if="uploadingLogo" class="animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></span>
                    <template v-else>
                      <span class="text-3xl text-gray-400 mb-2">📷</span>
                      <span class="text-sm text-gray-500">Klik untuk upload logo</span>
                      <span class="text-xs text-gray-400 mt-1">JPG, PNG, GIF, SVG, WebP (maks 2MB)</span>
                    </template>
                  </button>
                </div>
                
                <!-- Logo Size -->
                <div class="flex-shrink-0">
                  <label class="block text-sm font-medium text-gray-700 mb-2">Ukuran Logo (px)</label>
                  <input
                    v-model="settings.login_logo_size"
                    type="number"
                    min="40"
                    max="200"
                    class="w-24 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="80"
                  />
                </div>
              </div>
            </div>
            
            <InputField
              v-model="settings.login_title"
              label="Judul Halaman Login"
              placeholder="Selamat Datang"
            />
            
            <InputField
              v-model="settings.login_subtitle"
              label="Subtitle"
              placeholder="Silakan login untuk mengakses sistem"
            />
          </div>
        </div>

        <!-- Information Panel Settings -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-700 flex items-center gap-2">
              <span>ℹ️</span>
              Panel Informasi (Kiri)
            </h2>
            <button
              type="button"
              @click="addSection"
              class="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition-colors flex items-center gap-1"
            >
              <span>+</span>
              Tambah Seksi
            </button>
          </div>
          
          <div class="space-y-4">
            <InputField
              v-model="settings.login_info_title"
              label="Judul Panel Informasi"
              placeholder="Informasi"
            />
            
            <!-- Sections Editor -->
            <div class="space-y-4">
              <p class="text-sm font-medium text-gray-700">Konten Informasi</p>
              
              <!-- Empty state -->
              <div v-if="infoSections.length === 0" class="p-8 border-2 border-dashed border-gray-200 rounded-lg text-center">
                <span class="text-3xl text-gray-300">📋</span>
                <p class="text-gray-500 mt-2">Belum ada seksi. Klik "Tambah Seksi" untuk menambahkan.</p>
              </div>
              
              <!-- Section Cards -->
              <div
                v-for="(section, sectionIndex) in infoSections"
                :key="sectionIndex"
                class="border rounded-lg overflow-hidden"
              >
                <!-- Section Header -->
                <div class="bg-gray-50 px-4 py-3 flex items-center gap-3 border-b">
                  <!-- Move buttons -->
                  <div class="flex flex-col gap-0.5">
                    <button
                      type="button"
                      @click="moveSectionUp(sectionIndex)"
                      :disabled="sectionIndex === 0"
                      class="text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed"
                      title="Pindah ke atas"
                    >▲</button>
                    <button
                      type="button"
                      @click="moveSectionDown(sectionIndex)"
                      :disabled="sectionIndex === infoSections.length - 1"
                      class="text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed"
                      title="Pindah ke bawah"
                    >▼</button>
                  </div>
                  
                  <!-- Color border preview -->
                  <div 
                    class="w-1 h-10 rounded-full flex-shrink-0" 
                    :style="{ backgroundColor: section.color }"
                  ></div>
                  
                  <!-- Icon selector -->
                  <select
                    v-model="section.icon"
                    class="w-14 px-2 py-1 border border-gray-200 rounded text-center"
                  >
                    <option v-for="icon in sectionIcons" :key="icon.value" :value="icon.value">
                      {{ icon.value }}
                    </option>
                  </select>
                  
                  <!-- Title input -->
                  <input
                    v-model="section.title"
                    type="text"
                    class="flex-1 px-3 py-1.5 border border-gray-200 rounded focus:ring-1 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Judul Seksi"
                  />
                  
                  <!-- Color selector -->
                  <select
                    v-model="section.color"
                    class="px-3 py-1.5 border border-gray-200 rounded"
                  >
                    <option v-for="color in sectionColors" :key="color.value" :value="color.value">
                      {{ color.label }}
                    </option>
                  </select>
                  
                  <!-- Delete section -->
                  <button
                    type="button"
                    @click="removeSection(sectionIndex)"
                    class="text-red-500 hover:text-red-700 p-1"
                    title="Hapus Seksi"
                  >🗑️</button>
                </div>
                
                <!-- Section Items -->
                <div class="p-4 space-y-2">
                  <div
                    v-for="(item, itemIndex) in section.items"
                    :key="itemIndex"
                    class="flex items-start gap-2"
                  >
                    <span class="flex-shrink-0 w-6 h-6 rounded-full bg-gray-100 text-gray-500 flex items-center justify-center text-xs font-medium mt-1">
                      {{ itemIndex + 1 }}
                    </span>
                    
                    <div class="flex-1 space-y-2">
                      <input
                        v-model="item.text"
                        type="text"
                        class="w-full px-3 py-2 border border-gray-200 rounded focus:ring-1 focus:ring-blue-500 focus:border-transparent text-sm"
                        placeholder="Teks item..."
                      />
                      
                      <!-- Optional: Bold text and link settings -->
                      <div class="flex items-center gap-4 text-xs">
                        <label class="flex items-center gap-1 cursor-pointer">
                          <input
                            type="checkbox"
                            v-model="item.isBold"
                            class="rounded text-blue-500"
                          />
                          <span class="text-gray-500">Bold</span>
                        </label>
                        
                        <div class="flex items-center gap-2">
                          <input
                            v-model="item.linkText"
                            type="text"
                            class="w-24 px-2 py-1 border border-gray-200 rounded text-xs"
                            placeholder="Teks link"
                          />
                          <input
                            v-model="item.linkUrl"
                            type="text"
                            class="w-32 px-2 py-1 border border-gray-200 rounded text-xs"
                            placeholder="URL (opsional)"
                          />
                        </div>
                      </div>
                    </div>
                    
                    <button
                      type="button"
                      @click="removeItem(sectionIndex, itemIndex)"
                      class="text-red-400 hover:text-red-600 p-1 flex-shrink-0"
                      title="Hapus Item"
                    >✕</button>
                  </div>
                  
                  <!-- Add item button -->
                  <button
                    type="button"
                    @click="addItem(sectionIndex)"
                    class="w-full py-2 border-2 border-dashed border-gray-200 rounded text-gray-400 hover:border-blue-300 hover:text-blue-500 transition-colors text-sm"
                  >
                    + Tambah Item
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Save Button -->
        <div class="flex justify-end gap-3">
          <button
            type="button"
            @click="previewLogin"
            class="px-6 py-2.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
          >
            Preview
          </button>
          <button
            type="button"
            @click="saveSettingsWithSections"
            :disabled="saving"
            class="px-6 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors flex items-center gap-2"
          >
            <span v-if="saving" class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></span>
            {{ saving ? "Menyimpan..." : "Simpan Pengaturan" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
