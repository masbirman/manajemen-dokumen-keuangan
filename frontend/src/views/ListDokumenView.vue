<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { useAuthStore } from "@/stores/auth";
import { DataTable, Modal, Dropdown, InputField } from "@/components/ui";

interface Dokumen {
  id: string;
  nomor_dokumen: string;
  tanggal_dokumen: string;
  uraian: string;
  nilai: number;
  file_path: string;
  created_by: string;
  unit_kerja?: { nama: string };
  pptk?: { nama: string };
  sumber_dana?: { nama: string };
  jenis_dokumen?: { nama: string };
  creator?: { name: string };
}

interface UnitKerja {
  id: string;
  nama: string;
}

const toast = useToast();
const authStore = useAuthStore();
const router = useRouter();
const loading = ref(false);
const data = ref<Dokumen[]>([]);
const unitKerjaList = ref<UnitKerja[]>([]);
const showDetailModal = ref(false);
const showBulkDeleteModal = ref(false);
const showDeleteModal = ref(false);
const selectedDokumen = ref<Dokumen | null>(null);
const selectedIds = ref<string[]>([]);
const dokumenToDelete = ref<Dokumen | null>(null);

// Pagination
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(10);

const filters = ref({
  unit_kerja_id: "" as string,
  start_date: "",
  end_date: "",
});

const columns = [
  { key: "nomor_dokumen", label: "No. Dokumen", width: "180px" },
  { key: "tanggal_dokumen", label: "Tanggal", width: "120px" },
  { key: "uraian", label: "Uraian" },
  { key: "nilai", label: "Nilai", width: "150px" },
  { key: "unit_kerja", label: "Unit Kerja", width: "150px" },
  { key: "jenis_dokumen", label: "Jenis", width: "120px" },
  { key: "creator", label: "Operator", width: "150px" },
  { key: "actions", label: "Aksi", width: "130px" },
];

const unitKerjaOptions = computed(() => [
  { value: "", label: "Semua Unit Kerja" },
  ...unitKerjaList.value.map((u) => ({ value: u.id, label: u.nama })),
]);

const formatCurrency = (value: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR" }).format(
    value
  );
const formatDate = (date: string) => new Date(date).toLocaleDateString("id-ID");

