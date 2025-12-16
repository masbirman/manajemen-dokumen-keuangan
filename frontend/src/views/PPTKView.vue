<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { DataTable, Modal, InputField, Dropdown } from "@/components/ui";
import { useAuthStore } from "@/stores/auth";

interface UnitKerja {
  id: string;
  kode: string;
  nama: string;
}
interface PPTK {
  id: string;
  nip: string;
  nama: string;
  jabatan: string;
  unit_kerja_id: string;
  unit_kerja?: UnitKerja;
  avatar_path?: string;
}

const authStore = useAuthStore();
const toast = useToast();
const loading = ref(false);
const data = ref<PPTK[]>([]);
const unitKerjaList = ref<UnitKerja[]>([]);
const showModal = ref(false);
const showDeleteModal = ref(false);
const showBulkDeleteModal = ref(false);
const showAvatarModal = ref(false);
const editingItem = ref<PPTK | null>(null);
const deletingItem = ref<PPTK | null>(null);
const avatarPptk = ref<PPTK | null>(null);
const selectedIds = ref<string[]>([]);
const avatarFile = ref<File | null>(null);
const avatarPreview = ref<string>("");

// Pagination & Search
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(10);
const searchQuery = ref("");
const filterUnitKerja = ref("");

const form = ref({ nip: "", nama: "", jabatan: "", unit_kerja_id: "" });
const errors = ref<Record<string, string>>({});

const columns = computed(() => {
  const cols = [
    { key: "avatar", label: "", width: "60px" },
    { key: "nip", label: "NIP", width: "150px" },
    { key: "nama", label: "Nama" },
    { key: "jabatan", label: "Jabatan" },
    { key: "unit_kerja", label: "Unit Kerja" },
  ];
  
  if (authStore.isSuperAdmin) {
    cols.push({ key: "actions", label: "Aksi", width: "180px" });
  }

  return cols;
});

const unitKerjaOptions = computed(() =>
  unitKerjaList.value.map((uk) => ({ value: uk.id, label: uk.nama }))
);

const filterUnitKerjaOptions = computed(() => [
  { value: "", label: "Semua Unit Kerja" },
  ...unitKerjaList.value.map((uk) => ({ value: uk.id, label: uk.nama })),
]);

const getAvatarUrl = (path?: string) => {
  if (!path) return "";
  return `${
    import.meta.env.VITE_API_URL || "http://localhost:8000/api"
  }/files/${path}`;
};

