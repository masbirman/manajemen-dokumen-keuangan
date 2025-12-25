<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";

interface LoginContent {
  id: string;
  title: string;
  description: string;
  image_url: string;
  image_width: number;
  title_size: number;
  desc_size: number;
  start_date: string;
  end_date: string;
  is_active: boolean;
}

interface ContentForm {
  title: string;
  description: string;
  start_date: string;
  end_date: string;
  image_width: number;
  title_size: number;
  desc_size: number;
}

const toast = useToast();
const loading = ref(false);
const saving = ref(false);
const loginContents = ref<LoginContent[]>([]);
const showModal = ref(false);
const editingContent = ref<LoginContent | null>(null);
const previewImageUrl = ref<string | null>(null);
const selectedFile = ref<File | null>(null);

const contentForm = ref<ContentForm>({
  title: '',
  description: '',
  start_date: '',
  end_date: '',
  image_width: 400,
  title_size: 28,
  desc_size: 16
});

const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000';

const getImageUrl = (url: string) => {
  if (!url) return '';
  if (url.startsWith('http')) return url;
  return apiBaseUrl + url;
};

const fetchContents = async () => {
  loading.value = true;
  try {
    const response = await api.get('/login-content');
    if (response.data.success) {
      loginContents.value = response.data.data || [];
    }
  } catch {
    toast.error("Gagal memuat konten login");
  } finally {
    loading.value = false;
  }
};

const openModal = (content: LoginContent | null = null) => {
  previewImageUrl.value = null;
  selectedFile.value = null;
  
  if (content) {
    editingContent.value = content;
    contentForm.value = {
      title: content.title,
      description: content.description || '',
      start_date: content.start_date?.split('T')[0] || '',
      end_date: content.end_date?.split('T')[0] || '',
      image_width: content.image_width || 400,
      title_size: content.title_size || 28,
      desc_size: content.desc_size || 16
    };
  } else {
    editingContent.value = null;
    contentForm.value = {
      title: '',
      description: '',
      start_date: '',
      end_date: '',
      image_width: 400,
      title_size: 28,
      desc_size: 16
    };
  }
  showModal.value = true;
};

const saveContent = async () => {
  if (!contentForm.value.title || !contentForm.value.start_date || !contentForm.value.end_date) {
    toast.error('Judul dan tanggal wajib diisi');
    return;
  }
  
  saving.value = true;
  try {
    let contentId: string | null = null;
    
    if (editingContent.value) {
      await api.put(`/login-content/${editingContent.value.id}`, contentForm.value);
      contentId = editingContent.value.id;
      toast.success('Konten berhasil diupdate');
    } else {
      const response = await api.post('/login-content', contentForm.value);
      contentId = response.data.data?.id;
      toast.success('Konten berhasil ditambahkan');
    }
    
    // Upload image if selected
    if (selectedFile.value && contentId) {
      const formData = new FormData();
      formData.append('image', selectedFile.value);
      await api.post(`/login-content/${contentId}/image`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
    }
    
    showModal.value = false;
    previewImageUrl.value = null;
    selectedFile.value = null;
    await fetchContents();
  } catch {
    toast.error('Gagal menyimpan');
  } finally {
    saving.value = false;
  }
};

const deleteContent = async (content: LoginContent) => {
  if (!confirm(`Hapus konten "${content.title}"?`)) return;
  
  try {
    await api.delete(`/login-content/${content.id}`);
    toast.success('Konten berhasil dihapus');
    await fetchContents();
  } catch {
    toast.error('Gagal menghapus');
  }
};

const handleImageSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;
  
  if (file.size > 2 * 1024 * 1024) {
    toast.error('Ukuran file maksimal 2MB');
    return;
  }
  
  selectedFile.value = file;
  previewImageUrl.value = URL.createObjectURL(file);
};

const isContentActive = (content: LoginContent) => {
  const today = new Date().toISOString().split('T')[0];
  const start = content.start_date?.split('T')[0];
  const end = content.end_date?.split('T')[0];
  return start <= today && today <= end && content.is_active;
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
};

onMounted(fetchContents);
</script>

