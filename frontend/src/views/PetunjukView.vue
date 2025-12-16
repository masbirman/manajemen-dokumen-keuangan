<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { DataTable, Modal, InputField, Dropdown } from "@/components/ui";

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

interface Petunjuk {
  id: string;
  judul: string;
  konten: string;
  halaman: string;
  urutan: number;
  is_active: boolean;
  image_url?: string;
  image_size?: number;
}

const toast = useToast();
const loading = ref(false);
const data = ref<Petunjuk[]>([]);
const showModal = ref(false);
const showDeleteModal = ref(false);
const editingId = ref<string | null>(null);
const selectedIds = ref<string[]>([]);
const showBulkDeleteModal = ref(false);

// Pagination
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(10);
const search = ref("");

// Form data
const form = ref({
  judul: "",
  konten: "",
  halaman: "input_dokumen",
  urutan: 0,
  is_active: true,
  image_url: "",
  image_size: 200,
});

// Structured content
const infoSections = ref<InfoSection[]>([]);

// Image upload
const imageInput = ref<HTMLInputElement | null>(null);
const uploadingImage = ref(false);

// Section colors and icons
const sectionColors = [
  { value: "#3b82f6", label: "Biru" },
  { value: "#f97316", label: "Orange" },
  { value: "#22c55e", label: "Hijau" },
  { value: "#a855f7", label: "Ungu" },
  { value: "#ef4444", label: "Merah" },
  { value: "#6b7280", label: "Abu-abu" },
];

const sectionIcons = [
  { value: "✓", label: "Centang" },
  { value: "📋", label: "Clipboard" },
  { value: "ℹ️", label: "Info" },
  { value: "⚠️", label: "Warning" },
  { value: "💡", label: "Lampu" },
  { value: "📌", label: "Pin" },
  { value: "🔑", label: "Kunci" },
  { value: "👤", label: "User" },
  { value: "📷", label: "Kamera" },
  { value: "📁", label: "Folder" },
];

const columns = [
  { key: "judul", label: "Judul" },
  { key: "halaman", label: "Halaman", width: "150px" },
  { key: "urutan", label: "Urutan", width: "80px" },
  { key: "is_active", label: "Status", width: "100px" },
  { key: "actions", label: "Aksi", width: "120px" },
];

const halamanOptions = [
  { value: "input_dokumen", label: "Input Dokumen" },
  { value: "list_dokumen", label: "List Dokumen" },
  { value: "dashboard", label: "Dashboard" },
];

const halamanLabel = (value: string) => {
  const opt = halamanOptions.find((o) => o.value === value);
  return opt?.label || value;
};

// Parse content to sections
const parseContent = (content: string) => {
  if (!content) {
    infoSections.value = [];
    return;
  }
  try {
    const parsed = JSON.parse(content);
    infoSections.value = parsed.sections || [];
    form.value.image_url = parsed.image_url || "";
    form.value.image_size = parsed.image_size || 200;
  } catch {
    infoSections.value = [];
  }
};

// Convert sections to JSON
const sectionsToJson = () => {
  return JSON.stringify({
    sections: infoSections.value,
    image_url: form.value.image_url,
    image_size: form.value.image_size,
  });
};

// Section/Item management
const addSection = () => {
  infoSections.value.push({
    title: "Seksi Baru",
    icon: "ℹ️",
    color: "#3b82f6",
    items: [],
  });
};

const removeSection = (index: number) => {
  infoSections.value.splice(index, 1);
};

const addItem = (sectionIndex: number) => {
  infoSections.value[sectionIndex].items.push({
    text: "",
    isBold: false,
  });
};

const removeItem = (sectionIndex: number, itemIndex: number) => {
  infoSections.value[sectionIndex].items.splice(itemIndex, 1);
};

const moveSectionUp = (index: number) => {
  if (index > 0) {
    const temp = infoSections.value[index];
    infoSections.value[index] = infoSections.value[index - 1];
    infoSections.value[index - 1] = temp;
  }
};

const moveSectionDown = (index: number) => {
  if (index < infoSections.value.length - 1) {
    const temp = infoSections.value[index];
    infoSections.value[index] = infoSections.value[index + 1];
    infoSections.value[index + 1] = temp;
  }
};

// Image upload
const triggerImageUpload = () => {
  imageInput.value?.click();
};

const handleImageUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  const allowedTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];
  if (!allowedTypes.includes(file.type)) {
    toast.error("Format file tidak didukung. Gunakan JPG, PNG, GIF, atau WebP.");
    return;
  }

  if (file.size > 5 * 1024 * 1024) {
    toast.error("Ukuran file maksimal 5MB");
    return;
  }

  uploadingImage.value = true;
  try {
    const formData = new FormData();
    formData.append("file", file);

    const response = await api.post("/petunjuk/upload-image", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });

    form.value.image_url = response.data.url;
    toast.success("Gambar berhasil diupload");
  } catch {
    toast.error("Gagal mengupload gambar");
  } finally {
    uploadingImage.value = false;
    if (target) target.value = "";
  }
};

const removeImage = () => {
  form.value.image_url = "";
};

// Data fetching
const fetchData = async (page = 1) => {
  loading.value = true;
  try {
    const params: Record<string, string | number> = {
      page,
      page_size: perPage.value,
    };
    if (search.value) params.search = search.value;

    const response = await api.get("/petunjuk", { params });
    data.value = response.data.data || [];
    currentPage.value = response.data.page || 1;
    totalPages.value = response.data.total_pages || 1;
    totalItems.value = response.data.total || 0;
  } catch {
    toast.error("Gagal memuat data");
  } finally {
    loading.value = false;
  }
};

const onPageChange = (page: number) => fetchData(page);
const onPerPageChange = (newPerPage: number) => {
  perPage.value = newPerPage;
  fetchData(1);
};
const onSearch = (s: string) => {
  search.value = s;
  fetchData(1);
};

const openCreate = () => {
  editingId.value = null;
  form.value = {
    judul: "",
    konten: "",
    halaman: "input_dokumen",
    urutan: 0,
    is_active: true,
    image_url: "",
    image_size: 200,
  };
  infoSections.value = [];
  showModal.value = true;
};

const openEdit = (item: Petunjuk) => {
  editingId.value = item.id;
  form.value = {
    judul: item.judul,
    konten: item.konten,
    halaman: item.halaman,
    urutan: item.urutan,
    is_active: item.is_active,
    image_url: "",
    image_size: 200,
  };
  parseContent(item.konten);
  showModal.value = true;
};

const save = async () => {
  try {
    const payload = {
      judul: form.value.judul,
      konten: sectionsToJson(),
      halaman: form.value.halaman,
      urutan: form.value.urutan,
      is_active: form.value.is_active,
    };

    if (editingId.value) {
      await api.put(`/petunjuk/${editingId.value}`, payload);
      toast.success("Petunjuk berhasil diupdate");
    } else {
      await api.post("/petunjuk", payload);
      toast.success("Petunjuk berhasil ditambahkan");
    }
    showModal.value = false;
    fetchData(currentPage.value);
  } catch {
    toast.error("Gagal menyimpan data");
  }
};

const confirmDelete = (item: Petunjuk) => {
  editingId.value = item.id;
  showDeleteModal.value = true;
};

const deletePetunjuk = async () => {
  try {
    await api.delete(`/petunjuk/${editingId.value}`);
    toast.success("Petunjuk berhasil dihapus");
    showDeleteModal.value = false;
    fetchData(currentPage.value);
  } catch {
    toast.error("Gagal menghapus data");
  }
};

const onSelectionChange = (ids: (string | number)[]) => {
  selectedIds.value = ids as string[];
};

const confirmBulkDelete = async () => {
  try {
    await Promise.all(
      selectedIds.value.map((id) => api.delete(`/petunjuk/${id}`))
    );
    toast.success(`${selectedIds.value.length} petunjuk berhasil dihapus`);
    selectedIds.value = [];
    showBulkDeleteModal.value = false;
    fetchData();
  } catch {
    toast.error("Gagal menghapus data");
  }
};

// Get full image URL
const getImageUrl = (url: string) => {
  if (!url) return "";
  if (url.startsWith("/")) {
    return `http://localhost:8000${url}`;
  }
  return url;
};

