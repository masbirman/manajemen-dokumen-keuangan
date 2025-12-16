<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import {
  DataTable,
  Modal,
  InputField,
  Dropdown,
  RichEditor,
} from "@/components/ui";

interface Petunjuk {
  id: string;
  judul: string;
  konten: string;
  halaman: string;
  urutan: number;
  is_active: boolean;
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

const form = ref({
  judul: "",
  konten: "",
  halaman: "input_dokumen",
  urutan: 0,
  is_active: true,
});

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
  };
  showModal.value = true;
};

const openEdit = (item: Petunjuk) => {
  editingId.value = item.id;
  form.value = { ...item };
  showModal.value = true;
};

const save = async () => {
  try {
    if (editingId.value) {
      await api.put(`/petunjuk/${editingId.value}`, form.value);
      toast.success("Petunjuk berhasil diupdate");
    } else {
      await api.post("/petunjuk", form.value);
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
      size="2xl"
      @close="showModal = false"
    >
      <div class="space-y-4">
        <InputField v-model="form.judul" label="Judul" required />
        <RichEditor
          v-model="form.konten"
          label="Konten"
          placeholder="Tulis petunjuk di sini..."
          required
        />
        <Dropdown
          v-model="form.halaman"
          :options="halamanOptions"
          label="Halaman"
          required
        />
        <InputField v-model.number="form.urutan" label="Urutan" type="number" />
        <div class="flex items-center gap-2">
          <input
            type="checkbox"
            v-model="form.is_active"
            id="is_active"
            class="rounded"
          />
          <label for="is_active" class="text-sm text-gray-700">Aktif</label>
        </div>
      </div>
      <template #footer>
        <button
          @click="showModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="save"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          Simpan
        </button>
      </template>
    </Modal>

    <!-- Delete Modal -->
    <Modal
      :show="showDeleteModal"
      title="Konfirmasi Hapus"
      @close="showDeleteModal = false"
    >
      <p>Yakin ingin menghapus petunjuk ini?</p>
      <template #footer>
        <button
          @click="showDeleteModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="deletePetunjuk"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Hapus
        </button>
      </template>
    </Modal>

    <!-- Bulk Delete Modal -->
    <Modal
      :show="showBulkDeleteModal"
      title="Konfirmasi Hapus Massal"
      @close="showBulkDeleteModal = false"
    >
      <p>
        Yakin ingin menghapus
        <strong>{{ selectedIds.length }}</strong> petunjuk?
      </p>
      <template #footer>
        <button
          @click="showBulkDeleteModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="confirmBulkDelete"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Hapus
        </button>
      </template>
    </Modal>
  </div>
</template>