<template>
  <div class="max-w-4xl">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Konten Login Terjadwal</h1>
        <p class="text-gray-500 text-sm mt-1">Kelola konten halaman login berdasarkan jadwal</p>
      </div>
      <button
        @click="openModal(null)"
        class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 flex items-center gap-2"
      >
        <span>+</span>
        Tambah Konten
      </button>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center">
      <div class="animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent mx-auto"></div>
      <p class="text-gray-500 mt-4">Loading...</p>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-gray-100">
      <div v-if="loginContents.length === 0" class="p-8 text-center text-gray-500">
        <p>Belum ada konten. Klik "Tambah Konten" untuk menambah.</p>
      </div>
      <div v-else class="divide-y divide-gray-200">
        <div v-for="content in loginContents" :key="content.id" class="p-4 flex items-center justify-between hover:bg-gray-50">
          <div class="flex items-center gap-4">
            <div v-if="content.image_url" class="w-16 h-16 rounded-lg overflow-hidden bg-gray-100">
              <img :src="getImageUrl(content.image_url)" class="w-full h-full object-cover" />
            </div>
            <div v-else class="w-16 h-16 rounded-lg bg-gray-100 flex items-center justify-center text-2xl">🖼️</div>
            <div>
              <p class="font-medium text-gray-900">{{ content.title }}</p>
              <p class="text-sm text-gray-500">{{ formatDate(content.start_date) }} - {{ formatDate(content.end_date) }}</p>
              <span :class="isContentActive(content) ? 'text-green-600 text-xs' : 'text-gray-400 text-xs'">
                {{ isContentActive(content) ? '✅ Aktif Sekarang' : '⏳ Terjadwal' }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button @click="openModal(content)" class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg">
              ✏️
            </button>
            <button @click="deleteContent(content)" class="p-2 text-red-600 hover:bg-red-50 rounded-lg">
              🗑️
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="showModal = false"></div>
      <div class="relative bg-white rounded-xl shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">{{ editingContent ? 'Edit Konten' : 'Tambah Konten' }}</h3>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Judul <span class="text-red-500">*</span></label>
            <input v-model="contentForm.title" type="text" class="w-full px-4 py-2 border border-gray-300 rounded-lg" placeholder="Selamat Natal 2025! 🎄" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
            <textarea v-model="contentForm.description" class="w-full px-4 py-2 border border-gray-300 rounded-lg" rows="3" placeholder="Pesan singkat..."></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Mulai <span class="text-red-500">*</span></label>
              <input v-model="contentForm.start_date" type="date" class="w-full px-4 py-2 border border-gray-300 rounded-lg" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Selesai <span class="text-red-500">*</span></label>
              <input v-model="contentForm.end_date" type="date" class="w-full px-4 py-2 border border-gray-300 rounded-lg" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lebar Gambar (px)</label>
              <input v-model.number="contentForm.image_width" type="number" class="w-full px-4 py-2 border border-gray-300 rounded-lg" min="200" max="800" />
              <p class="text-xs text-gray-400 mt-1">200-800px. Landscape: 600-800px</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Font Judul (px)</label>
              <input v-model.number="contentForm.title_size" type="number" class="w-full px-4 py-2 border border-gray-300 rounded-lg" min="12" max="48" />
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Gambar Ilustrasi</label>
            <div v-if="editingContent?.image_url || previewImageUrl" class="mb-2">
              <img :src="previewImageUrl || getImageUrl(editingContent!.image_url)" class="max-w-full max-h-48 object-contain rounded-lg bg-gray-100" />
            </div>
            <input type="file" @change="handleImageSelect" accept="image/*" class="w-full px-4 py-2 border border-gray-300 rounded-lg" />
            <p class="text-xs text-gray-500 mt-1">PNG, JPG, WebP. Max 2MB. Landscape direkomendasikan.</p>
          </div>
        </div>
        <div class="p-6 border-t border-gray-200 flex justify-end gap-3">
          <button @click="showModal = false" class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Batal</button>
          <button @click="saveContent" :disabled="saving" class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50">
            {{ saving ? 'Menyimpan...' : 'Simpan' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
