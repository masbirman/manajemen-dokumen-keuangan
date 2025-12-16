<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { DataTable, Modal, InputField } from "@/components/ui";

interface UnitKerja {
  id: string;
  kode: string;
  nama: string;
  created_at: string;
}

const toast = useToast();
const loading = ref(false);
const data = ref<UnitKerja[]>([]);
const showModal = ref(false);
const showDeleteModal = ref(false);
const showBulkDeleteModal = ref(false);
const editingItem = ref<UnitKerja | null>(null);
const deletingItem = ref<UnitKerja | null>(null);
const selectedIds = ref<string[]>([]);

// Pagination & Search
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(10);
const searchQuery = ref("");

const form = ref({ kode: "", nama: "" });
const errors = ref<Record<string, string>>({});

const columns = [
  { key: "kode", label: "Kode", width: "150px" },
  { key: "nama", label: "Nama Unit Kerja" },
  { key: "actions", label: "Aksi", width: "150px" },
];

const fetchData = async (page = 1) => {
  loading.value = true;
  try {
    const response = await api.get("/unit-kerja", {
      params: {
        page,
        page_size: perPage.value,
        search: searchQuery.value || undefined,
      },
    });
    data.value = response.data.data || response.data || [];
    currentPage.value = response.data.page || 1;
    totalPages.value = response.data.total_pages || 1;
    totalItems.value = response.data.total || 0;
  } catch {
    toast.error("Gagal memuat data");
  } finally {
    loading.value = false;
  }
};

const onPageChange = (page: number) => {
  fetchData(page);
};

const onPerPageChange = (newPerPage: number) => {
  perPage.value = newPerPage;
  fetchData(1);
};

const onSearch = (query: string) => {
  searchQuery.value = query;
  fetchData(1);
};

const openCreate = () => {
  editingItem.value = null;
  form.value = { kode: "", nama: "" };
  errors.value = {};
  showModal.value = true;
};

const openEdit = (item: UnitKerja) => {
  editingItem.value = item;
  form.value = { kode: item.kode, nama: item.nama };
  errors.value = {};
  showModal.value = true;
};

const openDelete = (item: UnitKerja) => {
  deletingItem.value = item;
  showDeleteModal.value = true;
};

const save = async () => {
  errors.value = {};
  try {
    if (editingItem.value) {
      await api.put(`/unit-kerja/${editingItem.value.id}`, form.value);
      toast.success("Data berhasil diupdate");
    } else {
      await api.post("/unit-kerja", form.value);
      toast.success("Data berhasil ditambahkan");
    }
    showModal.value = false;
    fetchData(currentPage.value);
  } catch (e: unknown) {
    const err = e as {
      response?: { data?: { errors?: Record<string, string> } };
    };
    errors.value = err.response?.data?.errors || {};
    toast.error("Gagal menyimpan data");
  }
};

const confirmDelete = async () => {
  if (!deletingItem.value) return;
  try {
    await api.delete(`/unit-kerja/${deletingItem.value.id}`);
    toast.success("Data berhasil dihapus");
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
      selectedIds.value.map((id) => api.delete(`/unit-kerja/${id}`))
    );
    toast.success(`${selectedIds.value.length} data berhasil dihapus`);
    selectedIds.value = [];
    showBulkDeleteModal.value = false;
    fetchData(currentPage.value);
  } catch {
    toast.error("Gagal menghapus data");
  }
};

const downloadTemplate = async () => {
  try {
    const response = await api.get("/unit-kerja/template", {
      responseType: "blob",
    });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", "template_unit_kerja.xlsx");
    document.body.appendChild(link);
    link.click();
    link.remove();
  } catch {
    toast.error("Gagal download template");
  }
};

const importFile = ref<HTMLInputElement | null>(null);
const handleImport = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  const formData = new FormData();
  formData.append("file", file);
  try {
    await api.post("/unit-kerja/import", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    toast.success("Import berhasil");
    fetchData(1);
  } catch {
    toast.error("Import gagal");
  }
  input.value = "";
};

onMounted(() => fetchData());
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Unit Kerja</h1>
      <div class="flex gap-2">
        <button
          @click="downloadTemplate"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
        >
          📥 Template Excel
        </button>
        <input
          ref="importFile"
          type="file"
          accept=".xlsx"
          class="hidden"
          @change="handleImport"
        />
        <button
          @click="importFile?.click()"
          class="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700"
        >
          📤 Import Excel
        </button>
        <button
          @click="openCreate"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          + Tambah
        </button>
      </div>
    </div>

    <DataTable
      :columns="columns"
      :data="data"
      :loading="loading"
      :current-page="currentPage"
      :total-pages="totalPages"
      :total-items="totalItems"
      :per-page="perPage"
      selectable
      searchable
      search-placeholder="Cari kode atau nama..."
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
      <template #actions="{ row }">
        <div class="flex gap-2">
          <button
            @click.stop="openEdit(row as UnitKerja)"
            class="text-blue-600 hover:text-blue-800"
          >
            Edit
          </button>
          <button
            @click.stop="openDelete(row as UnitKerja)"
            class="text-red-600 hover:text-red-800"
          >
            Hapus
          </button>
        </div>
      </template>
    </DataTable>

    <Modal
      :show="showModal"
      :title="editingItem ? 'Edit Unit Kerja' : 'Tambah Unit Kerja'"
      @close="showModal = false"
    >
      <form @submit.prevent="save">
        <InputField
          v-model="form.kode"
          label="Kode"
          :error="errors.kode"
          required
        />
        <InputField
          v-model="form.nama"
          label="Nama"
          :error="errors.nama"
          required
        />
      </form>
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

    <Modal
      :show="showDeleteModal"
      title="Konfirmasi Hapus"
      @close="showDeleteModal = false"
    >
      <p>
        Yakin ingin menghapus <strong>{{ deletingItem?.nama }}</strong
        >?
      </p>
      <template #footer>
        <button
          @click="showDeleteModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="confirmDelete"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Hapus
        </button>
      </template>
    </Modal>

    <Modal
      :show="showBulkDeleteModal"
      title="Konfirmasi Hapus Massal"
      @close="showBulkDeleteModal = false"
    >
      <p>
        Yakin ingin menghapus <strong>{{ selectedIds.length }}</strong> data
        yang dipilih?
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
          Hapus Semua
        </button>
      </template>
    </Modal>
  </div>
</template>