const fetchData = async (page = 1) => {
  loading.value = true;
  try {
    const params: Record<string, unknown> = { page, page_size: perPage.value };
    if (searchQuery.value) params.search = searchQuery.value;
    if (filterUnitKerja.value) params.unit_kerja_id = filterUnitKerja.value;

    const [pptkRes, ukRes] = await Promise.all([
      api.get("/pptk", { params }),
      api.get("/unit-kerja/active"),
    ]);
    data.value = pptkRes.data.data || pptkRes.data || [];
    currentPage.value = pptkRes.data.page || 1;
    totalPages.value = pptkRes.data.total_pages || 1;
    totalItems.value = pptkRes.data.total || 0;
    unitKerjaList.value = ukRes.data.data || ukRes.data || [];
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

const onSearch = (query: string) => {
  searchQuery.value = query;
  fetchData(1);
};

const onFilterUnitKerja = () => {
  fetchData(1);
};

const openCreate = () => {
  editingItem.value = null;
  form.value = { nip: "", nama: "", jabatan: "", unit_kerja_id: "" };
  errors.value = {};
  showModal.value = true;
};

const openEdit = (item: PPTK) => {
  editingItem.value = item;
  form.value = {
    nip: item.nip,
    nama: item.nama,
    jabatan: item.jabatan,
    unit_kerja_id: item.unit_kerja_id,
  };
  errors.value = {};
  showModal.value = true;
};

const openDelete = (item: PPTK) => {
  deletingItem.value = item;
  showDeleteModal.value = true;
};

const openAvatarUpload = (pptk: PPTK) => {
  avatarPptk.value = pptk;
  avatarFile.value = null;
  avatarPreview.value = pptk.avatar_path ? getAvatarUrl(pptk.avatar_path) : "";
  showAvatarModal.value = true;
};

const onAvatarChange = (e: Event) => {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    avatarFile.value = file;
    avatarPreview.value = URL.createObjectURL(file);
  }
};

const uploadAvatar = async () => {
  if (!avatarPptk.value || !avatarFile.value) return;
  const formData = new FormData();
  formData.append("avatar", avatarFile.value);
  try {
    await api.post(`/pptk/${avatarPptk.value.id}/avatar`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    toast.success("Avatar berhasil diupload");
    showAvatarModal.value = false;
    fetchData(currentPage.value);
  } catch {
    toast.error("Gagal upload avatar");
  }
};

const save = async () => {
  errors.value = {};
  try {
    if (editingItem.value) {
      await api.put(`/pptk/${editingItem.value.id}`, form.value);
      toast.success("Data berhasil diupdate");
    } else {
      await api.post("/pptk", form.value);
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
    await api.delete(`/pptk/${deletingItem.value.id}`);
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
    await Promise.all(selectedIds.value.map((id) => api.delete(`/pptk/${id}`)));
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
    const response = await api.get("/pptk/template", { responseType: "blob" });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", "template_pptk.xlsx");
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
    await api.post("/pptk/import", formData, {
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
      <h1 class="text-2xl font-bold text-gray-800">PPTK</h1>
      <div v-if="authStore.isSuperAdmin" class="flex gap-2">
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
      :selectable="authStore.isSuperAdmin"
      searchable
      search-placeholder="Cari NIP atau nama..."
      :selected-ids="selectedIds"
      @selection-change="onSelectionChange"
      @page-change="onPageChange"
      @per-page-change="onPerPageChange"
      @search="onSearch"
    >
      <template #filters>
        <Dropdown
          v-model="filterUnitKerja"
          :options="filterUnitKerjaOptions"
          placeholder="Filter Unit Kerja"
          class="w-48"
          @update:model-value="onFilterUnitKerja"
        />
      </template>
      <template v-if="authStore.isSuperAdmin" #bulk-actions>
        <button
          @click="showBulkDeleteModal = true"
          class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700"
        >
          🗑️ Hapus Terpilih
        </button>
      </template>
      <template #avatar="{ row }">
        <div
          class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden"
        >
          <img
            v-if="(row as PPTK).avatar_path"
            :src="getAvatarUrl((row as PPTK).avatar_path)"
            class="w-full h-full object-cover"
          />
          <span v-else class="text-gray-500 text-lg">{{
            (row as PPTK).nama?.charAt(0)?.toUpperCase()
          }}</span>
        </div>
      </template>
      <template #unit_kerja="{ row }">{{
        (row as PPTK).unit_kerja?.nama || "-"
      }}</template>
      <template v-if="authStore.isSuperAdmin" #actions="{ row }">
        <div class="flex gap-2">
          <button
            @click.stop="openAvatarUpload(row as PPTK)"
            class="text-green-600 hover:text-green-800"
          >
            📷
          </button>
          <button
            @click.stop="openEdit(row as PPTK)"
            class="text-blue-600 hover:text-blue-800"
          >
            Edit
          </button>
          <button
            @click.stop="openDelete(row as PPTK)"
            class="text-red-600 hover:text-red-800"
          >
            Hapus
          </button>
        </div>
      </template>
    </DataTable>

    <Modal
      :show="showModal"
      :title="editingItem ? 'Edit PPTK' : 'Tambah PPTK'"
      size="lg"
      @close="showModal = false"
    >
      <form @submit.prevent="save">
        <InputField
          v-model="form.nip"
          label="NIP"
          :error="errors.nip"
          required
        />
        <InputField
          v-model="form.nama"
          label="Nama"
          :error="errors.nama"
          required
        />
        <InputField
          v-model="form.jabatan"
          label="Jabatan"
          :error="errors.jabatan"
          required
        />
        <Dropdown
          v-model="form.unit_kerja_id"
          :options="unitKerjaOptions"
          label="Unit Kerja"
          :error="errors.unit_kerja_id"
          searchable
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
      :show="showAvatarModal"
      title="Upload Avatar PPTK"
      @close="showAvatarModal = false"
    >
      <div class="flex flex-col items-center gap-4">
        <div
          class="w-32 h-32 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden border-4 border-gray-300"
        >
          <img
            v-if="avatarPreview"
            :src="avatarPreview"
            class="w-full h-full object-cover"
          />
          <span v-else class="text-gray-400 text-4xl">{{
            avatarPptk?.nama?.charAt(0)?.toUpperCase()
          }}</span>
        </div>
        <input
          type="file"
          accept="image/*"
          @change="onAvatarChange"
          class="text-sm"
        />
        <p class="text-sm text-gray-500">Format: JPG, PNG. Max 2MB</p>
      </div>
      <template #footer>
        <button
          @click="showAvatarModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="uploadAvatar"
          :disabled="!avatarFile"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          Upload
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