onMounted(fetchData);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Kelola Petunjuk</h1>
      <button
        @click="openCreate"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        + Tambah Petunjuk
      </button>
    </div>

    <DataTable
      :columns="columns"
      :data="data"
      :loading="loading"
      :current-page="currentPage"
      :total-pages="totalPages"
      :total-items="totalItems"
      :per-page="perPage"
      searchable
      selectable
      :selected-ids="selectedIds"
      @selection-change="onSelectionChange"
      @page-change="onPageChange"
      @per-page-change="onPerPageChange"
      @search="onSearch"
    >
      <template #bulk-actions>
        <button
          @click="showBulkDeleteModal = true"
          class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700"
        >
          🗑️ Hapus Terpilih
        </button>
      </template>
      <template #halaman="{ value }">
        <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
          {{ halamanLabel(value as string) }}
        </span>
      </template>
      <template #is_active="{ value }">
        <span
          :class="
            value ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
          "
          class="px-2 py-1 text-xs rounded"
        >
          {{ value ? "Aktif" : "Nonaktif" }}
        </span>
      </template>
      <template #actions="{ row }">
        <div class="flex gap-1">
          <button
            @click.stop="openEdit(row as Petunjuk)"
            class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded hover:bg-blue-200"
          >
            ✏️
          </button>
          <button
            @click.stop="confirmDelete(row as Petunjuk)"
            class="px-2 py-1 text-xs bg-red-100 text-red-600 rounded hover:bg-red-200"
          >
            🗑️
          </button>
        </div>
      </template>
    </DataTable>

    <!-- Create/Edit Modal -->
    <Modal
      :show="showModal"
      :title="editingId ? 'Edit Petunjuk' : 'Tambah Petunjuk'"
      size="3xl"
      @close="showModal = false"
    >
      <div class="space-y-6 max-h-[70vh] overflow-y-auto pr-2">
        <!-- Basic Info -->
        <div class="bg-gray-50 rounded-lg p-4 border">
          <h3 class="font-semibold text-gray-700 mb-3 flex items-center gap-2">
            <span>📝</span> Informasi Dasar
          </h3>
          <div class="grid grid-cols-2 gap-4">
            <InputField v-model="form.judul" label="Judul Petunjuk" required />
            <Dropdown
              v-model="form.halaman"
              :options="halamanOptions"
              label="Halaman"
              required
            />
          </div>
          <div class="grid grid-cols-2 gap-4 mt-4">
            <InputField v-model.number="form.urutan" label="Urutan" type="number" />
            <div class="flex items-center gap-2 pt-6">
              <input
                type="checkbox"
                v-model="form.is_active"
                id="is_active"
                class="rounded text-blue-500"
              />
              <label for="is_active" class="text-sm text-gray-700">Aktif</label>
            </div>
          </div>
        </div>

        <!-- Image Upload Section -->
        <div class="bg-gray-50 rounded-lg p-4 border">
          <h3 class="font-semibold text-gray-700 mb-3 flex items-center gap-2">
            <span>🖼️</span> Gambar (Opsional)
          </h3>
          
          <input
            ref="imageInput"
            type="file"
            accept="image/jpeg,image/png,image/gif,image/webp"
            class="hidden"
            @change="handleImageUpload"
          />
          
          <div class="flex items-start gap-4">
            <!-- Image Preview -->
            <div v-if="form.image_url" class="relative">
              <div class="p-2 bg-white rounded-lg border-2 border-dashed border-gray-200">
                <img 
                  :src="getImageUrl(form.image_url)" 
                  alt="Preview"
                  class="object-contain rounded"
                  :style="{ height: form.image_size + 'px', maxWidth: '300px' }"
                />
              </div>
              <div class="mt-2 flex gap-2">
                <button
                  type="button"
                  @click="triggerImageUpload"
                  :disabled="uploadingImage"
                  class="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200"
                >
                  Ganti
                </button>
                <button
                  type="button"
                  @click="removeImage"
                  class="px-3 py-1.5 text-sm bg-red-100 text-red-700 rounded-lg hover:bg-red-200"
                >
                  Hapus
                </button>
              </div>
            </div>
            
            <!-- Upload Button -->
            <div v-else>
              <button
                type="button"
                @click="triggerImageUpload"
                :disabled="uploadingImage"
                class="flex flex-col items-center justify-center w-48 h-32 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all"
              >
                <span v-if="uploadingImage" class="animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></span>
                <template v-else>
                  <span class="text-3xl text-gray-400 mb-2">📷</span>
                  <span class="text-sm text-gray-500">Klik untuk upload</span>
                  <span class="text-xs text-gray-400 mt-1">JPG, PNG, GIF (maks 5MB)</span>
                </template>
              </button>
            </div>
            
            <!-- Image Size Control -->
            <div v-if="form.image_url" class="flex-shrink-0">
              <label class="block text-sm font-medium text-gray-700 mb-2">Ukuran Gambar (px)</label>
              <input
                v-model.number="form.image_size"
                type="number"
                min="50"
                max="500"
                class="w-24 px-3 py-2 border border-gray-300 rounded-lg"
                placeholder="200"
              />
            </div>
          </div>
        </div>

        <!-- Sections Editor -->
        <div class="bg-gray-50 rounded-lg p-4 border">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-semibold text-gray-700 flex items-center gap-2">
              <span>📋</span> Konten Petunjuk
            </h3>
            <button
              type="button"
              @click="addSection"
              class="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 flex items-center gap-1"
            >
              <span>+</span> Tambah Seksi
            </button>
          </div>
          
          <!-- Empty State -->
          <div v-if="infoSections.length === 0" class="p-8 border-2 border-dashed border-gray-200 rounded-lg text-center bg-white">
            <span class="text-3xl text-gray-300">📋</span>
            <p class="text-gray-500 mt-2">Belum ada konten. Klik "Tambah Seksi" untuk menambahkan.</p>
          </div>
          
          <!-- Section Cards -->
          <div class="space-y-3">
            <div
              v-for="(section, sectionIndex) in infoSections"
              :key="sectionIndex"
              class="bg-white border rounded-lg overflow-hidden"
            >
              <!-- Section Header -->
              <div class="bg-gray-100 px-3 py-2 flex items-center gap-2 border-b">
                <div class="flex flex-col gap-0.5">
                  <button type="button" @click="moveSectionUp(sectionIndex)" :disabled="sectionIndex === 0" class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs">▲</button>
                  <button type="button" @click="moveSectionDown(sectionIndex)" :disabled="sectionIndex === infoSections.length - 1" class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs">▼</button>
                </div>
                
                <div class="w-1 h-8 rounded-full flex-shrink-0" :style="{ backgroundColor: section.color }"></div>
                
                <select v-model="section.icon" class="w-12 px-1 py-1 border border-gray-200 rounded text-center text-sm">
                  <option v-for="icon in sectionIcons" :key="icon.value" :value="icon.value">{{ icon.value }}</option>
                </select>
                
                <input v-model="section.title" type="text" class="flex-1 px-2 py-1 border border-gray-200 rounded text-sm" placeholder="Judul Seksi" />
                
                <select v-model="section.color" class="px-2 py-1 border border-gray-200 rounded text-sm">
                  <option v-for="color in sectionColors" :key="color.value" :value="color.value">{{ color.label }}</option>
                </select>
                
                <button type="button" @click="removeSection(sectionIndex)" class="text-red-500 hover:text-red-700 p-1">🗑️</button>
              </div>
              
              <!-- Section Items -->
              <div class="p-3 space-y-2">
                <div v-for="(item, itemIndex) in section.items" :key="itemIndex" class="flex items-center gap-2">
                  <span class="w-5 h-5 rounded-full bg-gray-100 text-gray-500 flex items-center justify-center text-xs">{{ itemIndex + 1 }}</span>
                  <input v-model="item.text" type="text" class="flex-1 px-2 py-1.5 border border-gray-200 rounded text-sm" placeholder="Teks..." />
                  <label class="flex items-center gap-1 text-xs">
                    <input type="checkbox" v-model="item.isBold" class="rounded" />
                    <span class="text-gray-500">Bold</span>
                  </label>
                  <input v-model="item.linkText" type="text" class="w-20 px-2 py-1 border border-gray-200 rounded text-xs" placeholder="Link teks" />
                  <button type="button" @click="removeItem(sectionIndex, itemIndex)" class="text-red-400 hover:text-red-600">✕</button>
                </div>
                
                <button type="button" @click="addItem(sectionIndex)" class="w-full py-1.5 border-2 border-dashed border-gray-200 rounded text-gray-400 hover:border-blue-300 hover:text-blue-500 text-sm">
                  + Tambah Item
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 border rounded-lg hover:bg-gray-50">
          Batal
        </button>
        <button @click="save" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          Simpan
        </button>
      </template>
    </Modal>

    <!-- Delete Modal -->
    <Modal :show="showDeleteModal" title="Konfirmasi Hapus" @close="showDeleteModal = false">
      <p>Yakin ingin menghapus petunjuk ini?</p>
      <template #footer>
        <button @click="showDeleteModal = false" class="px-4 py-2 border rounded-lg hover:bg-gray-50">Batal</button>
        <button @click="deletePetunjuk" class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700">Hapus</button>
      </template>
    </Modal>

    <!-- Bulk Delete Modal -->
    <Modal :show="showBulkDeleteModal" title="Konfirmasi Hapus Massal" @close="showBulkDeleteModal = false">
      <p>Yakin ingin menghapus <strong>{{ selectedIds.length }}</strong> petunjuk?</p>
      <template #footer>
        <button @click="showBulkDeleteModal = false" class="px-4 py-2 border rounded-lg hover:bg-gray-50">Batal</button>
        <button @click="confirmBulkDelete" class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700">Hapus</button>
      </template>
    </Modal>
  </div>
</template>