const fetchData = async (page = 1) => {
  loading.value = true;
  try {
    const params: Record<string, string | number> = {
      page,
      page_size: perPage.value,
    };
    if (filters.value.unit_kerja_id)
      params.unit_kerja_id = filters.value.unit_kerja_id;
    if (filters.value.start_date) params.start_date = filters.value.start_date;
    if (filters.value.end_date) params.end_date = filters.value.end_date;

    const response = await api.get("/dokumen", { params });
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

const onPageChange = (page: number) => fetchData(page);

const onPerPageChange = (newPerPage: number) => {
  perPage.value = newPerPage;
  fetchData(1);
};

const fetchUnitKerja = async () => {
  try {
    const response = await api.get("/unit-kerja/active");
    unitKerjaList.value = response.data.data || response.data || [];
  } catch {
    // ignore
  }
};

const openDetail = (dokumen: Dokumen) => {
  selectedDokumen.value = dokumen;
  showDetailModal.value = true;
};

const viewFile = async (dokumen: Dokumen) => {
  try {
    const response = await api.get(`/dokumen/${dokumen.id}/file`, {
      responseType: "blob",
    });
    const url = window.URL.createObjectURL(
      new Blob([response.data], { type: "application/pdf" })
    );
    window.open(url, "_blank");
  } catch {
    toast.error("Gagal membuka file");
  }
};

const onSelectionChange = (ids: (string | number)[]) => {
  selectedIds.value = ids as string[];
};

const confirmBulkDelete = async () => {
  try {
    await Promise.all(
      selectedIds.value.map((id) => api.delete(`/dokumen/${id}`))
    );
    toast.success(`${selectedIds.value.length} dokumen berhasil dihapus`);
    selectedIds.value = [];
    showBulkDeleteModal.value = false;
    fetchData();
  } catch {
    toast.error("Gagal menghapus dokumen");
  }
};

const applyFilter = () => fetchData(1);

const canEditDelete = (dokumen: Dokumen) => {
  // Admin/SuperAdmin can edit/delete all
  if (authStore.isAdmin) return true;
  // Operator can only edit/delete their own documents
  return dokumen.created_by === authStore.user?.id;
};

const editDokumen = (dokumen: Dokumen) => {
  router.push(`/dokumen/edit/${dokumen.id}`);
};

const confirmDelete = (dokumen: Dokumen) => {
  dokumenToDelete.value = dokumen;
  showDeleteModal.value = true;
};

const deleteDokumen = async () => {
  if (!dokumenToDelete.value) return;
  try {
    await api.delete(`/dokumen/${dokumenToDelete.value.id}`);
    toast.success("Dokumen berhasil dihapus");
    showDeleteModal.value = false;
    dokumenToDelete.value = null;
    fetchData(currentPage.value);
  } catch {
    toast.error("Gagal menghapus dokumen");
  }
};

onMounted(() => {
  fetchUnitKerja();
  fetchData();
});
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">List Dokumen</h1>
      <RouterLink
        to="/dokumen/input"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        + Input Dokumen
      </RouterLink>
    </div>

    <div class="bg-white rounded-lg shadow-sm p-4 mb-4 border border-gray-200">
      <div class="grid grid-cols-4 gap-4 items-end">
        <Dropdown
          v-model="filters.unit_kerja_id"
          :options="unitKerjaOptions"
          label="Unit Kerja"
          searchable
        />
        <InputField
          v-model="filters.start_date"
          label="Dari Tanggal"
          type="date"
        />
        <InputField
          v-model="filters.end_date"
          label="Sampai Tanggal"
          type="date"
        />
        <button
          @click="applyFilter"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 h-10"
        >
          Filter
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
      :selected-ids="selectedIds"
      @selection-change="onSelectionChange"
      @page-change="onPageChange"
      @per-page-change="onPerPageChange"
      @row-click="openDetail"
    >
      <template #bulk-actions>
        <button
          @click="showBulkDeleteModal = true"
          class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700"
        >
          🗑️ Hapus Terpilih
        </button>
      </template>
      <template #nomor_dokumen="{ value }">{{
        (value as string) || "-"
      }}</template>
      <template #tanggal_dokumen="{ value }">{{
        value ? formatDate(value as string) : "-"
      }}</template>
      <template #nilai="{ value }">{{
        formatCurrency(value as number)
      }}</template>
      <template #unit_kerja="{ row }">{{
        (row as Dokumen).unit_kerja?.nama || "-"
      }}</template>
      <template #jenis_dokumen="{ row }">{{
        (row as Dokumen).jenis_dokumen?.nama || "-"
      }}</template>
      <template #creator="{ row }">{{
        (row as Dokumen).creator?.name || "-"
      }}</template>
      <template #actions="{ row }">
        <div class="flex gap-1">
          <button
            @click.stop="viewFile(row as Dokumen)"
            class="px-2 py-1 text-xs bg-green-100 text-green-600 rounded hover:bg-green-200"
            title="Lihat PDF"
          >
            📄
          </button>
          <template v-if="canEditDelete(row as Dokumen)">
            <button
              @click.stop="editDokumen(row as Dokumen)"
              class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded hover:bg-blue-200"
              title="Edit"
            >
              ✏️
            </button>
            <button
              @click.stop="confirmDelete(row as Dokumen)"
              class="px-2 py-1 text-xs bg-red-100 text-red-600 rounded hover:bg-red-200"
              title="Hapus"
            >
              🗑️
            </button>
          </template>
        </div>
      </template>
    </DataTable>

    <Modal
      :show="showDetailModal"
      title="Detail Dokumen"
      size="lg"
      @close="showDetailModal = false"
    >
      <div v-if="selectedDokumen" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-gray-500">Nomor Dokumen</p>
            <p class="font-medium">{{ selectedDokumen.nomor_dokumen }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Tanggal</p>
            <p class="font-medium">
              {{ formatDate(selectedDokumen.tanggal_dokumen) }}
            </p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Unit Kerja</p>
            <p class="font-medium">
              {{ selectedDokumen.unit_kerja?.nama || "-" }}
            </p>
          </div>
          <div>
            <p class="text-sm text-gray-500">PPTK</p>
            <p class="font-medium">{{ selectedDokumen.pptk?.nama || "-" }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Sumber Dana</p>
            <p class="font-medium">
              {{ selectedDokumen.sumber_dana?.nama || "-" }}
            </p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Jenis Dokumen</p>
            <p class="font-medium">
              {{ selectedDokumen.jenis_dokumen?.nama || "-" }}
            </p>
          </div>
        </div>
        <div>
          <p class="text-sm text-gray-500">Uraian</p>
          <p class="font-medium">{{ selectedDokumen.uraian }}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Nilai</p>
          <p class="font-medium text-lg text-blue-600">
            {{ formatCurrency(selectedDokumen.nilai) }}
          </p>
        </div>
      </div>
      <template #footer>
        <button
          @click="showDetailModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Tutup
        </button>
        <button
          @click="viewFile(selectedDokumen!)"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          📄 Lihat PDF
        </button>
      </template>
    </Modal>

    <Modal
      :show="showBulkDeleteModal"
      title="Konfirmasi Hapus Massal"
      @close="showBulkDeleteModal = false"
    >
      <p>
        Yakin ingin menghapus <strong>{{ selectedIds.length }}</strong> dokumen
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

    <Modal
      :show="showDeleteModal"
      title="Konfirmasi Hapus"
      @close="showDeleteModal = false"
    >
      <p>
        Yakin ingin menghapus dokumen
        <strong>{{
          dokumenToDelete?.nomor_dokumen || dokumenToDelete?.uraian
        }}</strong
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
          @click="deleteDokumen"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Hapus
        </button>
      </template>
    </Modal>
  </div>
</template>
